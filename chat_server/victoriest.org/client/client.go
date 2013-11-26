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
		testMsg := probe.TestMsg{MsgInt: 5, ChatMessage: msg}
		broMsg := probe.VictoriestMsg{MsgType: probe.MSG_TYPE_TEST_MESSGAE, MsgContext: testMsg}

		strBuf, _ := jsonProbe.Serialize(broMsg, probe.MSG_TYPE_TEST_MESSGAE)
		writer.Write(strBuf)
		writer.Flush()
	}
}

func (self *VictoriestClient) readerPipe(conn *net.TCPConn) {
	reader := bufio.NewReader(conn)
	for {
		message, _, err := deserializeByReader(reader)
		log.Debug(message)
		// switch obj := (interface{}(message)).(type) {
		// case probe.VictoriestMsg:
		// 	log.Debug(obj.MsgContext)
		// default:
		// 	log.Debug("not a VictoriestMsg  ", message, "  ", msgType)
		// }
		// log.Debug((probe.VictoriestMsg(message)).MsgContext)
		utils.CheckError(err, true)
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

	if int32(reader.Buffered()) < length+8 {
		log.Error("int32(reader.Buffered()) < length + 8")
		return nil, -1, err
	}

	pack := make([]byte, int(8+length))
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
