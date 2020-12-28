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

	"log"

	bitmap "github.com/boljen/go-bitmap"
	pb "github.com/rule110-io/surge-ui/payloads"
	"github.com/rule110-io/surge-ui/surge/constants"
	"github.com/rule110-io/surge-ui/surge/platform"
	"github.com/rule110-io/surge-ui/surge/sessionmanager"
	"google.golang.org/protobuf/proto"

	open "github.com/skratchdot/open-golang/open"
)

//TestTopic only for testing
const TestTopic = "privateTest"

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
func RequestChunk(Session *sessionmanager.Session, FileID string, ChunkID int32) bool {
	if Session == nil || Session.Session == nil {
		return false
	}

	msg := &pb.SurgeMessage{
		FileID:  FileID,
		ChunkID: ChunkID,
	}
	msgSerialized, err := proto.Marshal(msg)
	if err != nil {
		log.Panic("Failed to encode surge message:", err)
	} else {
		fmt.Println(string("\033[31m"), "Request Chunk", FileID, ChunkID, string("\033[0m"))

		written, err := SessionWrite(Session, msgSerialized, surgeChunkID) //Client.Send(nkn.NewStringArray(Addr), msgSerialized, nil)
		if err != nil {
			log.Println("Failed to request chunk", err)
			return false
		}

		//Write add to upload
		bandwidthAccumulatorMapLock.Lock()
		uploadBandwidthAccumulator[FileID] += written
		bandwidthAccumulatorMapLock.Unlock()
	}

	return true
}

// TransmitChunk tranmits target file chunk to address
func TransmitChunk(Session *sessionmanager.Session, FileID string, ChunkID int32) {
	defer RecoverAndLog()

	//Open file

	fileWriteLock.Lock()
	fileInfo, err := dbGetFile(FileID)
	if err != nil {
		log.Println("Error on transmit chunk - file not in db", err.Error())
		return
	}
	fileInfo.ChunksShared++
	dbInsertFile(*fileInfo)
	fileWriteLock.Unlock()

	file, err := os.Open(fileInfo.Path)

	//When we have an OS read error on the file mark the file as missing, stop down and uploads on it
	if err != nil {
		log.Println("Error on transmit chunk - file read failure", err.Error())

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
			log.Println("Error on transmit chunk - read chunk failed: ", ChunkID, err.Error())
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
		log.Panic("Error on transmit chunk - chunk serialization error", err.Error())
		return
	}

	//Transmit the chunk
	fmt.Println(string("\033[31m"), "Transmit Chunk", FileID, ChunkID, string("\033[0m"))
	written, err := SessionWrite(Session, dateReplySerialized, surgeChunkID) //Client.Send(nkn.NewStringArray(Addr), dateReplySerialized, nil)
	if err != nil {
		log.Println("Error on transmit chunk - failed to write to session", err.Error())
		return
	}
	log.Println("Chunk transmitted: ", bytesread, " bytes")

	//Write add to upload
	bandwidthAccumulatorMapLock.Lock()
	uploadBandwidthAccumulator[FileID] += written
	bandwidthAccumulatorMapLock.Unlock()
}

// SendQueryRequest sends a query to a client on session
func SendQueryRequest(Addr string, Query string) bool {

	surgeSession, exists := sessionmanager.GetExistingSession(Addr, constants.SendQueryRequestSessionTimeout)

	if !exists {
		return false
	}

	msg := &pb.SurgeQuery{
		Query: Query,
	}
	msgSerialized, err := proto.Marshal(msg)
	if err != nil {
		log.Panic("Failed to encode surge message:", err)
		return false
	}

	fmt.Println(string("\033[31m"), "Send Query Request", Addr, string("\033[0m"))
	written, err := SessionWrite(surgeSession, msgSerialized, surgeQueryRequestID) //Client.Send(nkn.NewStringArray(Addr), msgSerialized, nil)
	if err != nil {
		log.Println("Failed to send Surge Request:", err)
		return false
	}

	//Write add to upload
	bandwidthAccumulatorMapLock.Lock()
	uploadBandwidthAccumulator["DISCOVERY"] += written
	bandwidthAccumulatorMapLock.Unlock()

	return true
}

// SendQueryResponse sends a query to a client on session
func SendQueryResponse(Session *sessionmanager.Session, Query string) {

	b := []byte(queryPayload)
	fmt.Println(string("\033[31m"), "Send Query Response", Session.Session.RemoteAddr().String(), string("\033[0m"))
	written, err := SessionWrite(Session, b, surgeQueryResponseID) //Client.Send(nkn.NewStringArray(Addr), msgSerialized, nil)
	if err != nil {
		log.Println("Failed to send Surge Ruquest:", err)
	}
	//Write add to upload
	bandwidthAccumulatorMapLock.Lock()
	uploadBandwidthAccumulator["DISCOVERY"] += written
	bandwidthAccumulatorMapLock.Unlock()
}

// AllocateFile Allocates a file on disk at path with size in bytes
func AllocateFile(path string, size int64) {

	fd, err := os.Create(path)
	if err != nil {
		log.Panic("Failed to create output")
	}
	_, err = fd.Seek(size-1, 0)
	if err != nil {
		log.Panic("Failed to seek")
	}
	_, err = fd.Write([]byte{0})
	if err != nil {
		log.Panic("Write failed")
	}
	err = fd.Close()
	if err != nil {
		log.Panic("Failed to close file")
	}
}

/*func closeSession(Session *sessionmanager.Session) {

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
	Session.Session.Close()
	Session.Session = nil
	Session.Reader = nil
	if Session.File != nil {
		err := Session.File.Close()
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
}*/

func processChunk(Session *sessionmanager.Session, Data []byte) {

	//Try to parse SurgeMessage
	surgeMessage := &pb.SurgeMessage{}
	if err := proto.Unmarshal(Data, surgeMessage); err != nil {
		log.Panic("Failed to parse surge message:", err)
	}
	fmt.Println(string("\033[31m"), "PROCESSING CHUNK", string("\033[0m"))

	//Write add to download
	bandwidthAccumulatorMapLock.Lock()
	downloadBandwidthAccumulator[surgeMessage.FileID] += len(Data)
	bandwidthAccumulatorMapLock.Unlock()

	//If this is the first file data over this session we need to set the session
	/*if Session.FileHash == "" {
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
	}*/

	//Data nill means its a request for data
	if surgeMessage.Data == nil {
		go TransmitChunk(Session, surgeMessage.FileID, surgeMessage.ChunkID)
	} else { //If data is not nill we are receiving data

		//When we receive a chunk mark it as no longer in transit
		chunkKey := surgeMessage.FileID + "_" + strconv.Itoa(int(surgeMessage.ChunkID))

		chunkInTransitLock.Lock()
		chunksInTransit[chunkKey] = false
		chunkInTransitLock.Unlock()

		go WriteChunk(surgeMessage.FileID, surgeMessage.ChunkID, surgeMessage.Data)
	}
}

func processQueryRequest(Session *sessionmanager.Session, Data []byte) {

	//Try to parse SurgeMessage
	surgeQuery := &pb.SurgeQuery{}
	if err := proto.Unmarshal(Data, surgeQuery); err != nil {
		log.Panic("Failed to parse surge message:", err)
	}
	log.Println("Query received", surgeQuery.Query)

	SendQueryResponse(Session, surgeQuery.Query)
}

func processQueryResponse(Session *sessionmanager.Session, Data []byte) {

	//Try to parse SurgeMessage
	s := string(Data)
	seeder := Session.Session.RemoteAddr().String()

	fmt.Println(string("\033[36m"), "file query response received", seeder, string("\033[0m"))

	ListedFilesLock.Lock()

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
			seeders:     []string{seeder},
			Path:        "",
			NumChunks:   numChunks,
			ChunkMap:    nil,
			seederCount: 1,
		}

		//Replace existing, or remove.
		var replace = false
		for l := 0; l < len(ListedFiles); l++ {
			if ListedFiles[l].FileHash == newListing.FileHash {

				//if the seeder is unique add it as an additional seeder for the file
				ListedFiles[l].seeders = append(ListedFiles[l].seeders, seeder)
				ListedFiles[l].seeders = distinctStringSlice(ListedFiles[l].seeders)
				ListedFiles[l].seederCount = len(ListedFiles[l].seeders)

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
			seeders:   seeder,
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
				ListedFiles[l].seeders = append(ListedFiles[l].seeders, seeder...)
				replace = true
				break
			}
		}
		//Unique listing so we add
		if replace == false {
			ListedFiles = append(ListedFiles, newListing)
		}
		ListedFilesLock.Unlock()

		log.Println("Program paramater new file: ", newListing.FileName, " seeder: ", newListing.seeders)

		go DownloadFile(newListing.FileHash)

		//Test gui
		//newButton := widget.NewButton(newListing.Filename+" | "+ByteCountSI(newListing.FileSize), func() {
		//	downloadFile(newListing.Seeder, newListing.FileSize, newListing.Filename)
		//})
		//fileBox.Append(newButton)
	}
}

//WriteChunk writes a chunk to disk
func WriteChunk(FileID string, ChunkID int32, Chunk []byte) {
	defer RecoverAndLog()

	workerCount--

	fileInfo, err := dbGetFile(FileID)
	if err != nil {
		log.Println("Error on write chunk (db get)", err.Error())
		return
	}
	remoteFolder, err := platform.GetRemoteFolder()
	if err != nil {
		pushError("Error on write chunk (GetRemoteFolder)", err.Error())
		return
	}
	path := remoteFolder + string(os.PathSeparator) + fileInfo.FileName
	osFile, err := os.OpenFile(path, os.O_RDWR, 0644)
	defer osFile.Close()

	if err != nil {
		pushError("Error on write chunk (sessionmanager.GetFile)", err.Error())
		return
	}

	chunkOffset := int64(ChunkID) * ChunkSize
	bytesWritten, err := osFile.WriteAt(Chunk, chunkOffset)
	if err != nil {
		pushError("Error on write chunk (file write)", err.Error())
		return
	}
	//Success
	log.Println("Chunk written to disk: ", bytesWritten, " bytes")

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
		log.Panic("Error on get filesize", err)
	}
	// get the size
	return fi.Size()
}

//BuildSeedString builds a string of seeded files to share with clients
func BuildSeedString(dbFiles []File) {

	newQueryPayload := ""
	for _, dbFile := range dbFiles {
		magnet := surgeGenerateMagnetLink(dbFile.FileName, dbFile.FileSize, dbFile.FileHash, client.Addr().String())
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

	//Make sure you're subscribed when seeding a file
	go subscribeToSurgeTopic()
}

//SeedFile generates everything needed to seed a file
func SeedFile(Path string) bool {

	log.Println("Seeding file", Path)

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Panic(err)
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
	dbInsertFile(*dbFile)

	//Add to payload
	AddToSeedString(*dbFile)
	pushNotification("Now seeding", dbFile.FileName)
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
		//log.Printf("%s\n", string(buf[0:stackSize]))

		var dir = platform.GetSurgeDir()
		var logPathOS = dir + string(os.PathSeparator) + "paniclog.txt"
		f, _ := os.OpenFile(logPathOS, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		w := bufio.NewWriter(f)
		w.WriteString(string(buf[0:stackSize]))
		w.Flush()

		pushError("Panic", "Please check your log file and paniclog for more info")

		panic("Panic dumped but not digested, please check your log")
	}
}
