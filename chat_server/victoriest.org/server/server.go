package server

import (
	"../probe"
	"../utils"
	"bufio"
	"bytes"
	log "code.google.com/p/log4go"
	"encoding/binary"
	"net"
)

type IVictoriestServer interface {
	Startup()
	Shutdown()
}

type ServerHandler func(*VictoriestServer, *probe.VictoriestMsg)

type VictoriestServer struct {
	// 服务端端口号
	port string
	// 退出信号量
	quitSp chan bool
	// 客户端连接Map
	connMap map[string]*net.TCPConn
	//
	Handler ServerHandler
}

func NewVictoriestServer(port string) *VictoriestServer {
	server := new(VictoriestServer)
	server.port = port
	// tcpConn的map
	server.connMap = make(map[string]*net.TCPConn)
	// 退出信号channel
	server.quitSp = make(chan bool)
	return server
}

/**
 * 客户端连接管理器
 */
func (self *VictoriestServer) initConnectionManager(tcpListener *net.TCPListener) {
	i := 0
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			log.Error(err.Error())
			continue
		}

		self.connMap[tcpConn.RemoteAddr().String()] = tcpConn
		i++

		log.Debug("A client connected : ", tcpConn.RemoteAddr().String())
		go self.tcpHandler(*tcpConn)
	}
}

/**
 * 开启服务器
 */
func (self *VictoriestServer) Startup() {
	strAddr := ":" + self.port

	// 构造tcpAddress
	tcpAddr, err := net.ResolveTCPAddr("tcp", strAddr)
	utils.CheckError(err, true)

	// 创建TCP监听
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	utils.CheckError(err, true)
	defer tcpListener.Close()
	log.Debug("Start listen ", tcpListener.Addr().String())

	// 连接管理
	self.initConnectionManager(tcpListener)
}

/**
 * 关闭服务器指令
 */
func (self *VictoriestServer) Shutdown() {
	// 关闭所有连接
	for _, conn := range self.connMap {
		log.Debug("close:" + conn.RemoteAddr().String())
		conn.Close()
	}
	log.Debug("Shutdown")
}

/**
 * 一客户端一线程
 */
func (self *VictoriestServer) tcpHandler(tcpConn net.TCPConn) {
	ipStr := tcpConn.RemoteAddr().String()
	defer func() {
		log.Debug("disconnected :" + ipStr)
		testMsg2 := probe.TestMsg{MsgInt: 5, ChatMessage: "disconnected :" + ipStr}
		broMsg2 := probe.VictoriestMsg{MsgType: probe.MSG_TYPE_TEST_MESSGAE, MsgContext: testMsg2}
		self.BroadcastMessage(broMsg2)
		tcpConn.Close()
		delete(self.connMap, ipStr)
	}()
	testMsg1 := probe.TestMsg{MsgInt: 5, ChatMessage: "A new connection :" + ipStr}
	broMsg1 := probe.VictoriestMsg{MsgType: probe.MSG_TYPE_TEST_MESSGAE, MsgContext: testMsg1}
	self.BroadcastMessage(broMsg1)
	reader := bufio.NewReader(&tcpConn)
	for {
		message, _, err := deserializeByReader(reader)
		if err != nil {
			return
		}
		// use pack do what you want ...
		self.Handler(self, message)
	}
}

func (self *VictoriestServer) BroadcastMessage(message probe.VictoriestMsg) {
	jsonProbe := new(probe.JsonProbe)
	buff, _ := jsonProbe.Serialize(message, probe.MSG_TYPE_TEST_MESSGAE)
	// 向所有人发话
	for _, conn := range self.connMap {
		conn.Write(buff)
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
