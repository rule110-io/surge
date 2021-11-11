package surge

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/rule110-io/surge/backend/constants"
	"github.com/rule110-io/surge/backend/messaging"
	"github.com/rule110-io/surge/backend/models"
	"github.com/rule110-io/surge/backend/mutexes"
)

const (
	MessageIDAnnounceFiles = iota
	MessageIDAnnounceFilesReply
	MessageIDAnnounceNewFile
)

func MessageReceived(msg *messaging.MessageReceivedObj) {
	fmt.Println(string("\033[36m"), "MESSAGE RECEIVED", string(msg.Data))
	fmt.Println(msg.Data)

	switch msg.Type {
	case MessageIDAnnounceFiles:
		SendAnnounceFilesReply(msg)
		processQueryResponse(msg.Sender, msg.Data)
	case MessageIDAnnounceFilesReply:
		//process file data
		processQueryResponse(msg.Sender, msg.Data)
	case MessageIDAnnounceNewFile:
		//process file data
		processQueryResponse(msg.Sender, msg.Data)
	}

}

func AnnounceFiles(topic string) {
	fmt.Println(string("\033[36m"), "REQUESTING FILES FOR TOPIC", topic)
	//Create the data object
	dataObj := messaging.MessageObj{
		Type:  MessageIDAnnounceFiles,
		Topic: topic,
		Data:  []byte(queryPayload),
	}

	messaging.Broadcast(&dataObj)
}

func SendAnnounceFilesReply(msg *messaging.MessageReceivedObj) {
	fmt.Println(string("\033[36m"), "SENDING FILE REQUEST REPLY", msg.Topic, msg.Sender)
	//Create the data object
	dataObj := messaging.MessageObj{
		Type:  MessageIDAnnounceFilesReply,
		Topic: msg.Topic,
		Data:  []byte(queryPayload),
	}
	msg.Reply(&dataObj)
}

func AnnounceNewFile(file *models.File) {
	fmt.Println(string("\033[36m"), "ANNOUNCE NEW FILE FOR TOPIC", file.Topic)

	//Create payload
	payload := surgeGenerateTopicPayload(file.FileName, file.FileSize, file.FileHash, file.Topic)

	//Create the data object
	dataObj := messaging.MessageObj{
		Type:  MessageIDAnnounceNewFile,
		Topic: file.Topic,
		Data:  []byte(payload),
	}

	messaging.Broadcast(&dataObj)
}

func processQueryResponse(seeder string, Data []byte) {

	//Try to parse SurgeMessage
	s := string(Data)
	fmt.Println(string("\033[36m"), "file query response received", seeder, string("\033[0m"))

	mutexes.ListedFilesLock.Lock()

	//Parse the response
	payloadSplit := strings.Split(s, "surge://")
	for j := 0; j < len(payloadSplit); j++ {
		data := strings.Split(payloadSplit[j], "|")

		if len(data) < 3 {
			continue
		}

		fileSize, _ := strconv.ParseInt(data[3], 10, 64)
		numChunks := int((fileSize-1)/int64(constants.ChunkSize)) + 1

		newListing := models.File{
			FileLocation: "remote",
			FileName:     data[2],
			FileSize:     fileSize,
			FileHash:     data[4],
			Path:         "",
			NumChunks:    numChunks,
			ChunkMap:     nil,
			Topic:        data[5],
		}

		//Replace existing, or remove.
		var replace = false
		for l := 0; l < len(ListedFiles); l++ {
			if ListedFiles[l].FileHash == newListing.FileHash {
				replace = true
				break
			}
		}
		//Unique listing so we add
		if !replace {
			ListedFiles = append(ListedFiles, newListing)
		}

		//Check if we have this file in the db already, if so add this new seeder
		mutexes.FileWriteLock.Lock()
		dbFile, err := dbGetFile(newListing.FileHash)
		if err == nil {
			log.Println("File in query was already known, seeder added!")
			dbInsertFile(*dbFile)
		}
		mutexes.FileWriteLock.Unlock()

		//We now add this seeder to our file seeders
		AddFileSeeder(dbFile.FileHash, seeder)

		fmt.Println(string("\033[33m"), "Filename", newListing.FileName, "FileHash", newListing.FileHash, string("\033[0m"))

		log.Println("Query response new file: ", newListing.FileName, " seeder: ", seeder)
	}
	mutexes.ListedFilesLock.Unlock()
}
