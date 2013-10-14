package main

import (
	"./goconfig"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

var quitSp chan bool

func main() {
	argLen := len(os.Args)
	var cmd string
	if argLen != 1 {
		cmd = "startup"
	} else {
		cmd = os.Args[1]
	}

	if cmd == "startup" {
		startUp()
	} else if cmd == "shutdown" {
		shutDown()
	}

}

func startUp() {
	strAddr := ":" + readServerConfig()
	tcpAddr, err := net.ResolveTCPAddr("tcp", strAddr)
	checkError(err)
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	defer tcpListener.Close()
	fmt.Println("start listen ", tcpListener.Addr().String())
	// 新连接channel
	connChan := make(chan *net.TCPConn)
	// tcpConn的map
	connMap := make(map[int]*net.TCPConn)
	// 退出信号channel
	quitSp = make(chan bool)
	go ConnectionManager(connMap, connChan, quitSp)
	for {
		var tcpConn *net.TCPConn
		tcpConn, err = tcpListener.AcceptTCP()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		connChan <- tcpConn

		fmt.Println("a client connected : ", tcpConn.RemoteAddr().String())
		go tcpHandler(*tcpConn)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
}

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
func ConnectionManager(connMap map[int]*net.TCPConn, connChan chan *net.TCPConn, quitSp chan bool) {
	// TODO : 弄个uuid呗
	i := 0
	for {
		select {
		case newConn := <-connChan:
			connMap[i] = newConn
			i++
		case <-quitSp:
			for _, conn := range connMap {
				conn.Close()
			}
			break
		}
	}
}

/**
 * 读取配置文件
 */
func readServerConfig() string {
	exefile, _ := exec.LookPath(os.Args[0])
	// exepath, _ := filepath.Abs(exefile)
	// dir, _ := path.Split(exepath)

	fmt.Println(filepath.Dir(exefile))
	filepath := path.Join(filepath.Dir(exefile), "./server_config.ini")
	cf, err := goconfig.LoadConfigFile(filepath)
	checkError(err)
	port, err := cf.GetValue(goconfig.DEFAULT_SECTION, "server.port")
	checkError(err)
	return port
}

/**
* 关闭服务器指令
 */
func shutDown() {
	quitSp <- true
}
