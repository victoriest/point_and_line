package main

import (
	"bufio"
	"fmt"
	"go_server/demo/codec"
	"net"
)

var quitSemaphore chan bool

func main() {
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:9999")

	conn, _ := net.DialTCP("tcp", nil, tcpAddr)
	defer conn.Close()
	fmt.Println("connected!")

	go onMessageRecived(conn)
	for {
		var msg string
		fmt.Scanln(&msg)
		b, _ := codec.Encode(msg)
		conn.Write(b)
	}
	// b := []byte("time\n")
	// conn.Write(b)

	// <-quitSemaphore
}

func onMessageRecived(conn *net.TCPConn) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := codec.Decode(reader)
		fmt.Println(msg)
		if err != nil {
			quitSemaphore <- true
			break
		}
	}
}
