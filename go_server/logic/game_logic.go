package logic

import (
	"../protocol"
	sev "../server"
	game "./games"
	log "code.google.com/p/log4go"
	proto "github.com/golang/protobuf/proto"
)

func processLinePoint(server *sev.Nexus, ipStr string, message *protocol.MobileSuiteModel) {
	to := inGameMap[ipStr]
	lpDto := &protocol.LineAPointDTO{}
	proto.Unmarshal(message.Message, lpDto)

	gameObj := gameObjMap[ipStr]
	if gameObj == nil {
		resByte, _ := proto.Marshal(result)
		resDtoMsg := &protocol.MobileSuiteModel{
			Type:    proto.Int32(int32(protocol.MessageType_MSG_TYPE_LINE_A_POINT_TO_REQUEST_RES)),
			Message: resByte,
		}
		server.SendTo(ipStr, resDtoMsg)
		return
	}

	result := gameObj.Line(lpDto.Row, lpDto.Col, lpDto.PlayerIndex)
	if result == 0 {
		byt, _ := proto.Marshal(lpDto)
		lpDtoMsg := &protocol.MobileSuiteModel{
			Type:    proto.Int32(int32(protocol.MessageType_MSG_TYPE_LINE_A_POINT_RES)),
			Message: byt,
		}
		server.SendTo(to, lpDtoMsg)
	}

	resByte, _ := proto.Marshal(result)
	resDtoMsg := &protocol.MobileSuiteModel{
		Type:    proto.Int32(int32(protocol.MessageType_MSG_TYPE_LINE_A_POINT_TO_REQUEST_RES)),
		Message: resByte,
	}
	server.SendTo(ipStr, resDtoMsg)

	// TODO 返回一条消息

}

func processSearchGame(server *sev.Nexus, ipStr string, message *protocol.MobileSuiteModel) {
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

		if !server.ConnectionIsOpen(strIp1) {
			return
		}

		jgIp2 := joinGameList.Front()
		strIp2 := jgIp2.Value.(string)
		joinGameList.Remove(jgIp2)

		if !server.ConnectionIsOpen(strIp2) {
			return
		}

		inGameMap[strIp1] = strIp2
		inGameMap[strIp2] = strIp1
		gameObj := games.NewPointAndLineGame(2)
		gameObjMap[strIp1] = gameObj
		gameObjMap[strIp2] = gameObj

		gsDto1 := &protocol.GameStartDTO{
			OpptName:    proto.String(ipMappingNick[strIp2]),
			PlayerIndex: proto.Int32(1),
		}
		byt1, _ := proto.Marshal(gsDto1)
		sendBack(server, strIp1, byt1, int32(protocol.MessageType_MSG_TYPE_START_RES))

		gsDto2 := &protocol.GameStartDTO{
			OpptName:    proto.String(ipMappingNick[strIp1]),
			PlayerIndex: proto.Int32(2),
		}
		byt2, _ := proto.Marshal(gsDto2)
		sendBack(server, strIp2, byt2, int32(protocol.MessageType_MSG_TYPE_START_RES))

	}
}

func processEndGame(server *sev.Nexus, ipStr string, message *protocol.MobileSuiteModel) {
	endGame(server, ipStr)
}
