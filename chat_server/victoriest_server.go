package main

import (
	"./goconfig"
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
	startUp()

	// 监测退出程序的信号量
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	s := <-c
	// shutDown()
	fmt.Println("Got signal:", s)

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
	quitSp <- true
	for _, conn := range connMap {
		fmt.Println("close")
		conn.Close()
	}
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
		buf := make([]byte, 256)
		n, err := tcpConn.Read(buf)
		if err != nil {
			break
		}
		// checkError(err)
		if n > 0 {
			fmt.Println(tcpConn.RemoteAddr().String(), " read ", n, "byte :", string(buf))
		} else {
			fmt.Println(tcpConn.RemoteAddr().String(), "nothing readed...")
		}
		n, err = tcpConn.Write(buf)
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
