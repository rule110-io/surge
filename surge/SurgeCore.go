package surge

import (
	"bufio"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/hex"
	"io"
	"math/rand"
	"os"
	"path/filepath"
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
	fileInfo, err := dbGetFile(Hash)
	if err != nil {
		pushError("Error on open file", err.Error())
		return
	}
	OpenOSPath(fileInfo.Path)
}

//OpenFolderByHash opens the folder containing the file by hash in os
func OpenFolderByHash(Hash string) {
	fileInfo, err := dbGetFile(Hash)
	if err != nil {
		pushError("Error on open folder", err.Error())
		return
	}
	OpenOSPath(filepath.Dir(fileInfo.Path))
}

// RequestChunk sends a request to an address for a specific chunk of a specific file
func RequestChunk(Session *Session, FileID string, ChunkID int32) {
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
			log.Fatalln("Failed to send Surge Ruquest:", err)
		}
	}
}

// TransmitChunk tranmits target file chunk to address
func TransmitChunk(Session *Session, FileID string, ChunkID int32) {
	//Open file
	fileInfo, err := dbGetFile(FileID)
	if err != nil {
		pushError("Error on transmit chunk", err.Error())
		return
	}

	file, err := os.Open(fileInfo.Path)

	if err != nil {
		pushError("Error on transmit chunk", err.Error())
		return
	}
	defer file.Close()

	//Read the requested chunk
	chunkOffset := int64(ChunkID) * ChunkSize
	buffer := make([]byte, ChunkSize)
	bytesread, err := file.ReadAt(buffer, chunkOffset)

	if err != nil {
		if err != io.EOF {
			pushError("Error on transmit chunk", err.Error())
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
		pushError("Error on transmit chunk", err.Error())
		return
	}

	//Transmit the chunk
	err = SessionWrite(Session, dateReplySerialized, surgeChunkID) //Client.Send(nkn.NewStringArray(Addr), dateReplySerialized, nil)
	if err != nil {
		pushError("Error on transmit chunk", err.Error())
		return
	}
	log.Println("Chunk transmitted: ", bytesread, " bytes")
	Session.Uploaded += int64(bytesread)
}

// SendQueryRequest sends a query to a client on session
func SendQueryRequest(Addr string, Query string) {

	var surgeSession *Session = nil

	//Check for sessions
	for i := 0; i < len(Sessions); i++ {
		if Sessions[i].session.RemoteAddr().String() == Addr {
			surgeSession = Sessions[i]
			break
		}
	}

	if surgeSession == nil {
		//Create session
		sessionConfig := nkn.GetDefaultSessionConfig()
		sessionConfig.MTU = 16384
		sessionConfig.CheckTimeoutInterval = 1
		sessionConfig.InitialRetransmissionTimeout = 1
		sessionConfig.MaxRetransmissionTimeout = 1

		dialConfig := &nkn.DialConfig{
			SessionConfig: sessionConfig,
			DialTimeout:   5000,
		}

		downloadSession, err := client.DialWithConfig(Addr, dialConfig)
		if err != nil {
			log.Printf("Peer with address %s is not online, stopped trying after 5000ms\n", Addr)
			return
		}
		log.Printf("Connected to peer %s requesting file listings\n", Addr)

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
	} else {
		err := SessionWrite(surgeSession, msgSerialized, surgeQueryRequestID) //Client.Send(nkn.NewStringArray(Addr), msgSerialized, nil)
		if err != nil {
			log.Fatalln("Failed to send Surge Ruquest:", err)
		}
	}
}

// SendQueryResponse sends a query to a client on session
func SendQueryResponse(Session *Session, Query string) {
	b := []byte(queryPayload)
	err := SessionWrite(Session, b, surgeQueryResponseID) //Client.Send(nkn.NewStringArray(Addr), msgSerialized, nil)
	if err != nil {
		log.Fatalln("Failed to send Surge Ruquest:", err)
	}
}

// AllocateFile Allocates a file on disk at path with size in bytes
func AllocateFile(path string, size int64) {
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
		go WriteChunk(Session, surgeMessage.FileID, surgeMessage.ChunkID, surgeMessage.Data)
	}
}

func processQueryRequest(Session *Session, Data []byte) {
	//Try to parse SurgeMessage
	surgeQuery := &pb.SurgeQuery{}
	if err := proto.Unmarshal(Data, surgeQuery); err != nil {
		log.Fatalln("Failed to parse surge message:", err)
	}
	log.Println("Query received", surgeQuery.Query)

	SendQueryResponse(Session, surgeQuery.Query)
}

func processQueryResponse(Session *Session, Data []byte) {
	//Try to parse SurgeMessage
	s := string(Data)
	seeder := Session.session.RemoteAddr().String()

	clientOnlineMapLock.Lock()
	clientOnlineMap[seeder] = true
	clientOnlineMapLock.Unlock()

	//Remove exisiting file seed listings for this user
	n := 0
	for _, x := range ListedFiles {
		nn := 0
		for _, y := range x.Seeders {
			if y != seeder {
				x.Seeders[nn] = y
				nn++
			}
			x.Seeders = x.Seeders[:nn]
		}

		//If there are no remaining seeders, remove listing.
		if len(ListedFiles[n].Seeders) == 0 {
			ListedFiles[n] = x
			n++
		}
	}
	ListedFiles = ListedFiles[:n]

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
			FileName:  data[2],
			FileSize:  fileSize,
			FileHash:  data[4],
			Seeders:   []string{seeder},
			Path:      "",
			NumChunks: numChunks,
			ChunkMap:  nil,
		}

		//Replace existing, or remove.
		var replace = false
		for l := 0; l < len(ListedFiles); l++ {
			if ListedFiles[l].FileHash == newListing.FileHash {
				//if the seeder is unique add it as an additional seeder for the file
				ListedFiles[l].Seeders = append(ListedFiles[l].Seeders, seeder)
				replace = true
				break
			}
		}
		//Unique listing so we add
		if replace == false {
			ListedFiles = append(ListedFiles, newListing)
		}

		log.Println("Query response new file: ", newListing.FileName, " seeder: ", seeder)

		//Test gui
		//newButton := widget.NewButton(newListing.Filename+" | "+ByteCountSI(newListing.FileSize), func() {
		//	downloadFile(newListing.Seeder, newListing.FileSize, newListing.Filename)
		//})
		//fileBox.Append(newButton)
	}
}

//ParsePayloadString parses payload of files
func ParsePayloadString(s string) {
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
	Session.Downloaded += int64(bytesWritten)

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
	fi, err := os.Stat(path)
	if err != nil {
		log.Panicln("Error on get filesize", err)
	}
	// get the size
	return fi.Size()
}

//BuildSeedString builds a string of seeded files to share with clients
func BuildSeedString(dbFiles []File) {
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
	//Add to payload
	payload := surgeGenerateTopicPayload(dbFile.FileName, dbFile.FileSize, dbFile.FileHash)
	//log.Println(payload)
	queryPayload += payload
}

//SeedFile generates everything needed to seed a file
func SeedFile(Path string) bool {
	log.Println("Seeding file", Path)

	hashString, err := HashFile(Path)
	if err != nil {
		log.Println(err)
		pushNotification("Seed failed", "Could not hash file at "+Path)
		return false
	}

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
		FileName:    fileName,
		FileSize:    fileSize,
		FileHash:    hashString,
		Path:        Path,
		NumChunks:   numChunks,
		ChunkMap:    chunkMap,
		IsUploading: true,
	}

	//Check if file is already seeded
	_, err = dbGetFile(localFile.FileHash)
	if err == nil {
		//File already seeding
		pushNotification("Seed failed", fileName+" already seeding.")
		return false
	}

	//When seeding a new file enter file into db
	dbInsertFile(localFile)

	//Add to payload
	payload := surgeGenerateTopicPayload(fileName, fileSize, hashString)
	queryPayload += payload
	pushNotification("Now seeding", fileName)
	return true
}

func restartDownload(Hash string) {
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
	if len(file.Seeders) == 0 {
		return
	}

	//Get missing chunk indices
	var missingChunks []int32
	for i := 0; i < file.NumChunks; i++ {
		if bitmap.Get(file.ChunkMap, i) == false {
			missingChunks = append(missingChunks, int32(i))
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

	//Create a session
	surgeSession, err := createSession(file)
	if err != nil {
		log.Println("Restarting Download Failed for", file.FileName)
		return
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

	log.Println("Restarting Download for", file.FileName)
	log.Println("Total Chunks", file.NumChunks)
	log.Println("Remaining Chunks", len(missingChunks))

	//Download missing chunks
	for i := 0; i < len(missingChunks); i++ {
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
		go RequestChunk(surgeSession, file.FileHash, missingChunks[i])

		for workerCount >= NumWorkers {
			time.Sleep(time.Millisecond)
		}
	}
}

func createSession(File *File) (*Session, error) {
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

		downloadSession, err := client.DialWithConfig(File.Seeders[0], dialConfig)
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
