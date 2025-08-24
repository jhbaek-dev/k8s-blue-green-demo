package main

import (
	"log"
	"net"
)

const VERSION = "v1.1.3"

func main() {
	log.Printf("[DEBUG] Service-C 시작 중... 버전: %s", VERSION)

	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("[DEBUG] TCP 리스너 시작 실패: %v", err)
	}
	defer listener.Close()

	log.Printf("[DEBUG] TCP 서버 시작: 포트 3000에서 대기 중...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("[DEBUG] 클라이언트 연결 수락 실패: %v", err)
			continue
		}

		log.Printf("[DEBUG] 새로운 클라이언트 연결 수락: %s", conn.RemoteAddr())

		go func(c net.Conn) {
			defer func() {
				log.Printf("[DEBUG] 클라이언트 연결 종료: %s", c.RemoteAddr())
				c.Close()
			}()

			buf := make([]byte, 1024)
			n, err := c.Read(buf)
			if err != nil {
				log.Printf("[DEBUG] 클라이언트로부터 데이터 읽기 실패: %v", err)
				return
			}

			request := string(buf[:n])
			log.Printf("[DEBUG] 요청 수신: '%s' from %s", request, c.RemoteAddr())

			_, err = c.Write([]byte(VERSION))
			if err != nil {
				log.Printf("[DEBUG] 클라이언트에게 응답 전송 실패: %v", err)
				return
			}

			log.Printf("[DEBUG] 버전 응답 전송 성공: %s to %s", VERSION, c.RemoteAddr())
		}(conn)
	}
}
