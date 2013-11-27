package main

import (
	"./goconfig"
	"./protocol"
	estClient "./victoriest.org/client"
	"./victoriest.org/probe"
	"./victoriest.org/utils"
	"bufio"
	log "code.google.com/p/log4go"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func main() {
	ip, port := readServerConfig()
	client := estClient.NewVictoriestClient(ip, port, msgReceivedHandler, msgSendHandler)
	client.Startup()
}

func readServerConfig() (string, string) {
	exefile, _ := exec.LookPath(os.Args[0])

	filepath := path.Join(filepath.Dir(exefile), "./server.config")
	cf, err := goconfig.LoadConfigFile(filepath)
	utils.CheckError(err, true)

	host, err := cf.GetValue(goconfig.DEFAULT_SECTION, "server.host")
	utils.CheckError(err, true)

	port, err := cf.GetValue(goconfig.DEFAULT_SECTION, "server.port")
	utils.CheckError(err, true)
	return host, port
}

func msgReceivedHandler(client *estClient.VictoriestClient, message *probe.VictoriestMsg) {
	log.Debug(message)
}

func msgSendHandler(client *estClient.VictoriestClient, writer *bufio.Writer, message string) {
	jsonProbe := new(probe.JsonProbe)
	chatMsg := protocol.ChatMsg{ChatMessage: message}
	broMsg := probe.VictoriestMsg{MsgType: protocol.MSG_TYPE_CHAT_MESSGAE, MsgContext: chatMsg}

	strBuf, _ := jsonProbe.Serialize(broMsg, protocol.MSG_TYPE_CHAT_MESSGAE)
	writer.Write(strBuf)
	writer.Flush()
}
