package client

import (
	"bufio"
	"fmt"
	"go_server/codec"
	"go_server/log"
	"go_server/protocol"
	"go_server/utils"
	"net"
	"strings"
)

// 消息处理托管
type MessageReceivedHandler func(*RobotClient, *protocol.MobileSuiteModel)
type MessageSenderHandler func(*RobotClient, *net.TCPConn, string)

type IClient interface {
	Startup()
}

type RobotClient struct {
	port            string                 // 服务端端口号
	ip              string                 // IP
	quitSemaphore   chan bool              // 退出信号量
	receivedHandler MessageReceivedHandler // 消息接收逻辑处理托管
	sendHandler     MessageSenderHandler   // 消息发送逻辑处理托管
	Probe           codec.ProtobufProbe    // 序列化实现
}

func NewClient(ip string, port string, receivedLogic MessageReceivedHandler, sendLogic MessageSenderHandler) *RobotClient {
	client := new(RobotClient)
	client.port = port
	client.ip = ip
	client.quitSemaphore = make(chan bool)
	client.receivedHandler = receivedLogic
	client.sendHandler = sendLogic
	client.Probe = *new(codec.ProtobufProbe)
	return client
}

func (self *RobotClient) Startup() {
	strAddr := self.ip + ":" + self.port
	tcpAddr, err := net.ResolveTCPAddr("tcp", strAddr)
	utils.CheckError(err, true)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	utils.CheckError(err, true)
	defer conn.Close()
	log.Info("connecting ", conn.RemoteAddr().String(), "...")

	go self.onSendMessage(conn)
	go self.onMessageReceived(conn)

	<-self.quitSemaphore
}

func (self *RobotClient) onSendMessage(conn *net.TCPConn) {
	for {
		var msg string
		fmt.Scanln(&msg)

		if strings.EqualFold(msg, "quit") {
			self.quitSemaphore <- true
			break
		}
		self.sendHandler(self, conn, msg)
	}
}

func (self *RobotClient) onMessageReceived(conn *net.TCPConn) {
	reader := bufio.NewReader(conn)
	for {
		message, _, err := self.Probe.DeserializeByReader(reader)
		utils.CheckError(err, true)
		self.receivedHandler(self, message)
	}
}
