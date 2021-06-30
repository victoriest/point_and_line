package main

import (
	pb "./victoriest.org/protobuf"
	"bufio"
	"bytes"
	"github.com/golang/protobuf/proto"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

var quitSp chan bool

func main() {
	//strAddr := "127.0.0.1" + ":" + "9596"
	//tcpAddr, _ := net.ResolveTCPAddr("tcp", strAddr)

	//quitSp = make(chan bool)

	//conn, _ := net.DialTCP("tcp", nil, tcpAddr)
	//defer conn.Close()

	//fmt.Println("connecting ", conn.RemoteAddr().String(), "...")

	//go readerPipe(conn)
	//go writerPipe(conn)

	//<-quitSp
}

func writerPipe(conn *net.TCPConn) {
	for {
		var msg string
		fmt.Scanln(&msg)

		if strings.EqualFold(msg, "quit") {
			quitSp <- true
			break
		}
		testMessage := &pb.TestMessage{
			TestInt:    proto.Int32(123),
			TestString: proto.String("est"),
		}
		fmt.Println(testMessage)
		byt, _ := proto.Marshal(testMessage)
		fmt.Println(byt)

		buff := &pb.MobileSuiteProtobuf{
			Type:    proto.Int32(321),
			Arena:   proto.Int32(111),
			Command: proto.Int32(0xa),
			Message: byt,
		}
		fmt.Println(buff)
		bybuf, _ := proto.Marshal(buff)
		fmt.Println(bybuf)
		conn.Write(bybuf)

	}
}

func readerPipe(conn *net.TCPConn) {
	reader := bufio.NewReader(conn)
	for {
		message, _ := deserializeByReader(reader)
		fmt.Println(message)
	}
}

func deserializeByReader(reader *bufio.Reader) (*pb.MobileSuiteProtobuf, error) {
	buff, _ := reader.Peek(4)
	data := bytes.NewBuffer(buff)
	var length int32
	err := binary.Read(data, binary.BigEndian, &length)
	if err != nil {
		return nil, err
	}
	fmt.Println(length)
	if int32(reader.Buffered()) < length+4 {
		return nil, err
	}

	pack := make([]byte, int(4+length))
	_, err = reader.Read(pack)
	if err != nil {
		return nil, err
	}
	msg := pack[4:]
	var dst pb.MobileSuiteProtobuf
	proto.Unmarshal(msg, &dst)
	fmt.Println(&dst)

	var testMsg pb.TestMessage
	proto.Unmarshal(dst.Message, &testMsg)
	fmt.Println(&testMsg)

	return &dst, nil
}
