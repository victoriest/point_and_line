package protocol

const MSG_TYPE_CHAT_MESSGAE int32 = 0xA

type ChatMsg struct {
	ChatMessage string
}
