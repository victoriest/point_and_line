package main

import (
	"./goconfig"
	"./protocol"
	sev "./server"
	"./utils"
	proto "code.google.com/p/goprotobuf/proto"
	log "code.google.com/p/log4go"
	"container/list"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

var inGameMap = make(map[string]string)

var ipMappingNick = make(map[string]string)

var joinGameList = list.New()

func main() {
	log.LoadConfiguration("./log4go.config")
	server := sev.NewNexus(readServerPort(), tcpHandler, connectedHandler, disconnectingHander)
	server.Startup()
}

// 读取配置文件
func readServerPort() string {
	exefile, _ := exec.LookPath(os.Args[0])
	log.Info(filepath.Dir(exefile))
	filepath := path.Join(filepath.Dir(exefile), "./server.config")
	cf, err := goconfig.LoadConfigFile(filepath)
	utils.CheckError(err, true)
	port, err := cf.GetValue(goconfig.DEFAULT_SECTION, "server.port")
	utils.CheckError(err, true)
	return port
}

// 处理消息具体实现
func tcpHandler(server *sev.Nexus, ipStr string, message *protocol.MobileSuiteModel) {
	log.Debug(message)
	switch int32(*message.Type) {

	case int32(protocol.MessageType_MSG_TYPE_CHAT_MESSGAE):
		server.BroadcastMessage(message)

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
		}
	}
}

func connectedHandler(server *sev.Nexus, conn *net.TCPConn) {
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

func disconnectingHander(server *sev.Nexus, conn *net.TCPConn) {
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
