package main

import (
	"./goconfig"
	"./protocol"
	"./victoriest.org/probe"
	estServer "./victoriest.org/server"
	"./victoriest.org/utils"
	log "code.google.com/p/log4go"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func main() {
	log.LoadConfiguration("./log4go.config")
	server := estServer.NewVictoriestServer(readServerPort(), tcpHandler, connectedHandler, disconnectingHander)
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
func tcpHandler(server *estServer.VictoriestServer, message *probe.VictoriestMsg) {
	log.Debug(message)
	server.BroadcastMessage(message)
}

func connectedHandler(server *estServer.VictoriestServer, conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()
	chatMsg := protocol.ChatMsg{ChatMessage: "A new connection :" + ipStr}
	broMsg := &probe.VictoriestMsg{MsgType: protocol.MSG_TYPE_CHAT_MESSGAE, MsgContext: chatMsg}

	server.BroadcastMessage(broMsg)
}

func disconnectingHander(server *estServer.VictoriestServer, conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()
	chatMsg := protocol.ChatMsg{ChatMessage: "disconnected :" + ipStr}
	broMsg := &probe.VictoriestMsg{MsgType: protocol.MSG_TYPE_CHAT_MESSGAE, MsgContext: chatMsg}
	server.BroadcastMessage(broMsg)
}
