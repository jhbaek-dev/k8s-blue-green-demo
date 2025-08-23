package main

import (
	"fmt"
	"net/http"
)

const VERSION = "v1.0.2"

func main() {
	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, VERSION)
	})

	http.ListenAndServe(":3000", nil)
}