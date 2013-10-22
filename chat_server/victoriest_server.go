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

func startUp() {
	// 从配置文件中读取port
	strAddr := ":" + readServerConfig()

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

	// // 连接管理
	ConnectionManager(connMap, quitSp, tcpListener)
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
		fmt.Println(string(pack[4:]))

		// buf := make([]byte, 256)
		// n, err := tcpConn.Read(buf)
		// if err != nil {
		// 	break
		// }
		// // checkError(err)
		// if n > 0 {
		// 	fmt.Println(tcpConn.RemoteAddr().String(), " read ", n, "byte :", string(buf))
		// } else {
		// 	fmt.Println(tcpConn.RemoteAddr().String(), "nothing readed...")
		// }
		_, err = tcpConn.Write(pack)
		// checkError(err)
		if err != nil {
			break
		}
	}
}

/**
* 客户端连接管理器
 */
func ConnectionManager(connMap map[int]*net.TCPConn, quitSp chan bool, tcpListener *net.TCPListener) {

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
func readServerConfig() string {
	exefile, _ := exec.LookPath(os.Args[0])
	fmt.Println(filepath.Dir(exefile))
	filepath := path.Join(filepath.Dir(exefile), "./server_config.ini")
	cf, err := goconfig.LoadConfigFile(filepath)
	checkError(err)
	port, err := cf.GetValue(goconfig.DEFAULT_SECTION, "server.port")
	checkError(err)
	return port
}
