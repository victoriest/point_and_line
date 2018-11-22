package logic

import (
	"container/list"
	"fmt"

	"../codec"
	"../protocol"
	sev "../server"
	game "./games"
	log "github.com/alecthomas/log4go"
	proto "github.com/golang/protobuf/proto"
)

var inGameMap = make(map[string]string)

var ipMappingNick = make(map[string]string)

var joinGameList = list.New()

var gameObjMap = make(map[string]*game.PointAndLineGame)

// TCPHandler 处理消息具体实现
func TCPHandler(server *sev.Nexus, ipStr string, message interface{}) {
	log.Debug(message)
	var messageType int32
	if server.ProtocolType == sev.ProtocolTypeTCP {
		messageType = int32(*message.(*protocol.MobileSuiteModel).Type)
	} else if server.ProtocolType == sev.ProtocolTypeWebSocket {
		messageType = int32(message.(*codec.VictoriestMsg).MsgType)
	}
	switch messageType {
	case int32(protocol.MessageType_MSG_TYPE_CHAT_MESSGAE_REQ):
		processChatMessage(server, message)
	case int32(protocol.MessageType_MSG_TYPE_LINE_A_POINT_REQ):
		processLinePoint(server, ipStr, message)
	case int32(protocol.MessageType_MSG_TYPE_SEARCH_A_GAME_REQ):
		processSearchGame(server, ipStr, message)
	case int32(protocol.MessageType_MSG_TYPE_CREATE_USER_REQ):
		processCreateUser(server, ipStr, message)
	case int32(protocol.MessageType_MSG_TYPE_LOGIN_REQ):
		processLogin(server, ipStr, message)
	case int32(protocol.MessageType_MSG_TYPE_END_GAME_REQ):
		processEndGame(server, ipStr, message)
	}
}

func ConnectedHandler(server *sev.Nexus, ipStr string) {
	// ipStr := conn.RemoteAddr().String()
	str := "A new connection :" + ipStr
	chatMsg := &protocol.ChatMsg{ChatContext: &str}
	byt, _ := proto.Marshal(chatMsg)
	broBack(server, byt, int32(protocol.MessageType_MSG_TYPE_CHAT_MESSAGE_RES))
}

func DisconnectingHander(server *sev.Nexus, ipStr string) {
	// ipStr := conn.RemoteAddr().String()
	str := "disconnected :" + ipStr
	chatMsg := &protocol.ChatMsg{ChatContext: &str}
	byt, _ := proto.Marshal(chatMsg)
	endGame(server, ipStr)
	broBack(server, byt, int32(protocol.MessageType_MSG_TYPE_CHAT_MESSAGE_RES))

}

func endGame(server *sev.Nexus, ipStr string) {
	opptIpStr, hasKey := inGameMap[ipStr]
	if !hasKey {
		return
	}
	logoutDto := &protocol.LogoutDTO{UserId: proto.Int64(1)}
	byt, _ := proto.Marshal(logoutDto)
	sendBack(server, opptIpStr, byt, int32(protocol.MessageType_MSG_TYPE_END_GAME_RES))
	fmt.Println(opptIpStr)
	delete(ipMappingNick, opptIpStr)
	delete(ipMappingNick, ipStr)
	delete(inGameMap, opptIpStr)
	delete(inGameMap, ipStr)
	delete(gameObjMap, ipStr)
	delete(gameObjMap, opptIpStr)
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
