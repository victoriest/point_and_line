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

func main() {
	strAddr := ":" + readServerConfig()
	tcpAddr, err := net.ResolveTCPAddr("tcp", strAddr)
	checkError(err)
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	defer tcpListener.Close()
	fmt.Println("start listen ", tcpListener.Addr().String())
	connChan := make(chan *net.TCPConn)
	connMap := make(map[int]*net.TCPConn)
	go ConnectionManager(connMap, connChan)
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

func ConnectionManager(connMap map[int]*net.TCPConn, connChan chan *net.TCPConn) {
	i := 0
	for {
		newConn := <-connChan
		connMap[i] = newConn
		i++
	}
}

/**
 * 读取配置文件
 */
func readServerConfig() string {
	exefile, _ := exec.LookPath(os.Args[0])
	// exepath, _ := filepath.Abs(exefile)
	// dir, _ := path.Split(exepath)
	// wd, _ := os.Getwd()
	// _, wd, _, _ := runtime.Caller(0)

	fmt.Println(filepath.Dir(exefile))
	filepath := path.Join(filepath.Dir(exefile), "./server_config.ini")
	cf, err := goconfig.LoadConfigFile(filepath)
	checkError(err)
	port, err := cf.GetValue(goconfig.DEFAULT_SECTION, "server.port")
	checkError(err)
	return port
}
