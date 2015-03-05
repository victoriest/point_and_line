package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

var quitSemaphore chan bool

func main() {
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:9999")

	conn, _ := net.DialTCP("tcp", nil, tcpAddr)
	defer conn.Close()
	fmt.Println("connected!")

	go onSendMessage(conn)
	go onMessageRecived(conn)

	<-quitSemaphore
}

func onSendMessage(conn *net.TCPConn) {
	for {
		var msg string
		fmt.Scanln(&msg)

		if strings.EqualFold(msg, "quit") {
			quitSemaphore <- true
			break
		}

		b := []byte(msg)
		conn.Write(b)
	}
}

func onMessageRecived(conn *net.TCPConn) {
	reader := bufio.NewReader(conn)
	for {
		msg, _, _ := reader.ReadLine()
		fmt.Println(string(msg))
	}
}
