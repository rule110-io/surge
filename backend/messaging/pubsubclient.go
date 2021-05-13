package messaging

import (
	"encoding/json"
	"fmt"

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
	fmt.Println("Marshalled Bytes:", jsonObj)

	err = nknClient.PublishBinary(msg.Topic, jsonObj, &nkn.MessageConfig{
		TxPool: true,
	})
	if err != nil {
		fmt.Println(err)
	}
}

func (msgReceived MessageReceivedObj) Reply(msg *MessageObj) {
	jsonObj, err := json.Marshal(msg)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Marshalled Bytes:", jsonObj)

	nknClient.SendBinary(nkn.NewStringArray(msgReceived.Sender), jsonObj, &nkn.MessageConfig{
		TxPool: true,
	})
}

func listen() {
	for true {
		//Wait for a message
		msg := <-nknClient.OnMessage.C

		//try to unmarshal
		msgObj := MessageReceivedObj{}
		err := json.Unmarshal(msg.Data, &msgObj)
		if err != nil {
			fmt.Println("Received invalid message:", string(msg.Data), "from:", msg.Src, "error:", err)
		} else {
			msgObj.Sender = msg.Src
			onMessageHandler(&msgObj)
		}
	}
}
