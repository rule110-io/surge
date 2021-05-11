package messaging

type MessageObj struct {
	Type  int
	Topic string
	Data  []byte
}

type MessageReceivedObj struct {
	Type   int
	Topic  string
	Data   []byte
	Sender string
}

//Message Types
const (
	MsgRequestFiles = iota
	MsgResponseFiles
	MsgRequestFileInfo
	MsgResponseFileInfo
)
