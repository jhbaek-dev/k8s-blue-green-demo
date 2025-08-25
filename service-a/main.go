package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

const VERSION = "v1.0.3"

func main() {
	log.Printf("[DEBUG] Service-A 시작 중... 버전: %s", VERSION)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[DEBUG] 요청 수신: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		w.Header().Set("Content-Type", "text/html")
		rootVersion := os.Getenv("ROOT_VERSION")
		serviceBHost := os.Getenv("SERVICE_B_HOST")
		serviceCHost := os.Getenv("SERVICE_C_HOST")

		log.Printf("[DEBUG] 환경 변수 확인 - ROOT_VERSION: %s, SERVICE_B_HOST: %s, SERVICE_C_HOST: %s", rootVersion, serviceBHost, serviceCHost)

		serviceBVersion := "unknown"
		if serviceBHost != "" {
			log.Printf("[DEBUG] Service-B에 연결 시도: %s:3000", serviceBHost)
			resp, err := http.Get("http://" + serviceBHost + ":3000/version")
			if err == nil {
				body, _ := io.ReadAll(resp.Body)
				serviceBVersion = string(body)
				resp.Body.Close()
				log.Printf("[DEBUG] Service-B 버전 응답 성공: %s", serviceBVersion)
			} else {
				log.Printf("[DEBUG] Service-B 연결 실패: %v", err)
			}
		} else {
			log.Printf("[DEBUG] SERVICE_B_HOST가 설정되지 않음")
		}

		serviceCVersion := "unknown"
		if serviceCHost != "" {
			log.Printf("[DEBUG] Service-C에 TCP 연결 시도: %s:3000", serviceCHost)
			conn, err := net.Dial("tcp", serviceCHost+":3000")
			if err == nil {
				log.Printf("[DEBUG] Service-C TCP 연결 성공, 버전 요청 전송")
				log.Printf("[DEBUG] TCP 연결 정보 - Local: %s, Remote: %s", conn.LocalAddr(), conn.RemoteAddr())
				log.Printf("[DEBUG] 연결 타입: %T", conn)
				
				log.Printf("[DEBUG] Service-C에 'version' 메시지 전송 시도")
				_, writeErr := conn.Write([]byte("version"))
				if writeErr != nil {
					log.Printf("[DEBUG] Service-C 메시지 전송 실패: %v", writeErr)
				} else {
					log.Printf("[DEBUG] Service-C 메시지 전송 완료")
				}
				
				log.Printf("[DEBUG] Service-C 응답 읽기 시도...")
				buf := make([]byte, 1024)
				n, _ := conn.Read(buf)
				serviceCVersion = string(buf[:n])
				conn.Close()
				log.Printf("[DEBUG] Service-C 버전 응답 성공: %s", serviceCVersion)
			} else {
				log.Printf("[DEBUG] Service-C TCP 연결 실패: %v", err)
			}
		} else {
			log.Printf("[DEBUG] SERVICE_C_HOST가 설정되지 않음")
		}

		response := fmt.Sprintf("Service-A Version: %s<br>Service-B Version: %s<br>Service-C Version: %s<br>Root Version: %s",
			VERSION, serviceBVersion, serviceCVersion, rootVersion)

		log.Printf("[DEBUG] 응답 반환: %s", response)
		fmt.Fprint(w, response)
	})

	log.Printf("[DEBUG] HTTP 서버 시작: 포트 3000에서 대기 중...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
