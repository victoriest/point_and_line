package main

import (
	pb "./protobuf"
	proto "code.google.com/p/goprotobuf/proto"
	"fmt"
	"net"
	"strings"
)

var quitSp chan bool

func main() {
	strAddr := "127.0.0.1" + ":" + "9596"
	tcpAddr, _ := net.ResolveTCPAddr("tcp", strAddr)

	quitSp = make(chan bool)

	conn, _ := net.DialTCP("tcp", nil, tcpAddr)
	defer conn.Close()

	fmt.Println("connecting ", conn.RemoteAddr().String(), "...")

	go writerPipe(conn)
	// go readerPipe(conn)

	<-quitSp
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

// func readerPipe(conn *net.TCPConn) {
// 	reader := bufio.NewReader(conn)
// 	for {
// 		message, _, err := deserializeByReader(reader)
// 		utils.CheckError(err, true)
// 		self.recivedHandler(self, message)
// 	}
// }
