package server

import (
	"../codec"
	"../dao"
	"../protocol"
	"../utils"
	"bufio"
	log "code.google.com/p/log4go"
	"net"
)

// 服务器接口
type INexus interface {
	// 启动服务器
	Startup()
	// 关闭服务器
	Shutdown()
	// 重启服务器
	// Restart()
}

// 消息处理托管
type MessageRecivedHandler func(*Nexus, string, *protocol.MobileSuiteModel)

// 连接状态处理托管
type ConnectionHandler func(*Nexus, *net.TCPConn)

type Nexus struct {
	port                 string                  // 服务端端口号
	quitSemaphore        chan bool               // 退出信号量
	connMap              map[string]*net.TCPConn // 客户端连接Map
	recivedHandler       MessageRecivedHandler   // 消息逻辑处理托管Handler
	newConnectionHandler ConnectionHandler       // 新连接处理Handler
	disconnectHandler    ConnectionHandler       // 断开连接处理Handler
	probe                codec.ProtobufProbe     // 序列化接口
	DbConnector          *dao.MysqlConnector
}

func NewNexus(port string, handler MessageRecivedHandler, connHander ConnectionHandler, disconnHander ConnectionHandler, dbCon *dao.MysqlConnector) *Nexus {
	nexus := new(Nexus)
	nexus.port = port
	nexus.connMap = make(map[string]*net.TCPConn)
	nexus.quitSemaphore = make(chan bool)
	nexus.recivedHandler = handler
	nexus.newConnectionHandler = connHander
	nexus.disconnectHandler = disconnHander
	nexus.probe = *new(codec.ProtobufProbe)
	nexus.DbConnector = dbCon
	return nexus
}

/**
 * 客户端连接管理器
 */
func (self *Nexus) initConnectionManager(tcpListener *net.TCPListener) {
	i := 0
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			log.Error(err.Error())
			continue
		}

		self.connMap[tcpConn.RemoteAddr().String()] = tcpConn
		i++

		log.Info("A client connected : ", tcpConn.RemoteAddr().String())
		go self.tcpPipe(tcpConn)
	}
}

func (self *Nexus) Startup() {
	strAddr := ":" + self.port

	tcpAddr, err := net.ResolveTCPAddr("tcp", strAddr)
	utils.CheckError(err, true)

	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	utils.CheckError(err, true)

	defer tcpListener.Close()
	log.Info("Start listen ", tcpListener.Addr().String())

	// 连接管理
	self.initConnectionManager(tcpListener)
}

/**
 * 关闭服务器指令
 */
func (self *Nexus) Shutdown() {
	// 关闭所有连接
	for _, conn := range self.connMap {
		log.Info("close:" + conn.RemoteAddr().String())
		conn.Close()
	}
	log.Info("Shutdown")
}

/**
 * 一客户端一线程
 */
func (self *Nexus) tcpPipe(tcpConn *net.TCPConn) {
	ipStr := tcpConn.RemoteAddr().String()
	defer func() {
		log.Info("disconnected :" + ipStr)
		self.disconnectHandler(self, tcpConn)

		tcpConn.Close()
		delete(self.connMap, ipStr)
	}()
	self.newConnectionHandler(self, tcpConn)

	reader := bufio.NewReader(tcpConn)

	for {
		message, _, err := self.probe.DeserializeByReader(reader)
		if err != nil {
			return
		}
		self.recivedHandler(self, ipStr, message)
	}
}

/**
 * 全局广播
 */
func (self *Nexus) BroadcastMessage(message *protocol.MobileSuiteModel) {
	buff, _ := self.probe.Serialize(message)
	// 向所有人发话
	for _, conn := range self.connMap {
		conn.Write(buff)
	}
}

/**
 * 向某人发消息
 */
func (self *Nexus) SendTo(sendTo string, message *protocol.MobileSuiteModel) {
	buff, _ := self.probe.Serialize(message)
	self.connMap[sendTo].Write(buff)
}
