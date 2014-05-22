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

// 服务器接口
type IVictoriestServer interface {
	Startup()
	Shutdown()
}

// 消息逻辑处理托管Handler
type MessageReceivedHandler func(*VictoriestServer, *probe.VictoriestMsg)

// 连接状态处理托管Handler
type ConnectionHandler func(*VictoriestServer, *net.TCPConn)

type VictoriestServer struct {
	// 服务端端口号
	port string
	// 退出信号量
	quitSp chan bool
	// 客户端连接Map
	connMap map[string]*net.TCPConn
	// 消息逻辑处理托管Handler
	recivedHandler MessageReceivedHandler
	// 新连接处理Handler
	connectedHandler ConnectionHandler
	// 断开连接处理Handler
	disconnectingHandler ConnectionHandler
}

// Server的构造器
func NewVictoriestServer(port string, handler MessageReceivedHandler, connHander ConnectionHandler, disconnHander ConnectionHandler) *VictoriestServer {
	server := new(VictoriestServer)
	server.port = port
	// tcpConn的map
	server.connMap = make(map[string]*net.TCPConn)
	// 退出信号channel
	server.quitSp = make(chan bool)

	server.recivedHandler = handler
	server.connectedHandler = connHander
	server.disconnectingHandler = disconnHander

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
		go self.tcpPipe(tcpConn)
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
func (self *VictoriestServer) tcpPipe(tcpConn *net.TCPConn) {
	ipStr := tcpConn.RemoteAddr().String()
	defer func() {
		log.Debug("disconnected :" + ipStr)
		self.disconnectingHandler(self, tcpConn)

		tcpConn.Close()
		delete(self.connMap, ipStr)
	}()
	self.connectedHandler(self, tcpConn)

	reader := bufio.NewReader(tcpConn)

	for {
		// var b []byte
		// tcpConn.Read(b)
		// jsonProbe := new(probe.JsonProbe)
		// var dst probe.VictoriestMsg
		// _, err := jsonProbe.Deserialize(b, &dst)
		// log.Debug(dst)
		message, _, err := deserializeByReader(reader)
		if err != nil {
			return
		}
		// use pack do what you want ...
		self.recivedHandler(self, message)
	}
}

/**
 * 全局广播
 */
func (self *VictoriestServer) BroadcastMessage(message *probe.VictoriestMsg) {
	jsonProbe := new(probe.JsonProbe)
	buff, _ := jsonProbe.Serialize(message)
	// 向所有人发话
	for _, conn := range self.connMap {
		conn.Write(buff)
	}
}

/**
 * 向某人发消息
 */
func (self *VictoriestServer) SendTo(sendTo string, message *probe.VictoriestMsg) {
	jsonProbe := new(probe.JsonProbe)
	buff, _ := jsonProbe.Serialize(message)
	self.connMap[sendTo].Write(buff)
}

func deserializeByReader(reader *bufio.Reader) (*probe.VictoriestMsg, int32, error) {
	jsonProbe := new(probe.JsonProbe)
	lengthByte, _ := reader.Peek(4)
	lengthBuff := bytes.NewBuffer(lengthByte)
	var length int32
	err := binary.Read(lengthBuff, binary.LittleEndian, &length)
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
	log.Debug(length, msgType, dst)
	return &dst, msgType, nil
}
