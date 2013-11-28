package client

import (
	"../probe"
	"../utils"
	"bufio"
	"bytes"
	log "code.google.com/p/log4go"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

// 消息逻辑处理托管Handler
type MessageReceivedHandler func(*VictoriestClient, *probe.VictoriestMsg)
type MessageSenderHandler func(*VictoriestClient, *net.TCPConn, string)

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
	// 消息逻辑处理托管Handler
	recivedHandler MessageReceivedHandler
	// 消息逻辑处理托管Handler
	sendHandler MessageSenderHandler
}

func NewVictoriestClient(ip string, port string, recivedLogic MessageReceivedHandler, sendLogic MessageSenderHandler) *VictoriestClient {
	client := new(VictoriestClient)
	client.ip = ip
	client.port = port
	client.quitSp = make(chan bool)
	client.recivedHandler = recivedLogic
	client.sendHandler = sendLogic
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
	for {
		var msg string
		fmt.Scanln(&msg)

		if strings.EqualFold(msg, "quit") {
			self.quitSp <- true
			break
		}
		self.sendHandler(self, conn, msg)
	}
}

func (self *VictoriestClient) readerPipe(conn *net.TCPConn) {
	reader := bufio.NewReader(conn)
	for {
		message, _, err := deserializeByReader(reader)
		utils.CheckError(err, true)
		self.recivedHandler(self, message)
	}
}

func deserializeByReader(reader *bufio.Reader) (*probe.VictoriestMsg, int32, error) {
	jsonProbe := new(probe.JsonProbe)
	buff, _ := reader.Peek(4)
	data := bytes.NewBuffer(buff)
	var length int32
	err := binary.Read(data, binary.LittleEndian, &length)
	if err != nil {
		log.Error("when deserializeByReader:", err.Error())
		return nil, -1, err
	}

	if int32(reader.Buffered()) < length+4 {
		log.Error("int32(reader.Buffered()) < length + 4")
		return nil, -1, err
	}

	pack := make([]byte, int(4+length))
	_, err = reader.Read(pack)
	if err != nil {
		log.Error("when deserializeByReader:", err.Error())
		return nil, -1, err
	}
	var dst probe.VictoriestMsg
	var msgType int32
	msgType, err = jsonProbe.Deserialize(pack, &dst)
	return &dst, msgType, nil
}
