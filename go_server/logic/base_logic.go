package logic

import (
	"encoding/json"

	log "github.com/alecthomas/log4go"
	"github.com/golang/protobuf/proto"
	"go_server/codec"
	"go_server/dao"
	"go_server/protocol"
	sev "go_server/server"
	//"strconv"
)

func processCreateUser(server *sev.Nexus, ipStr string, message interface{}) {
	createUserDto := &protocol.CreateUserDTO{}
	if server.ProtocolType == sev.ProtocolTypeTCP {
		proto.Unmarshal(message.(*protocol.MobileSuiteModel).Message, createUserDto)
	} else if server.ProtocolType == sev.ProtocolTypeWebSocket {
		bodyBytes, _ := json.Marshal(message.(*codec.VictoriestMsg).MsgContext)
		json.Unmarshal(bodyBytes, &createUserDto)
	}

	user := &dao.User{}
	user.Name = *createUserDto.UName
	user.Round = 0
	user.WinCount = 0
	user.Rank = 0
	user.Pwd = *createUserDto.Pwd
	user.OpenId = *createUserDto.UName

	result, err := server.DbConnector.Insert(user)
	if err != nil {
		log.Info(err)
		resp, _ := genResponseDTO(server, nil, int32(protocol.MessageType_MSG_TYPE_CREATE_USER_RES))
		server.SendTo(ipStr, resp)
		// sendBack(server, ipStr, nil,
		// 	int32(protocol.MessageType_MSG_TYPE_CREATE_USER_RES))
		return
	}
	createResult := &protocol.CreateResultDTO{
		UserId: proto.Int64(int64(result)),
	}
	resp, _ := genResponseDTO(server, createResult, int32(protocol.MessageType_MSG_TYPE_CREATE_USER_RES))
	server.SendTo(ipStr, resp)
	// byt, _ := proto.Marshal(createResult)
	// sendBack(server, ipStr, byt,
	// 	int32(protocol.MessageType_MSG_TYPE_CREATE_USER_RES))
}

func processLogin(server *sev.Nexus, ipStr string,
	message interface{}) {
	loginDto := &protocol.LoginDTO{}
	var user *dao.User
	if server.ProtocolType == sev.ProtocolTypeTCP {
		proto.Unmarshal(message.(*protocol.MobileSuiteModel).Message, loginDto)
		userArr, err := server.DbConnector.QueryByUserName(
			*loginDto.UName, *loginDto.Pwd)
		if err != nil || userArr == nil || len(userArr) < 1 {
			resp, _ := genResponseDTO(server, nil, int32(protocol.MessageType_MSG_TYPE_LOGIN_RES))
			server.SendTo(ipStr, resp)
			// sendBack(server, ipStr, nil,
			// 	int32(protocol.MessageType_MSG_TYPE_LOGIN_RES))
			return
		}
		user = &userArr[0]
	} else if server.ProtocolType == sev.ProtocolTypeWebSocket {
		bodyBytes, _ := json.Marshal(message.(*codec.VictoriestMsg).MsgContext)
		json.Unmarshal(bodyBytes, loginDto)
		// OPENID相关处理逻辑
		log.Info(*loginDto.UName)
		userArr, err := server.DbConnector.QueryByOpenId(*loginDto.UName)
		if err != nil || userArr == nil || len(userArr) < 1 {
			// TODO 如果没有用户 就创建一个
			user := &dao.User{}
			user.Name = *loginDto.UName
			user.Round = 0
			user.WinCount = 0
			user.Rank = 0
			user.Pwd = *loginDto.Pwd
			user.OpenId = *loginDto.UName
			result, _ := server.DbConnector.Insert(user)
			user.Id = result
		} else {
			user = &userArr[0]

			log.Info(&userArr[0])
			log.Info(user)
		}
	}

	//loginResultDto := &protocol.LoginResultDTO{}
	//loginResultDto.UserId = loginDto.UserId
	//strName := strconv.FormatInt(*loginResultDto.UserId, 10)
	//loginResultDto.UName = &strName

	// 登陆成功
	loginResultDto := &protocol.LoginResultDTO{}
	loginResultDto.UserId = proto.Int64(user.Id)
	loginResultDto.UName = &(user.Name)
	loginResultDto.Round = proto.Int32(int32(user.Round))
	loginResultDto.WinCount = proto.Int32(int32(user.WinCount))
	loginResultDto.Rank = proto.Int32(int32(user.Rank))
	log.Info(loginResultDto)
	// byt, _ := proto.Marshal(loginResultDto)
	// sendBack(server, ipStr, byt,
	// 	int32(protocol.MessageType_MSG_TYPE_LOGIN_RES))

	response, _ := genResponseDTO(server, loginResultDto, int32(protocol.MessageType_MSG_TYPE_LOGIN_RES))
	server.SendTo(ipStr, response)
}

func processChatMessage(server *sev.Nexus, message interface{}) {
	chatDto := &protocol.ChatMsg{}
	if server.ProtocolType == sev.ProtocolTypeTCP {
		proto.Unmarshal(message.(*protocol.MobileSuiteModel).Message, chatDto)
	} else if server.ProtocolType == sev.ProtocolTypeWebSocket {
		bodyBytes, _ := json.Marshal(message.(*codec.VictoriestMsg).MsgContext)
		json.Unmarshal(bodyBytes, &chatDto)
	}
	resp, _ := genResponseDTO(server, chatDto, int32(protocol.MessageType_MSG_TYPE_CHAT_MESSAGE_RES))
	server.BroadcastMessage(resp)
	// byt, _ := proto.Marshal(chatDto)
	// broBack(server, byt, int32(protocol.MessageType_MSG_TYPE_CHAT_MESSAGE_RES))
}

func genResponseDTO(server *sev.Nexus, response interface{}, messageType int32) (interface{}, error) {
	if server.ProtocolType == sev.ProtocolTypeTCP {
		var resByte []byte
		if response != nil {
			resByte, _ = proto.Marshal(response.(proto.Message))
		}
		resDtoMsg := &protocol.MobileSuiteModel{
			Type: proto.Int32(messageType),
		}
		if len(resByte) > 0 {
			resDtoMsg.Message = resByte
		}
		return resDtoMsg, nil
	}
	resDtoMsg := &codec.VictoriestMsg{
		MsgType:    messageType,
		MsgContext: response,
	}
	return resDtoMsg, nil
}
