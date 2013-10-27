package main

import (
	"../goconfig"
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"
)

func main() {
	host, port := readServerConfig()
	addr := host + ":" + port
	fmt.Println(addr)
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	fmt.Println("connecting ", conn.RemoteAddr().String(), "...")
	for {
		writer := bufio.NewWriterSize(conn, 128)
		timeNow := time.Now().String()

		strBuf := GenStringBuff(timeNow)
		writer.Write(strBuf.Bytes())
		writer.Flush()

		reader := bufio.NewReaderSize(conn, 128)
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

		time.Sleep(time.Second * 2)
	}
	conn.Close()
	os.Exit(0)
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

func GenStringBuff(str string) *bytes.Buffer {
	lenBuf := new(bytes.Buffer)
	var length int32 = int32(len(str))
	err := binary.Write(lenBuf, binary.LittleEndian, length)
	if err != nil {
		fmt.Println(err.Error())
		return lenBuf
	}
	lenBuf.WriteString(str)
	return lenBuf
}
