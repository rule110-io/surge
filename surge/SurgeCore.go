package surge

import (
	"bufio"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	pb "github.com/rule110-io/surge-ui/payloads"
	nkn "github.com/nknorg/nkn-sdk-go"
	"google.golang.org/protobuf/proto"
)

//TestTopic only for testing
const TestTopic = "poctest"

const surgeChunkID byte = 0x001
const surgeQueryRequestID byte = 0x002
const surgeQueryResponseID byte = 0x003

var queryPayload = ""

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
	//fileLocalMutex.Lock()
	file, err := os.Open(localPath + "/" + FileID)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//Read the requested chunk
	chunkOffset := int64(ChunkID) * ChunkSize
	buffer := make([]byte, ChunkSize)
	bytesread, err := file.ReadAt(buffer, chunkOffset)
	//fileLocalMutex.Unlock()

	if err != nil {
		if err != io.EOF {
			log.Fatal(err)
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
		log.Fatalln("Failed to encode surge message:", err)
	}

	//Transmit the chunk
	err = SessionWrite(Session, dateReplySerialized, surgeChunkID) //Client.Send(nkn.NewStringArray(Addr), dateReplySerialized, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Chunk transmitted: ", bytesread, " bytes")
	Session.Uploaded += int64(bytesread)
}

// SendQueryRequest sends a query to a client on session
func SendQueryRequest(Addr string, Query string) {

	//Create session
	sessionConfing := nkn.GetDefaultSessionConfig()
	sessionConfing.MTU = 16384
	
	dialConfig := &nkn.DialConfig{
		SessionConfig: sessionConfing,
		DialTimeout: 5000,
	}

	downloadSession, err := client.DialWithConfig(Addr, dialConfig)
	if err != nil {
		log.Printf("Peer with address %s is not online, stopped trying after 5000ms\n", Addr)
		return
	}
	log.Printf("Connected to peer %s requesting file listings\n", Addr)

	downloadReader := bufio.NewReader(downloadSession)

	surgeSession := &Session{
		reader:  downloadReader,
		session: downloadSession,
	}
	go initiateSession(surgeSession)

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
			log.Panic(err)
		}
		listenReader := bufio.NewReader(listenSession)

		surgeSession := &Session{
			reader:  listenReader,
			session: listenSession,
		}
		go initiateSession(surgeSession)

		time.Sleep(time.Millisecond)
	}
}

// Listen will listen to incoming requests for chunks
func Listen() {
	go listenForSession()
}

func initiateSession(Session *Session) {
	Sessions = append(Sessions, Session)

	for true {
		data, chunkType, err := SessionRead(Session)
		if err != nil {
			if err.Error() == "session closed" {
				break
			}
			log.Fatal(err)
		}

		log.Println(Session)

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

	//find index in Sessions
	var index = -1
	for i := 0; i < len(Sessions); i++ {
		if(Sessions[i] == Session) {
			index = i
			break
		}
	}
	if index == -1 {
		log.Panic("Couldnt find session to remove")
	}
	//Remove session
	Sessions[index] = Sessions[len(Sessions)-1]
	Sessions[len(Sessions)-1] = nil
	Sessions = Sessions[:len(Sessions)-1]

	log.Println("-=Session closed=-")
}

func processChunk(Session *Session, Data []byte) {
	//Try to parse SurgeMessage
	surgeMessage := &pb.SurgeMessage{}
	if err := proto.Unmarshal(Data, surgeMessage); err != nil {
		log.Fatalln("Failed to parse surge message:", err)
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
	log.Println(surgeQuery.Query)

	SendQueryResponse(Session, surgeQuery.Query)
}

func processQueryResponse(Session *Session, Data []byte) {
	//Try to parse SurgeMessage
	s := string(Data)
	log.Println("Query Reponse: ", s)

	//Parse the response
	payloadSplit := strings.Split(s, "surge://")
	for j := 0; j < len(payloadSplit); j++ {
		data := strings.Split(payloadSplit[j], "|")

		if len(data) < 3 {
			continue
		}

		fileSize, _ := strconv.ParseInt(data[3], 10, 64)

		newListing := File{data[2], fileSize, data[4], Session.session.RemoteAddr().String()}

		ListedFiles = append(ListedFiles, newListing)

		//Test gui
		//newButton := widget.NewButton(newListing.Filename+" | "+ByteCountSI(newListing.FileSize), func() {
		//	downloadFile(newListing.Seeder, newListing.FileSize, newListing.Filename)
		//})
		//fileBox.Append(newButton)
	}
}

//WriteChunk writs a chunk to disk
func WriteChunk(Session *Session, FileID string, ChunkID int32, Chunk []byte) {
	var path = remotePath + "/" + FileID

	_, err := os.Stat(path)

	//Open file
	//fileRemoteMutex.Lock()
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	chunkOffset := int64(ChunkID) * ChunkSize

	bytesWritten, err := file.WriteAt(Chunk, chunkOffset)
	//fileRemoteMutex.Unlock()
	if err != nil {
		log.Fatal(err)
		return
	}
	//Success
	log.Println("Chunk written to disk: ", bytesWritten, " bytes")
	Session.Downloaded += int64(bytesWritten)

	workerCount--
	chunksReceived++
}

//TopicEncode .
func TopicEncode(topic string) string {
	return "SRG_" + strings.ReplaceAll(b64.StdEncoding.EncodeToString([]byte(topic)), "=", "-")
}

func surgeGenerateTopicPayload(fileName string, sizeInBytes int64, hash string) string {
	//Example payload
	//surge://|file|The_Two_Towers-The_Purist_Edit-Trailer.avi|14997504|965c013e991ee246d63d45ea71954c4d|/

	//Append to local files
	localFile := File{fileName, sizeInBytes, hash, "local"}
	LocalFiles = append(LocalFiles, localFile)

	return "surge://|file|" + fileName + "|" + strconv.FormatInt(sizeInBytes, 10) + "|" + hash + "|/"
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
		log.Fatal(err)
	}
	// get the size
	return fi.Size()
}

//ScanLocal scans local files
func ScanLocal() {
	var files []string

	root := localPath
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if path != root {
			name := filepath.Base(path)
			if len(strings.Split(name, ".")[0]) > 0 {
				files = append(files, path)
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	var topic = TestTopic
	queryPayload = ""
	for _, file := range files {
		fmt.Println(file)

		hashString, err := HashFile(file)
		if err != nil {
			log.Panicln(err)
			continue
		}


		payload := surgeGenerateTopicPayload(filepath.Base(file), surgeGetFileSize(file), hashString)

		queryPayload += payload

	}
	topicEncoded := TopicEncode(topic)
	sendSeedSubscription(topicEncoded, queryPayload)
	log.Println("Seeding to Topic: ", topicEncoded)
}
