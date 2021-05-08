package surge

import (
	"fmt"

	"github.com/rule110-io/surge/backend/messaging"
)

const (
	MessageIDRequestFiles = iota
	MessageIDReplyFiles
)

func MessageReceived(msg *messaging.MessageReceivedObj) {
	fmt.Println(string(msg.Data))

	switch msg.Type {
	case MessageIDRequestFiles:
		SendRequestFilesReply(msg)
		break
	case MessageIDReplyFiles:
		//process file data
		break
	}

}

func RequestFiles(topic string) {
	//Create the data object
	dataObj := messaging.MessageObj{
		Type:  MessageIDRequestFiles,
		Topic: topic,
	}

	messaging.Broadcast(&dataObj)
}

func SendRequestFilesReply(msg *messaging.MessageReceivedObj) {
	//Create the data object
	dataObj := messaging.MessageObj{
		Type:  MessageIDReplyFiles,
		Topic: msg.Topic,
		Data:  []byte(queryPayload),
	}
	msg.Reply(&dataObj)
}
