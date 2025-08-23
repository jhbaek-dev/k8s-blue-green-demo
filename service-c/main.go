package main

import (
	"net"
)

const VERSION = "v1.1.0"

func main() {
	listener, _ := net.Listen("tcp", ":3000")
	defer listener.Close()
	
	for {
		conn, _ := listener.Accept()
		go func(c net.Conn) {
			defer c.Close()
			buf := make([]byte, 1024)
			c.Read(buf) // 요청을 읽고
			c.Write([]byte(VERSION)) // 버전을 응답
		}(conn)
	}
}