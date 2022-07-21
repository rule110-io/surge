package messaging

import (
	"encoding/json"
	"log"

	nkn "github.com/nknorg/nkn-sdk-go"
)

var nknClient *nkn.MultiClient
var nknAccount *nkn.Account
var onMessageHandler func(*MessageReceivedObj)

//Initializes provides client with required nkn objects
func Initialize(client *nkn.MultiClient, account *nkn.Account, onMsgHandler func(*MessageReceivedObj)) {
	nknClient = client
	nknAccount = account
	onMessageHandler = onMsgHandler

	go listen()
}

//Broadcast sends a message to all subscribers
func Broadcast(msg *MessageObj) {
	jsonObj, err := json.Marshal(msg)
	if err != nil {
		log.Println("Broadcast json marshal:", err)
	}

	err = nknClient.PublishBinary(msg.TopicEncoded, jsonObj, &nkn.MessageConfig{
		TxPool: true,
	})
	if err != nil {
		log.Println("Broadcast send binary:", err)
	}
}

func (msgReceived MessageReceivedObj) Reply(msg *MessageObj) {
	jsonObj, err := json.Marshal(msg)

	if err != nil {
		log.Println("Reply json marshal:", err)
		return
	}

	_, err = nknClient.SendBinary(nkn.NewStringArray(msgReceived.Sender), jsonObj, &nkn.MessageConfig{
		TxPool: true,
	})
	if err != nil {
		log.Println("Reply send binary:", err)
		return
	}

}

func listen() {
	for {
		//Wait for a message
		msg := <-nknClient.OnMessage.C

		if msg != nil && msg.Data != nil {
			//try to unmarshal
			msgObj := MessageReceivedObj{}
			err := json.Unmarshal(msg.Data, &msgObj)
			if err != nil {
				log.Println("Received invalid message:", string(msg.Data), "from:", msg.Src, "error:", err)
				//} else if msg.Src == nknClient.Address() {
				//We exclude messages from ourselves
			} else {
				msgObj.Sender = msg.Src
				onMessageHandler(&msgObj)
			}
		}
	}
}
