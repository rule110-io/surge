package surge

import (
	"bufio"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	bitmap "github.com/boljen/go-bitmap"
	nkn "github.com/nknorg/nkn-sdk-go"
	pb "github.com/rule110-io/surge-ui/payloads"
	"google.golang.org/protobuf/proto"

	open "github.com/skratchdot/open-golang/open"
)

//TestTopic only for testing
const TestTopic = "poctest"

const surgeChunkID byte = 0x001
const surgeQueryRequestID byte = 0x002
const surgeQueryResponseID byte = 0x003

var queryPayload = ""
var fileWriteLock = &sync.Mutex{}

//OpenOSPath Open a file, directory, or URI using the OS's default application for that object type. Don't wait for the open command to complete.
func OpenOSPath(Path string) {
	open.Start(Path)
}

//OpenFileByHash opens a file with OS default application for object type
func OpenFileByHash(Hash string) {
	defer RecoverAndLog()
	fileInfo, err := dbGetFile(Hash)
	if err != nil {
		pushError("Error on open file", err.Error())
		return
	}
	OpenOSPath(fileInfo.Path)
}

//OpenFolderByHash opens the folder containing the file by hash in os
func OpenFolderByHash(Hash string) {
	defer RecoverAndLog()
	fileInfo, err := dbGetFile(Hash)
	if err != nil {
		pushError("Error on open folder", err.Error())
		return
	}
	OpenOSPath(filepath.Dir(fileInfo.Path))
}

// RequestChunk sends a request to an address for a specific chunk of a specific file
func RequestChunk(Session *Session, FileID string, ChunkID int32) bool {
	if Session == nil || Session.session == nil {
		return false
	}

	defer RecoverAndLog()
	msg := &pb.SurgeMessage{
		FileID:  FileID,
		ChunkID: ChunkID,
	}
	msgSerialized, err := proto.Marshal(msg)
	if err != nil {
		log.Fatalln("Failed to encode surge message:", err)
	} else {
		err := SessionWrite(Session, msgSerialized, surgeChunkID) //Client.Send(nkn.NewStringArray(Addr), msgSerialized, nil)
		if err != nil {
			log.Error("Failed to request chunk", err)
			return false
		}
	}

	return true
}

// TransmitChunk tranmits target file chunk to address
func TransmitChunk(Session *Session, FileID string, ChunkID int32) {
	defer RecoverAndLog()
	//Open file

	fileWriteLock.Lock()
	fileInfo, err := dbGetFile(FileID)
	if err != nil {
		log.Error("Error on transmit chunk - file not in db", err.Error())
		return
	}
	fileInfo.ChunksShared++
	dbInsertFile(*fileInfo)
	fileWriteLock.Unlock()

	file, err := os.Open(fileInfo.Path)

	//When we have an OS read error on the file mark the file as missing, stop down and uploads on it
	if err != nil {
		log.Error("Error on transmit chunk - file read failure", err.Error())

		fileWriteLock.Lock()
		fileInfo.IsMissing = true
		fileInfo.IsDownloading = false
		fileInfo.IsUploading = false
		dbInsertFile(*fileInfo)
		fileWriteLock.Unlock()

		return
	}
	defer file.Close()

	//Read the requested chunk
	chunkOffset := int64(ChunkID) * ChunkSize
	buffer := make([]byte, ChunkSize)
	bytesread, err := file.ReadAt(buffer, chunkOffset)

	if err != nil {
		if err != io.EOF {
			log.Error("Error on transmit chunk - read chunk failed: ", ChunkID, err.Error())
			return
		}
	}

	//Create the proto data
	dataReply := &pb.SurgeMessage{
		FileID:  FileID,
		ChunkID: ChunkID,
		Data:    buffer[:bytesread],
	}
	dateReplySerialized, err := proto.Marshal(dataReply)
	if err != nil {
		log.Error("Error on transmit chunk - chunk serialization error", err.Error())
		return
	}

	//Transmit the chunk
	err = SessionWrite(Session, dateReplySerialized, surgeChunkID) //Client.Send(nkn.NewStringArray(Addr), dateReplySerialized, nil)
	if err != nil {
		log.Error("Error on transmit chunk - failed to write to session", err.Error())
		return
	}
	log.Println("Chunk transmitted: ", bytesread, " bytes")
	Session.Uploaded += int64(bytesread)
}

// SendQueryRequest sends a query to a client on session
func SendQueryRequest(Addr string, Query string) bool {
	defer RecoverAndLog()

	var surgeSession *Session = nil

	//Check for sessions
	sessionsWriteLock.Lock()
	for i := 0; i < len(Sessions); i++ {
		if Sessions[i].session.RemoteAddr().String() == Addr {
			surgeSession = Sessions[i]
			break
		}
	}
	sessionsWriteLock.Unlock()

	if surgeSession == nil {
		//Create session
		sessionConfig := nkn.GetDefaultSessionConfig()
		sessionConfig.MTU = 16384
		sessionConfig.CheckTimeoutInterval = 1
		sessionConfig.InitialRetransmissionTimeout = 1
		sessionConfig.MaxRetransmissionTimeout = 1

		dialConfig := &nkn.DialConfig{
			SessionConfig: sessionConfig,
			DialTimeout:   60000,
		}

		downloadSession, err := client.DialWithConfig(Addr, dialConfig)
		if err != nil {
			fmt.Println(string("\033[35m"), "Peer with address is not online, stopped trying after 60000ms", Addr, string("\033[0m"))
			go setClientOnlineMap(Addr, false)
			return false
		}
		fmt.Println(string("\033[36m"), "Connected to peer %s requesting file listings", Addr, string("\033[0m"))

		go setClientOnlineMap(Addr, true)

		downloadReader := bufio.NewReader(downloadSession)

		surgeSession = &Session{
			reader:  downloadReader,
			session: downloadSession,
		}
		go initiateSession(surgeSession)
	}

	msg := &pb.SurgeQuery{
		Query: Query,
	}
	msgSerialized, err := proto.Marshal(msg)
	if err != nil {
		log.Fatalln("Failed to encode surge message:", err)
		return false
	}

	err = SessionWrite(surgeSession, msgSerialized, surgeQueryRequestID) //Client.Send(nkn.NewStringArray(Addr), msgSerialized, nil)
	if err != nil {
		log.Println("Failed to send Surge Request:", err)
		return false
	}

	return true
}

// SendQueryResponse sends a query to a client on session
func SendQueryResponse(Session *Session, Query string) {
	defer RecoverAndLog()
	b := []byte(queryPayload)
	err := SessionWrite(Session, b, surgeQueryResponseID) //Client.Send(nkn.NewStringArray(Addr), msgSerialized, nil)
	if err != nil {
		log.Println("Failed to send Surge Ruquest:", err)
	}
}

// AllocateFile Allocates a file on disk at path with size in bytes
func AllocateFile(path string, size int64) {
	defer RecoverAndLog()
	fd, err := os.Create(path)
	if err != nil {
		log.Fatal("Failed to create output")
	}
	_, err = fd.Seek(size-1, 0)
	if err != nil {
		log.Fatal("Failed to seek")
	}
	_, err = fd.Write([]byte{0})
	if err != nil {
		log.Fatal("Write failed")
	}
	err = fd.Close()
	if err != nil {
		log.Fatal("Failed to close file")
	}
}

func listenForSession() {
	defer RecoverAndLog()
	for !client.IsClosed() {
		listenSession, err := client.Accept()
		if err != nil {
			pushError("Error on client accept", err.Error())
			continue
		}

		var surgeSession *Session = nil
		//Check for existing sessions that are not for a file
		for i := 0; i < len(Sessions); i++ {
			if Sessions[i].session.RemoteAddr().String() == listenSession.RemoteAddr().String() && Sessions[i].FileSize == 0 {
				closeSession(Sessions[i])
			}
		}
		if surgeSession == nil {
			listenReader := bufio.NewReader(listenSession)
			surgeSession = &Session{
				reader:  listenReader,
				session: listenSession,
			}
			go initiateSession(surgeSession)
		}

		time.Sleep(time.Millisecond)
	}
}

// Listen will listen to incoming requests for chunks
func Listen() {
	go listenForSession()
}

func initiateSession(Session *Session) {
	defer RecoverAndLog()

	sessionsWriteLock.Lock()
	Sessions = append(Sessions, Session)
	sessionsWriteLock.Unlock()

	for true {
		data, chunkType, err := SessionRead(Session)
		if err != nil {
			log.Println("Session read failed, closing session error:", err)
			break
		}

		switch chunkType {
		case surgeChunkID:
			processChunk(Session, data)
			break
		case surgeQueryRequestID:
			processQueryRequest(Session, data)
			break
		case surgeQueryResponseID:
			processQueryResponse(Session, data)
			break
		}
	}

	closeSession(Session)
}

func closeSession(Session *Session) {
	defer RecoverAndLog()
	//find index in Sessions
	sessionsWriteLock.Lock()
	var index = -1
	for i := 0; i < len(Sessions); i++ {
		if Sessions[i] == Session {
			index = i
			break
		}
	}

	if index == -1 {
		log.Println("Session already removed")
		sessionsWriteLock.Unlock()
		return
	}

	//Close nkn session, nill out the pointers
	Session.session.Close()
	Session.session = nil
	Session.reader = nil
	if Session.file != nil {
		err := Session.file.Close()
		if err != nil {
			log.Println("File no longer exists")
		}
	}

	//Replace index of session to be removed with last element in slice
	Sessions[index] = Sessions[len(Sessions)-1]
	//Nul out the pointer to the surge session
	Sessions[len(Sessions)-1] = nil
	//Slice off the last element
	Sessions = Sessions[:len(Sessions)-1]
	sessionsWriteLock.Unlock()

	log.Println("-=Session closed=-")
}

func processChunk(Session *Session, Data []byte) {
	defer RecoverAndLog()
	//Try to parse SurgeMessage
	surgeMessage := &pb.SurgeMessage{}
	if err := proto.Unmarshal(Data, surgeMessage); err != nil {
		log.Fatalln("Failed to parse surge message:", err)
	}

	//If this is the first file data over this session we need to set the session
	if Session.FileHash == "" {
		dbFile, err := dbGetFile(surgeMessage.FileID)
		if err != nil {
			log.Println("Chunk requested by someone for a file which we do not have in our db")
			return
		}

		if !dbFile.IsUploading {
			log.Println("Chunk requested by someone for a file which we have not marked as uploading")
			return
		}

		Session.FileHash = dbFile.FileHash
		Session.FileSize = dbFile.FileSize
	}

	//Data nill means its a request for data
	if surgeMessage.Data == nil {
		go TransmitChunk(Session, surgeMessage.FileID, surgeMessage.ChunkID)
	} else { //If data is not nill we are receiving data

		//When we receive a chunk mark it as no longer in transit
		chunkKey := surgeMessage.FileID + "_" + strconv.Itoa(int(surgeMessage.ChunkID))

		chunkInTransitLock.Lock()
		chunksInTransit[chunkKey] = false
		chunkInTransitLock.Unlock()

		go WriteChunk(Session, surgeMessage.FileID, surgeMessage.ChunkID, surgeMessage.Data)
	}
}

func processQueryRequest(Session *Session, Data []byte) {
	defer RecoverAndLog()
	//Try to parse SurgeMessage
	surgeQuery := &pb.SurgeQuery{}
	if err := proto.Unmarshal(Data, surgeQuery); err != nil {
		log.Fatalln("Failed to parse surge message:", err)
	}
	log.Println("Query received", surgeQuery.Query)

	SendQueryResponse(Session, surgeQuery.Query)
}

func processQueryResponse(Session *Session, Data []byte) {
	defer RecoverAndLog()

	//Try to parse SurgeMessage
	s := string(Data)
	seeder := Session.session.RemoteAddr().String()

	fmt.Println(string("\033[36m"), "file query response received", seeder, string("\033[0m"))

	go setClientOnlineMap(seeder, true)

	ListedFilesLock.Lock()

	//Remove seeders that match current seeder from file listings
	for i := 0; i < len(ListedFiles); i++ {
		for j := 0; j < len(ListedFiles[i].Seeders); j++ {
			if ListedFiles[i].Seeders[j] == seeder {
				// Remove the element at index i from a.
				ListedFiles[i].Seeders[j] = ListedFiles[i].Seeders[len(ListedFiles[i].Seeders)-1] // Copy last element to index i.
				ListedFiles[i].Seeders[len(ListedFiles[i].Seeders)-1] = ""                        // Erase last element (write zero value).
				ListedFiles[i].Seeders = ListedFiles[i].Seeders[:len(ListedFiles[i].Seeders)-1]   // Truncate slice.
			}
		}
	}

	//Remove empty seeders listings
	/*
		for i := 0; i < len(ListedFiles); i++ {
			if len(ListedFiles[i].Seeders) == 0 {
				// Remove the element at index i from a.
				ListedFiles[i] = ListedFiles[len(ListedFiles)-1] // Copy last element to index i.
				ListedFiles[len(ListedFiles)-1] = File{}         // Erase last element (write zero value).
				ListedFiles = ListedFiles[:len(ListedFiles)-1]   // Truncate slice.
			}
		}*/

	//Parse the response
	payloadSplit := strings.Split(s, "surge://")
	for j := 0; j < len(payloadSplit); j++ {
		data := strings.Split(payloadSplit[j], "|")

		if len(data) < 3 {
			continue
		}

		fileSize, _ := strconv.ParseInt(data[3], 10, 64)
		numChunks := int((fileSize-1)/int64(ChunkSize)) + 1

		newListing := File{
			FileName:    data[2],
			FileSize:    fileSize,
			FileHash:    data[4],
			Seeders:     []string{seeder},
			Path:        "",
			NumChunks:   numChunks,
			ChunkMap:    nil,
			SeederCount: 1,
		}

		//Replace existing, or remove.
		var replace = false
		for l := 0; l < len(ListedFiles); l++ {
			if ListedFiles[l].FileHash == newListing.FileHash {

				//if the seeder is unique add it as an additional seeder for the file
				ListedFiles[l].Seeders = append(ListedFiles[l].Seeders, seeder)
				ListedFiles[l].SeederCount = len(ListedFiles[l].Seeders)
				replace = true
				break
			}
		}
		//Unique listing so we add
		if replace == false {
			ListedFiles = append(ListedFiles, newListing)
		}

		fmt.Println(string("\033[33m"), "Filename", newListing.FileName, "FileHash", newListing.FileHash, string("\033[0m"))

		log.Println("Query response new file: ", newListing.FileName, " seeder: ", seeder)

		//Test gui
		//newButton := widget.NewButton(newListing.Filename+" | "+ByteCountSI(newListing.FileSize), func() {
		//	downloadFile(newListing.Seeder, newListing.FileSize, newListing.Filename)
		//})
		//fileBox.Append(newButton)
	}
	ListedFilesLock.Unlock()
}

//ParsePayloadString parses payload of files
func ParsePayloadString(s string) {
	defer RecoverAndLog()
	payloadSplit := strings.Split(s, "surge://")
	for j := 0; j < len(payloadSplit); j++ {
		data := strings.Split(payloadSplit[j], "|")

		if len(data) < 3 {
			continue
		}

		fileSize, _ := strconv.ParseInt(data[3], 10, 64)
		numChunks := int((fileSize-1)/int64(ChunkSize)) + 1

		seeder := strings.Split(data[5], ",")

		newListing := File{
			FileName:  data[2],
			FileSize:  fileSize,
			FileHash:  data[4],
			Seeders:   seeder,
			Path:      "",
			NumChunks: numChunks,
			ChunkMap:  nil,
		}

		ListedFilesLock.Lock()
		//Replace existing, or remove.
		var replace = false
		for l := 0; l < len(ListedFiles); l++ {
			if ListedFiles[l].FileHash == newListing.FileHash {
				//if the seeder is unique add it as an additional seeder for the file
				ListedFiles[l].Seeders = append(ListedFiles[l].Seeders, seeder...)
				replace = true
				break
			}
		}
		//Unique listing so we add
		if replace == false {
			ListedFiles = append(ListedFiles, newListing)
		}
		ListedFilesLock.Unlock()

		log.Println("Program paramater new file: ", newListing.FileName, " seeder: ", newListing.Seeders)

		go DownloadFile(newListing.FileHash)

		//Test gui
		//newButton := widget.NewButton(newListing.Filename+" | "+ByteCountSI(newListing.FileSize), func() {
		//	downloadFile(newListing.Seeder, newListing.FileSize, newListing.Filename)
		//})
		//fileBox.Append(newButton)
	}
}

//WriteChunk writes a chunk to disk
func WriteChunk(Session *Session, FileID string, ChunkID int32, Chunk []byte) {
	defer RecoverAndLog()
	workerCount--

	if Session.file == nil {
		fileInfo, err := dbGetFile(FileID)
		if err != nil {
			pushError("Error on write chunk (db get)", err.Error())
			return
		}

		var path = remoteFolder + string(os.PathSeparator) + fileInfo.FileName

		//Open file
		Session.file, err = os.OpenFile(path, os.O_RDWR, 0644)
		if err != nil {
			pushError("Error on write chunk (os open)", err.Error())
			return
		}
	}

	chunkOffset := int64(ChunkID) * ChunkSize
	bytesWritten, err := Session.file.WriteAt(Chunk, chunkOffset)
	if err != nil {
		pushError("Error on write chunk (file write)", err.Error())
		return
	}
	//Success
	log.Println("Chunk written to disk: ", bytesWritten, " bytes")

	//Multiple sessions can be downloading this file so we add to all
	for i := 0; i < len(Sessions); i++ {
		if Sessions[i].FileHash == FileID {
			Sessions[i].Downloaded += int64(bytesWritten)
		}
	}

	//Update bitmap async as this has a lock in it but does not have to be waited for
	setBitMap := func() {
		fileWriteLock.Lock()

		//Set chunk to available in the map
		fileInfo, err := dbGetFile(FileID)
		if err != nil {
			pushError("Error on chunk write (db get)", err.Error())
			return
		}
		bitmap.Set(fileInfo.ChunkMap, int(ChunkID), true)
		dbInsertFile(*fileInfo)

		fileWriteLock.Unlock()
	}
	go setBitMap()
}

//TopicEncode .
func TopicEncode(topic string) string {
	return "SRG_" + strings.ReplaceAll(b64.StdEncoding.EncodeToString([]byte(topic)), "=", "-")
}

func surgeGenerateTopicPayload(fileName string, sizeInBytes int64, hash string) string {
	//Example payload
	//surge://|file|The_Two_Towers-The_Purist_Edit-Trailer.avi|14997504|965c013e991ee246d63d45ea71954c4d|/

	return "surge://|file|" + fileName + "|" + strconv.FormatInt(sizeInBytes, 10) + "|" + hash + "|/"
}

func surgeGenerateMagnetLink(fileName string, sizeInBytes int64, hash string, seeder string) string {
	//Example payload
	//surge://|file|The_Two_Towers-The_Purist_Edit-Trailer.avi|14997504|965c013e991ee246d63d45ea71954c4d|/
	if seeder == "" {
		seeder = client.Addr().String()
	}

	return "surge://|file|" + fileName + "|" + strconv.FormatInt(sizeInBytes, 10) + "|" + hash + "|" + seeder + "|/"
}

// HashFile generates hash for file given filepath
func HashFile(filePath string) (string, error) {
	defer RecoverAndLog()
	//Open the passed argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	//Tell the program to call the following function when the current function returns
	defer file.Close()

	//Open a new hash interface to write to
	hash := sha256.New()

	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)

	//Convert the bytes to a string
	hashString := hex.EncodeToString(hashInBytes)

	return hashString, nil

}

func surgeGetFileSize(path string) int64 {
	defer RecoverAndLog()
	fi, err := os.Stat(path)
	if err != nil {
		log.Panicln("Error on get filesize", err)
	}
	// get the size
	return fi.Size()
}

//BuildSeedString builds a string of seeded files to share with clients
func BuildSeedString(dbFiles []File) {
	defer RecoverAndLog()
	newQueryPayload := ""
	for _, dbFile := range dbFiles {
		magnet := surgeGenerateMagnetLink(dbFile.FileName, dbFile.FileSize, dbFile.FileHash, strings.Join(dbFile.Seeders, ","))
		log.Println("Magnet:", magnet)

		if dbFile.IsUploading {
			//Add to payload
			payload := surgeGenerateTopicPayload(dbFile.FileName, dbFile.FileSize, dbFile.FileHash)
			//log.Println(payload)
			newQueryPayload += payload
		}
	}
	queryPayload = newQueryPayload
}

//AddToSeedString adds to existing seed string
func AddToSeedString(dbFile File) {
	defer RecoverAndLog()
	//Add to payload
	payload := surgeGenerateTopicPayload(dbFile.FileName, dbFile.FileSize, dbFile.FileHash)
	//log.Println(payload)
	queryPayload += payload

	//Make sure you're subscribed when seeding a file
	go subscribeToSurgeTopic()
}

//SeedFile generates everything needed to seed a file
func SeedFile(Path string) bool {
	defer RecoverAndLog()
	log.Println("Seeding file", Path)

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	randomHash := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	fileName := filepath.Base(Path)
	fileSize := surgeGetFileSize(Path)
	numChunks := int((fileSize-1)/int64(ChunkSize)) + 1
	chunkMap := bitmap.NewSlice(numChunks)

	//Local files are always fully available, set all chunks to 1
	for i := 0; i < numChunks; i++ {
		bitmap.Set(chunkMap, i, true)
	}

	//Append to local files
	localFile := File{
		FileName:      fileName,
		FileSize:      fileSize,
		FileHash:      randomHash,
		Path:          Path,
		NumChunks:     numChunks,
		ChunkMap:      chunkMap,
		IsUploading:   false,
		IsDownloading: false,
		IsHashing:     true,
		SeederCount:   0,
	}

	//Check if file is already seeded
	_, err = dbGetFile(localFile.FileHash)
	if err == nil {
		//File already seeding
		pushError("Seed failed", fileName+" already seeding.")
		return false
	}

	//When seeding a new file enter file into db
	dbInsertFile(localFile)

	go hashFileJob(randomHash)

	return true
}

func hashFileJob(randomHash string) {
	//Clean up after were done here, even when we fail we dont want these randomhash files in db
	defer dbDeleteFile(randomHash)

	dbFile, err := dbGetFile(randomHash)
	if err != nil {
		pushError("File Hash Failed", "Could find dbEntry for hash "+randomHash)
	}

	hashString, err := HashFile(dbFile.Path)
	if err != nil {
		log.Println(err)
		pushError("File Hash Failed", "Could not hash file at "+dbFile.Path)
	}

	dbFile.IsUploading = true
	dbFile.IsHashing = false
	dbFile.FileHash = hashString
	dbFile.SeederCount = 1
	dbInsertFile(*dbFile)

	//Add to payload
	AddToSeedString(*dbFile)
	pushNotification("Now seeding", dbFile.FileName)
}

func restartDownload(Hash string) {
	defer RecoverAndLog()

	file, err := dbGetFile(Hash)
	if err != nil {
		pushError("Error on restart download", err.Error())
		return
	}

	//If its not downloading we do not have to do anything
	//if !file.IsDownloading {
	//	return
	//}

	if file.IsPaused == true {
		return
	}

	//TODO: Seed discovery?

	//Early out if we have no seeder
	//if len(file.Seeders) == 0 {
	//	return
	//}

	//Get missing chunk indices
	var missingChunks []int
	for i := 0; i < file.NumChunks; i++ {
		if bitmap.Get(file.ChunkMap, i) == false {
			missingChunks = append(missingChunks, i)
		}
	}

	//Nothing more to download
	if len(missingChunks) == 0 {
		//TODO: set flag so we dont keep doing this?
		return
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(missingChunks), func(i, j int) { missingChunks[i], missingChunks[j] = missingChunks[j], missingChunks[i] })

	log.Println("Restarting Download Creation Session for", file.FileName)

	downloadSessions := []*Session{}
	// Create  sessions
	for i := 0; i < len(file.Seeders); i++ {
		surgeSession, err := createSession(file, file.Seeders[i])
		if err != nil {
			log.Println("Restarting Download Failed for", file.FileName, file.Seeders[i])
			continue
		}

		//Prime the session with known bytes downloaded
		surgeSession.Downloaded = int64(file.NumChunks-len(missingChunks)) * ChunkSize
		//If the last chunk is set, we want to deduct the missing bytes because its not a complete chunk
		lastChunkSet := bitmap.Get(file.ChunkMap, file.NumChunks-1)
		if lastChunkSet {
			overshotBytes := int64(file.NumChunks)*int64(ChunkSize) - file.FileSize
			surgeSession.Downloaded -= overshotBytes
		}

		go initiateSession(surgeSession)

		downloadSessions = append(downloadSessions, surgeSession)
	}

	if len(downloadSessions) == 0 {
		pushNotification("Restart download Session Failed, failed to connect to all seeders.", file.FileName)
		return
	}

	log.Println("Restarting Download for", file.FileName)
	log.Println("Total Chunks", file.NumChunks)
	log.Println("Remaining Chunks", len(missingChunks))

	seederAlternator := 0
	mutateSeederLock := sync.Mutex{}

	appendChunkLock := sync.Mutex{}
	chunksRemaining := len(missingChunks)

	downloadJob := func(terminateFlag *bool) {
		//Used to terminate the rescanning of peers
		terminate := func(flag *bool) {
			*flag = true
		}
		defer terminate(terminateFlag)

		for i := 0; i < chunksRemaining; i++ {
			newFileData := getListedFileByHash(Hash)
			if newFileData != nil {
				file = newFileData
			}

			for len(downloadSessions) == 0 {
				time.Sleep(time.Second * 5)
			}

			//Pause if file is paused
			dbFile, err := dbGetFile(file.FileHash)
			for err == nil && dbFile.IsPaused {
				time.Sleep(time.Second * 5)
				dbFile, err = dbGetFile(file.FileHash)
				if err != nil {
					break
				}
			}

			workerCount++

			//Create a async job to download a chunk
			requestChunkJob := func(chunkID int) {
				defer RecoverAndLog()

				success := false
				downloadSeeder := &Session{}
				downloadSeederAddr := ""

				if len(downloadSessions) > seederAlternator {
					//Get seeder
					downloadSeeder = downloadSessions[seederAlternator]
					if downloadSeeder != nil && downloadSeeder.session != nil {
						downloadSeederAddr = downloadSeeder.session.RemoteAddr().String()
						success = RequestChunk(downloadSeeder, file.FileHash, int32(chunkID))
					}
				}

				//if download fails append the chunk to remaining to retry later
				if !success {
					appendChunkLock.Lock()
					missingChunks = append(missingChunks, chunkID)
					chunksRemaining++
					appendChunkLock.Unlock()

					workerCount--

					//Drop the seeder
					mutateSeederLock.Lock()
					downloadSessions = removeAndCloseSessionOrdered(downloadSessions, downloadSeederAddr)
					log.Println("Lost connection", "Dropping 1 Session for Download "+file.FileName)
					mutateSeederLock.Unlock()

					if len(downloadSessions) == 0 {
						pushNotification("Download stopped, no more remote connections.", file.FileName)
						return
					}
				}

				//If chunk is requested add to transit map
				chunkKey := file.FileHash + "_" + strconv.Itoa(chunkID)

				chunkInTransitLock.Lock()
				chunksInTransit[chunkKey] = true
				chunkInTransitLock.Unlock()

				//Sleep for 30 seconds, check if entry still exists in transit map.
				time.Sleep(time.Second * 60)
				inTransit := chunksInTransit[chunkKey]

				//If its still in transit abort
				if inTransit {
					appendChunkLock.Lock()
					missingChunks = append(missingChunks, chunkID)
					chunksRemaining++
					appendChunkLock.Unlock()

					workerCount--

					//Drop the seeder
					mutateSeederLock.Lock()
					downloadSessions = removeAndCloseSessionOrdered(downloadSessions, downloadSeederAddr)
					log.Println("Lost connection", "Dropping 1 Session for Download "+file.FileName)
					mutateSeederLock.Unlock()

					if len(downloadSessions) == 0 {
						pushNotification("Download stopped, no more remote connections.", file.FileName)
						return
					}
				}
			}

			//get chunk id
			appendChunkLock.Lock()
			chunkid := missingChunks[i]
			appendChunkLock.Unlock()

			go requestChunkJob(chunkid)

			mutateSeederLock.Lock()
			seederAlternator++
			if seederAlternator > len(downloadSessions)-1 {
				seederAlternator = 0
			}
			mutateSeederLock.Unlock()

			for workerCount >= NumWorkers {
				time.Sleep(time.Millisecond)
				//log.Println("Active Workers:", workerCount)
				//fmt.Println("Active Workers:", workerCount)
			}
		}
	}

	scanForSeeders := func(terminateFlag *bool) {
		//While we are not terminated scan for new peers
		for *terminateFlag == false {
			time.Sleep(time.Second * 5)

			newFile := getListedFileByHash(Hash)
			if newFile != nil {
				//Check for new sessions
				for i := 0; i < len(newFile.Seeders); i++ {
					//Check if the newFile seeder is not already part of the downloadSessions
					alreadyAdded := false
					for j := 0; j < len(downloadSessions); j++ {
						if downloadSessions[j].session == nil {
							continue
						}
						if downloadSessions[j].session.RemoteAddr().String() == newFile.Seeders[i] {
							alreadyAdded = true
							break
						}
					}

					//Skip this entry
					if alreadyAdded {
						continue
					}

					surgeSession, err := createSession(newFile, newFile.Seeders[i])
					if err != nil {
						log.Println("Could not create session for download", Hash, newFile.Seeders[i])
						continue
					}

					dbFile, err := dbGetFile(Hash)
					if err != nil && dbFile != nil {
						//Prime the session with known bytes downloaded
						surgeSession.Downloaded = int64(dbFile.NumChunks-len(missingChunks)) * ChunkSize
						//If the last chunk is set, we want to deduct the missing bytes because its not a complete chunk
						lastChunkSet := bitmap.Get(dbFile.ChunkMap, dbFile.NumChunks-1)
						if lastChunkSet {
							overshotBytes := int64(dbFile.NumChunks)*int64(ChunkSize) - dbFile.FileSize
							surgeSession.Downloaded -= overshotBytes
						}

						go initiateSession(surgeSession)

						mutateSeederLock.Lock()
						downloadSessions = append(downloadSessions, surgeSession)
						mutateSeederLock.Unlock()
					}
				}
			}
		}
	}

	terminateFlag := false
	go downloadJob(&terminateFlag)
	go scanForSeeders(&terminateFlag)
}

func removeAndCloseSessionOrdered(slice []*Session, s string) []*Session {

	//Find empty session
	emptyIndex := -1
	for i := 0; i < len(slice); i++ {
		if slice[i] == nil || slice[i].session == nil {
			emptyIndex = i
			break
		}
	}
	if emptyIndex != -1 {
		slice = append(slice[:emptyIndex], slice[emptyIndex+1:]...)
	}

	//Find target session
	index := -1
	for i := 0; i < len(slice); i++ {
		//Find the target s
		if slice[i] != nil && slice[i].session != nil {
			if slice[i].session.RemoteAddr().String() == s {
				index = i
				break
			}
		}
	}
	if index != -1 {
		closeSession(slice[index])
		return append(slice[:index], slice[index+1:]...)
	}
	return slice
}

func createSession(File *File, Seeder string) (*Session, error) {
	defer RecoverAndLog()

	//Check if nkn session exists with address
	var downloadSession *Session
	/*for i := 0; i < len(Sessions); i++ {
		if Sessions[i].session.RemoteAddr().String() == File.Seeder {
			downloadSession = Sessions[i]
			break
		}
	}*/

	//There is no nkn session with this client yet, create a new session
	if downloadSession == nil {
		sessionConfing := nkn.GetDefaultSessionConfig()
		sessionConfing.MTU = 16384
		dialConfig := &nkn.DialConfig{
			SessionConfig: sessionConfing,
			DialTimeout:   60000,
		}

		downloadSession, err := client.DialWithConfig(Seeder, dialConfig)
		if err != nil {
			log.Println("Download Session timout for", File.FileName)
			return nil, err
		}
		log.Println("Download Session created for: ", File.FileName)
		downloadReader := bufio.NewReader(downloadSession)

		return &Session{
			reader:   downloadReader,
			session:  downloadSession,
			FileSize: File.FileSize,
			FileHash: File.FileHash,
		}, nil
	}

	//nkn session already exists create a new file session and include existing nkn session
	return &Session{
		reader:   downloadSession.reader,
		session:  downloadSession.session,
		FileSize: File.FileSize,
		FileHash: File.FileHash,
	}, nil
}

// FileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

//RecoverAndLog Recovers and then logs the stack
func RecoverAndLog() {
	if r := recover(); r != nil {
		fmt.Println("Panic digested from ", r)

		log.Printf("Internal error: %v", r)
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, true)
		log.Printf("%s\n", string(buf[0:stackSize]))

		var dir = GetSurgeDir()
		var logPathOS = dir + string(os.PathSeparator) + "paniclog.txt"
		f, _ := os.OpenFile(logPathOS, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		w := bufio.NewWriter(f)
		w.WriteString(string(buf[0:stackSize]))
		w.Flush()

		pushError("Panic", "Please check your log file and paniclog for more info")

		panic("Panic dumped but not digested, please check your log")
	}
}
