package main

import (
	"../goconfig"
	"../probe"
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// 退出信号量
var quitSp chan bool

func main() {
	host, port := readServerConfig()
	addr := host + ":" + port
	fmt.Println(addr)
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	checkError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	defer conn.Close()
	fmt.Println("connecting ", conn.RemoteAddr().String(), "...")

	// 退出信号channel
	quitSp = make(chan bool)

	go readerPipe(conn)
	go writerPipe(conn)

	<-quitSp
}

type StrMsg struct {
	Msg string
}

func writerPipe(conn *net.TCPConn) {
	for {
		var msg string
		fmt.Scanln(&msg)

		if strings.EqualFold(msg, "quit") {
			quitSp <- true
			break
		}
		writer := bufio.NewWriter(conn)
		strBuf, _ := probe.Encoding(msg)
		fmt.Println(strBuf)
		writer.Write(strBuf)
		writer.Flush()
	}

}

func readerPipe(conn *net.TCPConn) {
	reader := bufio.NewReader(conn)
	for {
		message, _ := probe.Decoding(reader)
		fmt.Println(message)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
}

func readServerConfig() (string, string) {
	exefile, _ := exec.LookPath(os.Args[0])
	fmt.Println(filepath.Dir(exefile))

	filepath := path.Join(filepath.Dir(exefile), "./client.config")
	cf, err := goconfig.LoadConfigFile(filepath)
	checkError(err)

	host, err := cf.GetValue(goconfig.DEFAULT_SECTION, "server.host")
	checkError(err)

	port, err := cf.GetValue(goconfig.DEFAULT_SECTION, "server.port")
	checkError(err)

	return host, port
}

// func GenStringBuff(str string) *bytes.Buffer {
// 	lenBuf := new(bytes.Buffer)
// 	var length int32 = int32(len(str))
// 	err := binary.Write(lenBuf, binary.LittleEndian, length)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return lenBuf
// 	}
// 	lenBuf.WriteString(str)

// 	return lenBuf
// }
