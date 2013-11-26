package probe

const SEND_BRO_GLOBAL string = "BRO_GLOBAL"

const SEND_BRO_GROUP string = "BRO_GROUP"

const SEND_TO_SERVER string = "TO_SERVER"

type VictoriestMsg struct {
	MsgType    int32
	MsgContext interface{}
}

type TestMsg struct {
	MsgInt      int32
	ChatMessage string
}
