package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

const VERSION = "v1.0.1"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		serviceBHost := os.Getenv("SERVICE_B_HOST")
		serviceCHost := os.Getenv("SERVICE_C_HOST")
		
		serviceBVersion := "unknown"
		if serviceBHost != "" {
			resp, err := http.Get("http://" + serviceBHost + ":3000/version")
			if err == nil {
				body, _ := io.ReadAll(resp.Body)
				serviceBVersion = string(body)
				resp.Body.Close()
			}
		}
		
		serviceCVersion := "unknown"
		if serviceCHost != "" {
			conn, err := net.Dial("tcp", serviceCHost+":3000")
			if err == nil {
				conn.Write([]byte("version"))
				buf := make([]byte, 1024)
				n, _ := conn.Read(buf)
				serviceCVersion = string(buf[:n])
				conn.Close()
			}
		}
		
		fmt.Fprintf(w, "Service-A Version: %s<br>Service-B Version: %s<br>Service-C Version: %s", 
			VERSION, serviceBVersion, serviceCVersion)
	})
	
	http.ListenAndServe(":3000", nil)
}