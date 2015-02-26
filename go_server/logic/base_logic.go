package logic

import (
	"../dao"
	"../protocol"
	sev "../server"
	log "code.google.com/p/log4go"
	proto "github.com/golang/protobuf/proto"
	//"strconv"
)

func processCreateUser(server *sev.Nexus, ipStr string, message *protocol.MobileSuiteModel) {
	createUserDto := &protocol.CreateUserDTO{}
	proto.Unmarshal(message.Message, createUserDto)

	user := &dao.User{}
	user.Name = *createUserDto.UName
	user.Round = 0
	user.WinCount = 0
	user.Rank = 0
	user.Pwd = *createUserDto.Pwd

	result, err := server.DbConnector.Insert(user)
	if err != nil {
		log.Info(err)
		sendBack(server, ipStr, nil,
			int32(protocol.MessageType_MSG_TYPE_CREATE_USER_RES))
		return
	}
	createResult := &protocol.CreateResultDTO{
		UserId: proto.Int64(int64(result)),
	}
	byt, _ := proto.Marshal(createResult)

	sendBack(server, ipStr, byt,
		int32(protocol.MessageType_MSG_TYPE_CREATE_USER_RES))
}

func processLogin(server *sev.Nexus, ipStr string,
	message *protocol.MobileSuiteModel) {
	loginDto := &protocol.LoginDTO{}
	proto.Unmarshal(message.Message, loginDto)

	userArr, err := server.DbConnector.QueryByUserName(
		*loginDto.UName, *loginDto.Pwd)
	if err != nil || userArr == nil || len(userArr) < 1 {
		sendBack(server, ipStr, nil,
			int32(protocol.MessageType_MSG_TYPE_LOGIN_RES))
		return
	}

	//loginResultDto := &protocol.LoginResultDTO{}
	//loginResultDto.UserId = loginDto.UserId
	//strName := strconv.FormatInt(*loginResultDto.UserId, 10)
	//loginResultDto.UName = &strName

	user := dao.User(userArr[0])

	// 登陆成功
	loginResultDto := &protocol.LoginResultDTO{}
	loginResultDto.UserId = proto.Int64(user.Id)
	loginResultDto.UName = &user.Name
	loginResultDto.Round = proto.Int32(int32(user.Round))
	loginResultDto.WinCount = proto.Int32(int32(user.WinCount))
	loginResultDto.Rank = proto.Int32(int32(user.Rank))

	byt, _ := proto.Marshal(loginResultDto)
	sendBack(server, ipStr, byt,
		int32(protocol.MessageType_MSG_TYPE_LOGIN_RES))
}

func processChatMessage(server *sev.Nexus, message *protocol.MobileSuiteModel) {
	chatDto := &protocol.ChatMsg{}
	proto.Unmarshal(message.Message, chatDto)

	byt, _ := proto.Marshal(chatDto)
	broBack(server, byt, int32(protocol.MessageType_MSG_TYPE_CHAT_MESSAGE_RES))
}
