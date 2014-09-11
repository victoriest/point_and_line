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

	// result, _ := server.DbConnector.Insert(user)
	// createResult := &protocol.CreateResultDTO{
	// 	UserId: result,
	// }

	// byt, _ := proto.Marshal(*createResult)

	// lpDtoMsg := &protocol.MobileSuiteModel{
	// 	Type:    proto.Int32(int32(protocol.MessageType_MSG_TYPE_CREATE_USER_RES)),
	// 	Message: byt,
	// }
	// server.SendTo(to, lpDtoMsg)

}

func procerssLogin(server *sev.Nexus, ipStr string, message *protocol.MobileSuiteModel) {
	to := inGameMap[ipStr]
	loginDto := &protocol.LoginDTO{}
	proto.Unmarshal(message.Message, loginDto)

	userArr, err := server.DbConnector.QueryByUserId(int(*loginDto.UserId))
	loginResultDto := &protocol.LoginResultDTO{}
	if len(userArr) > 0 {
		user := dao.User(userArr[0])
		// 登陆成功
		loginResultDto.UserId = user.Id
		loginResultDto.Name = user.Name
		loginResultDto.Round = user.Round
		loginResultDto.WinCount = user.WinCount
		loginResultDto.Rank = user.Rank
	}

	byt, _ := proto.Marshal(loginResultDto)
	lpDtoMsg := &protocol.MobileSuiteModel{
		Type:    proto.Int32(int32(protocol.MessageType_MSG_TYPE_LOGIN_RES)),
		Message: byt,
	}

	server.SendTo(to, lpDtoMsg)

}
