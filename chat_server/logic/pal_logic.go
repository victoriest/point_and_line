package logic

import (
	"../protocol"
	sev "../server"
	proto "code.google.com/p/goprotobuf/proto"
	log "code.google.com/p/log4go"
	"container/list"
	"net"
)

var inGameMap = make(map[string]string)

var ipMappingNick = make(map[string]string)

var joinGameList = list.New()

// 处理消息具体实现
func TcpHandler(server *sev.Nexus, ipStr string, message *protocol.MobileSuiteModel) {
	log.Debug(message)
	switch int32(*message.Type) {

	case int32(protocol.MessageType_MSG_TYPE_CHAT_MESSGAE):
		server.BroadcastMessage(message)

	case int32(protocol.MessageType_MSG_TYPE_LINE_A_POINT_REQ):
		to := inGameMap[ipStr]
		lpDto := &protocol.LineAPointDTO{}
		proto.Unmarshal(message.Message, lpDto)

		byt, _ := proto.Marshal(lpDto)
		lpDtoMsg := &protocol.MobileSuiteModel{
			Type:    proto.Int32(int32(protocol.MessageType_MSG_TYPE_LINE_A_POINT_RES)),
			Message: byt,
		}
		server.SendTo(to, lpDtoMsg)

	case int32(protocol.MessageType_MSG_TYPE_SEARCH_A_GAME_REQ):
		msg := &protocol.ChatMsg{}
		proto.Unmarshal(message.Message, msg)
		log.Info(*msg.ChatContext)
		ipMappingNick[ipStr] = *msg.ChatContext
		joinGameList.PushBack(ipStr)
		// 如果人多了 就开一场游戏
		if joinGameList.Len() >= 2 {
			jgIp1 := joinGameList.Front()
			strIp1 := jgIp1.Value.(string)
			joinGameList.Remove(jgIp1)

			jgIp2 := joinGameList.Front()
			strIp2 := jgIp2.Value.(string)
			joinGameList.Remove(jgIp2)

			inGameMap[strIp1] = strIp2
			inGameMap[strIp2] = strIp1

			gsDto1 := &protocol.GameStartDTO{
				OpptName:    proto.String(ipMappingNick[strIp2]),
				PlayerIndex: proto.Int32(1),
			}
			byt1, _ := proto.Marshal(gsDto1)
			broMsg1 := &protocol.MobileSuiteModel{
				Type:    proto.Int32(int32(protocol.MessageType_MSG_TYPE_START_RES)),
				Message: byt1,
			}
			server.SendTo(strIp1, broMsg1)

			gsDto2 := &protocol.GameStartDTO{
				OpptName:    proto.String(ipMappingNick[strIp1]),
				PlayerIndex: proto.Int32(2),
			}
			byt2, _ := proto.Marshal(gsDto2)
			broMsg2 := &protocol.MobileSuiteModel{
				Type:    proto.Int32(int32(protocol.MessageType_MSG_TYPE_START_RES)),
				Message: byt2,
			}
			server.SendTo(strIp2, broMsg2)
		}
	}
}

func ConnectedHandler(server *sev.Nexus, conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()
	str := "A new connection :" + ipStr
	chatMsg := &protocol.ChatMsg{ChatContext: &str}
	byt, _ := proto.Marshal(chatMsg)
	broMsg := &protocol.MobileSuiteModel{
		Type:    proto.Int32(int32(protocol.MessageType_MSG_TYPE_CHAT_MESSGAE)),
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
		Type:    proto.Int32(int32(protocol.MessageType_MSG_TYPE_CHAT_MESSGAE)),
		Message: byt,
	}
	server.BroadcastMessage(broMsg)
}
