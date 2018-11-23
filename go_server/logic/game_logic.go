package logic

import (
	"encoding/json"

	"../codec"
	"../protocol"
	sev "../server"
	"./games"
	log "github.com/alecthomas/log4go"
	proto "github.com/golang/protobuf/proto"
)

func processLinePoint(server *sev.Nexus, ipStr string, message interface{}) {
	to := inGameMap[ipStr]
	lpDto := &protocol.LineAPointDTO{}
	if server.ProtocolType == sev.ProtocolTypeTCP {
		proto.Unmarshal(message.(*protocol.MobileSuiteModel).Message, lpDto)
	} else if server.ProtocolType == sev.ProtocolTypeWebSocket {
		bodyBytes, _ := json.Marshal(message.(*codec.VictoriestMsg).MsgContext)
		json.Unmarshal(bodyBytes, &lpDto)
	}

	gameObj := gameObjMap[ipStr]
	if gameObj == nil {
		response := &protocol.LineAPointResponseDTO{
			Result: proto.Int32(int32(-1)),
		}
		resByte, _ := proto.Marshal(response)
		resDtoMsg := &protocol.MobileSuiteModel{
			Type:    proto.Int32(int32(protocol.MessageType_MSG_TYPE_LINE_A_POINT_TO_REQUEST_RES)),
			Message: resByte,
		}
		server.SendTo(ipStr, resDtoMsg)
		return
	}

	result := gameObj.Line(int(*lpDto.Row), int(*lpDto.Col), int(*lpDto.PlayerIndex))
	if result == 0 {
		byt, _ := proto.Marshal(lpDto)
		lpDtoMsg := &protocol.MobileSuiteModel{
			Type:    proto.Int32(int32(protocol.MessageType_MSG_TYPE_LINE_A_POINT_RES)),
			Message: byt,
		}
		server.SendTo(to, lpDtoMsg)
	}

	response := &protocol.LineAPointResponseDTO{
		Result: proto.Int32(int32(result)),
	}
	resByte, _ := proto.Marshal(response)
	resDtoMsg := &protocol.MobileSuiteModel{
		Type:    proto.Int32(int32(protocol.MessageType_MSG_TYPE_LINE_A_POINT_TO_REQUEST_RES)),
		Message: resByte,
	}
	server.SendTo(ipStr, resDtoMsg)

}

func processSearchGame(server *sev.Nexus, ipStr string, message interface{}) {
	msg := &protocol.ChatMsg{}
	if server.ProtocolType == sev.ProtocolTypeTCP {
		proto.Unmarshal(message.(*protocol.MobileSuiteModel).Message, msg)
	} else if server.ProtocolType == sev.ProtocolTypeWebSocket {
		bodyBytes, _ := json.Marshal(message.(*codec.VictoriestMsg).MsgContext)
		json.Unmarshal(bodyBytes, &msg)
	}

	log.Info(*msg.ChatContext)
	ipMappingNick[ipStr] = *msg.ChatContext
	joinGameList.PushBack(ipStr)

	// 如果人多了 就开一场游戏
	if joinGameList.Len() >= 2 {
		jgIP1 := joinGameList.Front()
		strIP1 := jgIP1.Value.(string)
		joinGameList.Remove(jgIP1)

		if !server.ConnectionIsOpen(strIP1) {
			return
		}

		jgIP2 := joinGameList.Front()
		strIP2 := jgIP2.Value.(string)
		joinGameList.Remove(jgIP2)

		if !server.ConnectionIsOpen(strIP2) {
			return
		}

		inGameMap[strIP1] = strIP2
		inGameMap[strIP2] = strIP1
		gameObj := games.NewPointAndLineGame(2)
		gameObjMap[strIP1] = gameObj
		gameObjMap[strIP2] = gameObj

		gsDto1 := &protocol.GameStartDTO{
			OpptName:    proto.String(ipMappingNick[strIP2]),
			PlayerIndex: proto.Int32(1),
		}
		resp, _ := genResponseDTO(server, gsDto1, int32(protocol.MessageType_MSG_TYPE_START_RES))
		server.SendTo(strIP1, resp)
		// byt1, _ := proto.Marshal(gsDto1)
		// sendBack(server, strIP1, byt1, int32(protocol.MessageType_MSG_TYPE_START_RES))

		gsDto2 := &protocol.GameStartDTO{
			OpptName:    proto.String(ipMappingNick[strIP1]),
			PlayerIndex: proto.Int32(2),
		}
		resp, _ = genResponseDTO(server, gsDto2, int32(protocol.MessageType_MSG_TYPE_START_RES))
		server.SendTo(strIP2, resp)
		// byt2, _ := proto.Marshal(gsDto2)
		// sendBack(server, strIP2, byt2, int32(protocol.MessageType_MSG_TYPE_START_RES))

	}
}

func processEndGame(server *sev.Nexus, ipStr string, message interface{}) {
	endGame(server, ipStr)
}
