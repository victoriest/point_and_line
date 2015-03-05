package main

import (
	"fmt"
	"net"
)

func main() {
	var tcpAddr *net.TCPAddr
	// var err error

	tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:9999")

	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)

	defer tcpListener.Close()

	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			continue
		}

		fmt.Println(tcpConn.RemoteAddr().String())

	}
}
