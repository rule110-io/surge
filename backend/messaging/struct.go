package messaging

type MessageObj struct {
	Type         int
	TopicEncoded string
	Data         []byte
}

type MessageReceivedObj struct {
	Type         int
	TopicEncoded string
	Data         []byte
	Sender       string
}

//Message Types
const (
	MsgRequestFiles = iota
	MsgResponseFiles
	MsgRequestFileInfo
	MsgResponseFileInfo
)
