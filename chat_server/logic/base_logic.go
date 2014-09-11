package logic

import (
	"../dao"
	"../protocol"
	sev "../server"
	proto "code.google.com/p/goprotobuf/proto"
)

func processCreateUser(server *sev.Nexus, ipStr string, message *protocol.MobileSuiteModel) {
	to := inGameMap[ipStr]
	createUserDto := &protocol.CreateUserDTO{}
	proto.Unmarshal(message.Message, createUserDto)

	user := &dao.User{}
	user.Name = *createUserDto.Name
	user.Round = 0
	user.WinCount = 0
	user.Rank = 0

	result, err := server.DbConnector.Insert(user)
	msgRes := proto.Int32(int32(result))
	byt, _ := proto.Marshal(*msgRes)

	lpDtoMsg := &protocol.MobileSuiteModel{
		Type:    proto.Int32(int32(protocol.MessageType_MSG_TYPE_CREATE_USER_RES)),
		Message: byt,
	}
	server.SendTo(to, lpDtoMsg)

}

func procerssLogin(server *sev.Nexus, ipStr string, message *protocol.MobileSuiteModel) {
	to := inGameMap[ipStr]
	var userId int
	proto.Unmarshal(message.Message, userId)

	userArr, err := server.DbConnector.QueryByUserId(userId)
	loginDto := &protocol.LoginDTO{}
	if len(userArr) > 0 {
		// 登陆成功
		loginDto.UserId = userArr[0].Id
		loginDto.Name = userArr[0].Name
		loginDto.Round = userArr[0].Round
		loginDto.WinCount = userArr[0].WinCound
		loginDto.Rank = userArr[0].Rank
	}

	byt, _ := proto.Marshal(loginDto)
	lpDtoMsg := &protocol.MobileSuiteModel{
		Type:    proto.Int32(int32(protocol.MessageType_MSG_TYPE_CREATE_USER_RES)),
		Message: byt,
	}

	server.SendTo(to, lpDtoMsg)

}
