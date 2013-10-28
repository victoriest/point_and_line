package main

import (
	"./goconfig"
	"bufio"
	"bytes"
	"encoding/binary"
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
var connMap map[int]*net.TCPConn

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
	checkError(err)

	// 创建TCP监听
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	defer tcpListener.Close()
	fmt.Println("start listen ", tcpListener.Addr().String())

	// tcpConn的map
	connMap = make(map[int]*net.TCPConn)
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

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
}

/**
 * 一客户端一线程
 */
func tcpHandler(tcpConn net.TCPConn) {
	defer tcpConn.Close()
	for {
		reader := bufio.NewReaderSize(&tcpConn, 128)
		buff, _ := reader.Peek(4)
		data := bytes.NewBuffer(buff)
		var length int32
		err := binary.Read(data, binary.LittleEndian, &length)
		checkError(err)
		fmt.Println(length)
		if int32(reader.Buffered()) < length+4 {
			fmt.Println("int32(reader.Buffered()) < length+4")
			_, err := reader.Peek(int(4 + length))
			if err != nil {
				return
			}
		}
		pack := make([]byte, int(4+length))
		_, err = reader.Read(pack)
		if err != nil {
			break
		}
		message := string(pack[4:])
		message = tcpConn.RemoteAddr().String() + ":" + message
		fmt.Println(message)

		// use pack do what you want ...
		broadcastMessage(message)
	}
}

/**
 * 客户端连接管理器
 */
func initConnectionManager(connMap map[int]*net.TCPConn, tcpListener *net.TCPListener) {

	i := 0
	for {
		// var tcpConn *net.TCPConn
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		// connChan <- tcpConn

		connMap[i] = tcpConn
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
	checkError(err)
	port, err := cf.GetValue(goconfig.DEFAULT_SECTION, "server.port")
	checkError(err)
	return port
}

func broadcastMessage(message string) {
	msgLength := int32(len(message))
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.LittleEndian, msgLength)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	buff.WriteString(message)

	// 向所有人发话
	for _, conn := range connMap {
		conn.Write(buff.Bytes())
	}

}
