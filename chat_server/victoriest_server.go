package main

import (
	"./goconfig"
	"./probe"
	"bufio"
	// "bytes"
	// "encoding/binary"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
)

// 退出信号量
var quitSp chan bool

// 客户端连接Map
var connMap map[string]*net.TCPConn

func main() {
	// 启动服务goroutine
	go startUp()
	shutDown()
}

/**
 * 初始化服务器
 */
func startUp() {
	// 从配置文件中读取port
	strAddr := ":" + readServerPort()

	// 构造tcpAddress
	tcpAddr, err := net.ResolveTCPAddr("tcp", strAddr)
	checkError(err, true)

	// 创建TCP监听
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err, true)
	defer tcpListener.Close()
	fmt.Println("start listen ", tcpListener.Addr().String())

	// tcpConn的map
	connMap = make(map[string]*net.TCPConn)
	// 退出信号channel
	quitSp = make(chan bool)

	// 连接管理
	initConnectionManager(connMap, tcpListener)
}

/**
 * 关闭服务器指令
 */
func shutDown() {
	// 监测退出程序的信号量
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, os.Interrupt, os.Kill)
	<-sign
	fmt.Println(len(connMap))

	// 关闭所有连接
	for _, conn := range connMap {
		fmt.Println("close:", conn.RemoteAddr().String())
		conn.Close()
	}
	fmt.Println("shutdown")
}

/**
 * 一客户端一线程
 */
func tcpHandler(tcpConn net.TCPConn) {
	ipStr := tcpConn.RemoteAddr().String()
	defer func() {
		fmt.Println("disconnected :" + ipStr)
		broadcastMessage("disconnected :" + ipStr)
		tcpConn.Close()
		delete(connMap, ipStr)
		fmt.Println(len(connMap))
	}()
	broadcastMessage("A new connection :" + ipStr)
	reader := bufio.NewReader(&tcpConn)
	for {
		message, err := probe.Decoding(reader)
		if err != nil {
			return
		}
		// message = tcpConn.RemoteAddr().String() + ":" + string(message)
		fmt.Println(message)

		// use pack do what you want ...
		broadcastMessage(message)
	}
}

/**
 * 客户端连接管理器
 */
func initConnectionManager(connMap map[string]*net.TCPConn, tcpListener *net.TCPListener) {

	i := 0
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		connMap[tcpConn.RemoteAddr().String()] = tcpConn
		i++

		fmt.Println("a client connected : ", tcpConn.RemoteAddr().String())
		go tcpHandler(*tcpConn)
	}
}

/**
 * 读取配置文件
 */
func readServerPort() string {
	exefile, _ := exec.LookPath(os.Args[0])
	fmt.Println(filepath.Dir(exefile))
	filepath := path.Join(filepath.Dir(exefile), "./server.config")
	cf, err := goconfig.LoadConfigFile(filepath)
	checkError(err, true)
	port, err := cf.GetValue(goconfig.DEFAULT_SECTION, "server.port")
	checkError(err, true)
	return port
}

func broadcastMessage(message interface{}) {
	buff, _ := probe.Encoding(message)
	// 向所有人发话
	for _, conn := range connMap {
		conn.Write(buff)
	}

}

func boradcastConnectingMessage(conn *net.TCPConn) {
	message := "A new connection :" + conn.RemoteAddr().String()
	broadcastMessage(message)
}

func checkError(err error, isQuit bool) {
	if err != nil {
		fmt.Println(err.Error())
		if isQuit {
			os.Exit(2)
		}
	}
}
