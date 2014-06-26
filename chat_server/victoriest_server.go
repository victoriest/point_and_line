package main

import (
	"./goconfig"
	"./protocol"
	sev "./server"
	"./utils"
	proto "code.google.com/p/goprotobuf/proto"
	log "code.google.com/p/log4go"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

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
func tcpHandler(server *sev.Nexus, message *protocol.MobileSuiteModel) {
	log.Debug(message)
	server.BroadcastMessage(message)
}

func connectedHandler(server *sev.Nexus, conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()
	str := "A new connection :" + ipStr
	chatMsg := &protocol.ChatMsg{ChatContext: &str}
	byt, _ := proto.Marshal(chatMsg)
	broMsg := &protocol.MobileSuiteModel{
		Type:    proto.Int32(protocol.MSG_TYPE_CHAT_MESSGAE),
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
		Type:    proto.Int32(protocol.MSG_TYPE_CHAT_MESSGAE),
		Message: byt,
	}
	server.BroadcastMessage(broMsg)
}
