package main

import (
	"./goconfig"
	"./protocol"
	estClient "./victoriest.org/client"
	"./victoriest.org/probe"
	"./victoriest.org/utils"
	log "code.google.com/p/log4go"
	// "image"
	// "image/jpeg"
	// "io"
	// "io/ioutil"
	"net"
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
	obj, _ := message.MsgContext.(map[string]interface{})
	log.Debug(obj["ChatMessage"])

	// if message.MsgType == 123 {
	// 	f3, err := os.Create("123321.jpg")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer f3.Close()

	// 	jpeg.Encode(f3, message.MsgContext.(image.Image), &jpeg.Options{90})
	// }

	// err := ioutil.WriteFile("temp.jpg", obj["ChatMessage"].([]byte), 0600)
	// if err != nil {
	// 	log.Error("err")
	// }
}

func msgSendHandler(client *estClient.VictoriestClient, writer *net.TCPConn, message string) {
	jsonProbe := new(probe.JsonProbe)
	chatMsg := protocol.ChatMsg{ChatMessage: message}
	broMsg := &probe.VictoriestMsg{MsgType: protocol.MSG_TYPE_CHAT_MESSGAE, MsgContext: chatMsg}

	strBuf, _ := jsonProbe.Serialize(broMsg)
	writer.Write(strBuf)

	// exefile, _ := exec.LookPath(os.Args[0])
	// filepath := path.Join(filepath.Dir(exefile), "./123.jpg")
	// file, err := os.Open(filepath)
	// if err != nil {
	// 	utils.CheckError(err, true)
	// }

	// m, _, err := image.Decode(file)
	// mm := &probe.VictoriestMsg{MsgType: protocol.MSG_TYPE_CHAT_MESSGAE, MsgContext: m}
	// strBuf, _ = jsonProbe.Serialize(mm)
	// log.Debug(mm)

	// var dest probe.VictoriestMsg
	// jsonProbe.Deserialize(strBuf, &dest)

	// f3, err := os.Create("123321.jpg")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f3.Close()
	// jpeg.Encode(f3, mm.MsgContext.(image.Image), &jpeg.Options{50})

	// err = ioutil.WriteFile("temp.jpg", b, 0600)
	// if err != nil {
	// 	log.Error("err")
	// }

	// broMsg := &probe.VictoriestMsg{MsgType: 123, MsgContext: m}
	// strBuf, _ := jsonProbe.Serialize(broMsg)
	// writer.Write(strBuf)

}
