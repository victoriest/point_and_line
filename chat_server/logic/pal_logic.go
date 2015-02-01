package logic

import (
	"../protocol"
	sev "../server"
	log "code.google.com/p/log4go"
	"container/list"
	proto "github.com/golang/protobuf/proto"
	"net"
)

var inGameMap = make(map[string]string)

var ipMappingNick = make(map[string]string)

var joinGameList = list.New()

// 处理消息具体实现
func TcpHandler(server *sev.Nexus, ipStr string, message *protocol.MobileSuiteModel) {
	log.Debug(message)
	switch int32(*message.Type) {
	case int32(protocol.MessageType_MSG_TYPE_CHAT_MESSGAE_REQ):
		server.BroadcastMessage(message)
	case int32(protocol.MessageType_MSG_TYPE_LINE_A_POINT_REQ):
		processLinePoint(server, ipStr, message)
	case int32(protocol.MessageType_MSG_TYPE_SEARCH_A_GAME_REQ):
		processSearchGame(server, ipStr, message)
	case int32(protocol.MessageType_MSG_TYPE_CREATE_USER_REQ):
		processCreateUser(server, ipStr, message)
	case int32(protocol.MessageType_MSG_TYPE_LOGIN_REQ):
		processLogin(server, ipStr, message)
	}
}

func ConnectedHandler(server *sev.Nexus, conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()
	str := "A new connection :" + ipStr
	chatMsg := &protocol.ChatMsg{ChatContext: &str}
	byt, _ := proto.Marshal(chatMsg)
	broBack(server, byt, int32(protocol.MessageType_MSG_TYPE_CHAT_MESSAGE_RES))
}

func DisconnectingHander(server *sev.Nexus, conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()
	str := "disconnected :" + ipStr
	chatMsg := &protocol.ChatMsg{ChatContext: &str}
	byt, _ := proto.Marshal(chatMsg)
	endGame(ipStr)
	broBack(server, byt, int32(protocol.MessageType_MSG_TYPE_CHAT_MESSAGE_RES))
	broBack(server, byt, int32(protocol.MessageType_MSG_TYPE_LOGOUT_RES))
}

func sendBack(server *sev.Nexus, ipStr string, byt []byte, msgType int32) {
	lpDtoMsg := &protocol.MobileSuiteModel{
		Type: proto.Int32(msgType),
	}
	if len(byt) > 0 {
		lpDtoMsg.Message = byt
	}
	server.SendTo(ipStr, lpDtoMsg)
}

func broBack(server *sev.Nexus, byt []byte, msgType int32) {
	broMsg := &protocol.MobileSuiteModel{
		Type:    proto.Int32(msgType),
		Message: byt,
	}
	if len(byt) > 0 {
		broMsg.Message = byt
	}
	server.BroadcastMessage(broMsg)
}
