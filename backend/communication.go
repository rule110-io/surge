package surge

import (
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
	MessageIDAnnounceRemoveFile
)

func MessageReceived(msg *messaging.MessageReceivedObj) {
	switch msg.Type {
	case MessageIDAnnounceFiles:
		if msg.Sender != GetAccountAddress() {
			go SendAnnounceFilesReply(msg)
		}
		go processQueryResponse(msg.Sender, msg.Data)
	case MessageIDAnnounceFilesReply:
		go processQueryResponse(msg.Sender, msg.Data)
	case MessageIDAnnounceNewFile:
		go processQueryResponse(msg.Sender, msg.Data)
	case MessageIDAnnounceRemoveFile:
		go processRemoveFile(string(msg.Data), msg.Sender)
	}

}

func AnnounceFiles(topicEncoded string) {
	payload := getTopicPayload(topicEncoded)

	dataObj := messaging.MessageObj{
		Type:         MessageIDAnnounceFiles,
		TopicEncoded: topicEncoded,
		Data:         []byte(payload),
	}

	messaging.Broadcast(&dataObj)
}

func SendAnnounceFilesReply(msg *messaging.MessageReceivedObj) {
	payload := getTopicPayload(msg.TopicEncoded)

	if len(payload) > 0 {
		//Create the data object
		dataObj := messaging.MessageObj{
			Type:         MessageIDAnnounceFilesReply,
			TopicEncoded: msg.TopicEncoded,
			Data:         []byte(payload),
		}
		msg.Reply(&dataObj)
	}
}

func AnnounceNewFile(file *models.File) {
	//Create payload
	payload := surgeGenerateTopicPayload(file.FileName, file.FileSize, file.FileHash, file.Topic)

	//Create the data object
	dataObj := messaging.MessageObj{
		Type:         MessageIDAnnounceNewFile,
		TopicEncoded: TopicEncode(file.Topic),
		Data:         []byte(payload),
	}

	messaging.Broadcast(&dataObj)
}

func AnnounceRemoveFile(topic string, fileHash string) {
	//Create the data object
	dataObj := messaging.MessageObj{
		Type:         MessageIDAnnounceRemoveFile,
		TopicEncoded: TopicEncode(topic),
		Data:         []byte(fileHash),
	}

	messaging.Broadcast(&dataObj)
}

func processRemoveFile(hash string, seeder string) {
	RemoveFileSeeder(hash, seeder)

	mutexes.ListedFilesLock.Lock()
	defer mutexes.ListedFilesLock.Unlock()

	//Remove empty seeders listings
	for i := 0; i < len(ListedFiles); i++ {
		if !AnySeeders(ListedFiles[i].FileHash) {
			// Remove the element at index i from a.
			ListedFiles[i] = ListedFiles[len(ListedFiles)-1] // Copy last element to index i.
			ListedFiles[len(ListedFiles)-1] = models.File{}  // Erase last element (write zero value).
			ListedFiles = ListedFiles[:len(ListedFiles)-1]   // Truncate slice.
			i--
		}
	}
}

func processQueryResponse(seeder string, Data []byte) {

	//Try to parse SurgeMessage
	s := string(Data)
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
			FileName:  data[2],
			FileSize:  fileSize,
			FileHash:  data[4],
			Path:      "",
			NumChunks: numChunks,
			ChunkMap:  nil,
			Topic:     data[5],
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

		//We now add this seeder to our file seeders
		AddFileSeeder(newListing.FileHash, seeder)
	}
	mutexes.ListedFilesLock.Unlock()
}

func getTopicPayload(topicEncoded string) string {
	dbFiles := dbGetAllFiles()

	payload := ""
	for _, dbFile := range dbFiles {

		if TopicEncode(dbFile.Topic) != topicEncoded {
			continue
		}

		if dbFile.IsUploading {
			//Add to payload
			payload += surgeGenerateTopicPayload(dbFile.FileName, dbFile.FileSize, dbFile.FileHash, dbFile.Topic)
		}
	}
	return payload
}
