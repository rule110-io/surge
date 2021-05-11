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
	MessageIDRequestFiles = iota
	MessageIDReplyFiles
)

func MessageReceived(msg *messaging.MessageReceivedObj) {
	fmt.Println(string("\033[36m"), "MESSAGE RECEIVED", string(msg.Data))
	fmt.Println(msg.Data)

	switch msg.Type {
	case MessageIDRequestFiles:
		SendRequestFilesReply(msg)
		break
	case MessageIDReplyFiles:
		//process file data
		processQueryResponse(msg.Sender, msg.Data)
		break
	}

}

func RequestFiles(topic string) {
	fmt.Println(string("\033[36m"), "REQUESTING FILES FOR TOPIC", topic)
	//Create the data object
	dataObj := messaging.MessageObj{
		Type:  MessageIDRequestFiles,
		Topic: topic,
	}

	messaging.Broadcast(&dataObj)
}

func SendRequestFilesReply(msg *messaging.MessageReceivedObj) {
	fmt.Println(string("\033[36m"), "SENDING FILE REQUEST REPLY", msg.Topic, msg.Sender)
	//Create the data object
	dataObj := messaging.MessageObj{
		Type:  MessageIDReplyFiles,
		Topic: msg.Topic,
		Data:  []byte(queryPayload),
	}
	msg.Reply(&dataObj)
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
			Seeders:      []string{seeder},
			Path:         "",
			NumChunks:    numChunks,
			ChunkMap:     nil,
			SeederCount:  1,
			Topic:        data[5],
		}

		//Replace existing, or remove.
		var replace = false
		for l := 0; l < len(ListedFiles); l++ {
			if ListedFiles[l].FileHash == newListing.FileHash {

				//if the seeder is unique add it as an additional seeder for the file
				ListedFiles[l].Seeders = append(ListedFiles[l].Seeders, seeder)
				ListedFiles[l].Seeders = distinctStringSlice(ListedFiles[l].Seeders)
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
	}
	mutexes.ListedFilesLock.Unlock()
}
