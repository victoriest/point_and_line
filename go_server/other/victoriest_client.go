package main

import (
	log "github.com/alecthomas/log4go"
	"github.com/golang/protobuf/proto"
	"go_server/client"
	"go_server/goconfig"
	"go_server/protocol"
	"go_server/utils"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func main() {
	ip, port := readServerConfig()
	client := client.NewClient(ip, port, msgReceivedHandler, msgSendHandler)
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

func msgReceivedHandler(client *client.RobotClient, message *protocol.MobileSuiteModel) {
	// obj, _ := message.MsgContext.(map[string]interface{})

	log.Info(*message.Type)
	log.Info(message.Message)
	msg := &protocol.ChatMsg{}
	proto.Unmarshal(message.Message, msg)
	log.Info(msg)
	log.Info(*msg.ChatContext)
}

func msgSendHandler(client *client.RobotClient, writer *net.TCPConn, message string) {
	// probe := new(codec.ProtobufProbe)
	testMessage := &protocol.ChatMsg{
		ChatContext: proto.String(message),
	}
	byt, _ := proto.Marshal(testMessage)
	msg := &protocol.MobileSuiteModel{
		Type:    proto.Int32(protocol.MSG_TYPE_CHAT_MESSGAE),
		Message: byt,
	}
	bybuf, _ := client.Probe.Serialize(msg)
	writer.Write(bybuf)
}
