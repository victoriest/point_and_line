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
		procerssLogin(server, ipStr, message)
	}
}

func ConnectedHandler(server *sev.Nexus, conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()
	str := "A new connection :" + ipStr
	chatMsg := &protocol.ChatMsg{ChatContext: &str}
	byt, _ := proto.Marshal(chatMsg)
	broMsg := &protocol.MobileSuiteModel{
		Type:    proto.Int32(int32(protocol.MessageType_MSG_TYPE_CHAT_MESSAGE_RES)),
		Message: byt,
	}
	server.BroadcastMessage(broMsg)
}

func DisconnectingHander(server *sev.Nexus, conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()
	str := "disconnected :" + ipStr
	chatMsg := &protocol.ChatMsg{ChatContext: &str}
	byt, _ := proto.Marshal(chatMsg)
	broMsg := &protocol.MobileSuiteModel{
		Type:    proto.Int32(int32(protocol.MessageType_MSG_TYPE_CHAT_MESSAGE_RES)),
		Message: byt,
	}
	server.BroadcastMessage(broMsg)
}
