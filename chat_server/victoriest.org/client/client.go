package client

import (
	"../probe"
	"../utils"
	"bufio"
	log "code.google.com/p/log4go"
	"fmt"
	"net"
	"strings"
)

type IVictoriestClient interface {
	Startup()
}

type VictoriestClient struct {
	// 服务端端口号
	port string
	// IP
	ip string
	// 退出信号量
	quitSp chan bool
}

func NewVictoriestClient(ip string, port string) *VictoriestClient {
	client := new(VictoriestClient)
	client.ip = ip
	client.port = port
	client.quitSp = make(chan bool)
	return client
}

func (self *VictoriestClient) Startup() {
	strAddr := self.ip + ":" + self.port
	tcpAddr, err := net.ResolveTCPAddr("tcp", strAddr)
	utils.CheckError(err, true)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	utils.CheckError(err, true)
	defer conn.Close()
	fmt.Println("connecting ", conn.RemoteAddr().String(), "...")

	go self.readerPipe(conn)
	go self.writerPipe(conn)

	<-self.quitSp
}

func (self *VictoriestClient) writerPipe(conn *net.TCPConn) {
	jsonProbe := new(probe.JsonProbe)
	writer := bufio.NewWriter(conn)
	for {
		var msg string
		fmt.Scanln(&msg)

		if strings.EqualFold(msg, "quit") {
			self.quitSp <- true
			break
		}

		msgObj := probe.VictoriestMsg{MsgType: 1, MsgContext: msg}
		// switch obj := interface{}(msgObj).(type) {
		// case probe.VictoriestMsg:
		// 	log.Debug("write", obj)
		// default:
		// 	log.Debug("write not a VictoriestMsg")
		// }

		strBuf, _ := jsonProbe.Serialize(msgObj)
		writer.Write(strBuf)
		writer.Flush()
	}
}

func (self *VictoriestClient) readerPipe(conn *net.TCPConn) {
	reader := bufio.NewReader(conn)
	jsonProbe := new(probe.JsonProbe)
	for {
		var message interface{}
		err := jsonProbe.DeserializeByReader(reader, message)

		switch obj := (interface{}(message)).(type) {
		case probe.VictoriestMsg:
			log.Debug(obj.MsgContext)
		default:
			log.Debug("not a VictoriestMsg  ", obj)
		}
		utils.CheckError(err, true)
	}
}
