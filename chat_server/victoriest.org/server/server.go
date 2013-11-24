package server

import (
	"../probe"
	"../utils"
	"bufio"
	log "code.google.com/p/log4go"
	"net"
)

type IVictoriestServer interface {
	Startup()
	Shutdown()
}

type VictoriestServer struct {
	// 服务端端口号
	port string
	// 退出信号量
	quitSp chan bool
	// 客户端连接Map
	connMap map[string]*net.TCPConn
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
		self.broadcastMessage("disconnected :" + ipStr)
		tcpConn.Close()
		delete(self.connMap, ipStr)
	}()
	self.broadcastMessage("A new connection :" + ipStr)
	reader := bufio.NewReader(&tcpConn)
	for {
		jsonProbe := new(probe.JsonProbe)
		var message interface{}
		err := jsonProbe.DeserializeByReader(reader, message)
		if err != nil {
			return
		}

		// switch obj := (interface{}(message)).(type) {
		// case probe.VictoriestMsg:
		// 	log.Debug("est", obj.MsgContext)
		// default:
		// 	log.Debug("not a VictoriestMsg  ", obj)
		// }

		// log.Debug(message.(probe.VictoriestMsg).MegContext)

		// use pack do what you want ...
		self.broadcastMessage(message)
	}
}

func (self *VictoriestServer) broadcastMessage(message interface{}) {
	jsonProbe := new(probe.JsonProbe)
	buff, _ := jsonProbe.Serialize(message)
	// 向所有人发话
	for _, conn := range self.connMap {
		conn.Write(buff)
	}
}
