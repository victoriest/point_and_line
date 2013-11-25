package probe

const SEND_BRO_GLOBAL int32 = 0x100

const SEND_BRO_GROUP int32 = 0x101

const SEND_TO_SERVER int32 = 0x102

type VictoriestMsg struct {
	sendFrom string
	sendTo string
	MsgContext interface{}
}
