package main

import (
	"bufio"
	"crypto/md5"
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

	pb "nsurge/proto/SurgeMessage"
	pbq "nsurge/proto/SurgeQuery"
	nkn "github.com/nknorg/nkn-sdk-go"
	"google.golang.org/protobuf/proto"
)

const testTopic = "poctest"

const surgeChunkID byte = 0x001
const surgeQueryRequestID byte = 0x002
const surgeQueryResponseID byte = 0x003

var queryPayload = ""

// SurgeRequestChunk sends a request to an address for a specific chunk of a specific file
func SurgeRequestChunk(Session SurgeSession, FileID string, ChunkID int32) {
	msg := &pb.SurgeMessage{
		FileID:  FileID,
		ChunkID: ChunkID,
	}
	msgSerialized, err := proto.Marshal(msg)
	if err != nil {
		log.Fatalln("Failed to encode surge message:", err)
	} else {
		err := surgeSessionWrite(Session, msgSerialized, surgeChunkID) //Client.Send(nkn.NewStringArray(Addr), msgSerialized, nil)
		if err != nil {
			log.Fatalln("Failed to send Surge Ruquest:", err)
		}
	}
}

// SurgeTransmitChunk tranmits target file chunk to address
func SurgeTransmitChunk(Session SurgeSession, FileID string, ChunkID int32) {

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
	err = surgeSessionWrite(Session, dateReplySerialized, surgeChunkID) //Client.Send(nkn.NewStringArray(Addr), dateReplySerialized, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Chunk transmitted: ", bytesread, " bytes")
}

// SurgeSendQueryRequest sends a query to a client on session
func SurgeSendQueryRequest(Addr string, Query string) {

	//Create session
	sessionConfing := nkn.GetDefaultSessionConfig()
	sessionConfing.MTU = 16384
	dialConfig := &nkn.DialConfig{
		SessionConfig: sessionConfing,
	}

	downloadSession, err := client.DialWithConfig(Addr, dialConfig)
	if err != nil {
		log.Fatal(err)
	}
	downloadReader := bufio.NewReader(downloadSession)

	surgeSession := SurgeSession{
		Reader:  downloadReader,
		Session: downloadSession,
	}
	go initiateSession(surgeSession)

	msg := &pbq.SurgeQuery{
		Query: Query,
	}
	msgSerialized, err := proto.Marshal(msg)
	if err != nil {
		log.Fatalln("Failed to encode surge message:", err)
	} else {
		err := surgeSessionWrite(surgeSession, msgSerialized, surgeQueryRequestID) //Client.Send(nkn.NewStringArray(Addr), msgSerialized, nil)
		if err != nil {
			log.Fatalln("Failed to send Surge Ruquest:", err)
		}
	}
}

// SurgeSendQueryResponse sends a query to a client on session
func SurgeSendQueryResponse(Session SurgeSession, Query string) {
	b := []byte(queryPayload)
	err := surgeSessionWrite(Session, b, surgeQueryResponseID) //Client.Send(nkn.NewStringArray(Addr), msgSerialized, nil)
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

		surgeSession := SurgeSession{
			Reader:  listenReader,
			Session: listenSession,
		}
		go initiateSession(surgeSession)

		time.Sleep(time.Millisecond)
	}
}

// Listen will listen to incoming requests for chunks
func Listen() {
	go listenForSession()

	//Listen as long as client is alive
	/*for SurgeActive {
		for _, session := range sessions {

			//TODO: deal with sessions closing
			if testSession == nil {
				time.Sleep(time.Millisecond * 100)
				continue
			}

			//msg := <-client.OnMessage.C

		}
	}*/
}

func initiateSession(Session SurgeSession) {
	sessions = append(sessions, Session)

	for true {
		data, chunkType, err := surgeSessionRead(Session)
		if err != nil {
			if err.Error() == "session closed" {
				break
			}
			log.Fatal(err)
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
}

func processChunk(Session SurgeSession, Data []byte) {
	//Try to parse SurgeMessage
	surgeMessage := &pb.SurgeMessage{}
	if err := proto.Unmarshal(Data, surgeMessage); err != nil {
		log.Fatalln("Failed to parse surge message:", err)
	}

	//Data nill means its a request for data
	if surgeMessage.Data == nil {
		go SurgeTransmitChunk(Session, surgeMessage.FileID, surgeMessage.ChunkID)
	} else { //If data is not nill we are receiving data
		go surgeWriteChunk(surgeMessage.FileID, surgeMessage.ChunkID, surgeMessage.Data)
	}
}

func processQueryRequest(Session SurgeSession, Data []byte) {
	//Try to parse SurgeMessage
	surgeQuery := &pbq.SurgeQuery{}
	if err := proto.Unmarshal(Data, surgeQuery); err != nil {
		log.Fatalln("Failed to parse surge message:", err)
	}
	log.Println(surgeQuery.Query)

	SurgeSendQueryResponse(Session, surgeQuery.Query)
}

func processQueryResponse(Session SurgeSession, Data []byte) {
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

		newListing := SurgeFile{data[2], fileSize, data[4], Session.Session.RemoteAddr().String()}

		listedFiles = append(listedFiles, newListing)

		//Test gui
		//newButton := widget.NewButton(newListing.Filename+" | "+ByteCountSI(newListing.FileSize), func() {
		//	downloadFile(newListing.Seeder, newListing.FileSize, newListing.Filename)
		//})
		//fileBox.Append(newButton)
	}
}

func surgeWriteChunk(FileID string, ChunkID int32, Chunk []byte) {
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
	workerCount--
	chunksReceived++
}

func surgeTopicEncode(topic string) string {
	return "SRG_" + strings.ReplaceAll(b64.StdEncoding.EncodeToString([]byte(topic)), "=", "-")
}

func surgeGenerateTopicPayload(fileName string, sizeInBytes int64, md5 string) string {
	//Example payload
	//surge://|file|The_Two_Towers-The_Purist_Edit-Trailer.avi|14997504|965c013e991ee246d63d45ea71954c4d|/
	return "surge://|file|" + fileName + "|" + strconv.FormatInt(sizeInBytes, 10) + "|" + md5 + "|/"
}

// SurgeHashFile generates hash for file given filepath
func SurgeHashFile(filePath string) (string, error) {
	//Initialize variable returnMD5String now in case an error has to be returned
	var returnMD5String string

	//Open the passed argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}

	//Tell the program to call the following function when the current function returns
	defer file.Close()

	//Open a new hash interface to write to
	hash := md5.New()

	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}

	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]

	//Convert the bytes to a string
	returnMD5String = hex.EncodeToString(hashInBytes)

	return returnMD5String, nil

}

func surgeGetFileSize(path string) int64 {
	fi, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	// get the size
	return fi.Size()
}

func surgeScanLocal() {
	var files []string

	root := localPath
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if path != root {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	var topic = testTopic
	queryPayload = ""
	for _, file := range files {
		fmt.Println(file)

		md5, err := SurgeHashFile(file)
		if err != nil {
			log.Panicln(err)
			continue
		}
		payload := surgeGenerateTopicPayload(filepath.Base(file), surgeGetFileSize(file), md5)

		queryPayload += payload

	}
	topicEncoded := surgeTopicEncode(topic)
	sendSeedSubscription(topicEncoded, queryPayload)
	log.Println("Seeding to Topic: ", topicEncoded)
}
