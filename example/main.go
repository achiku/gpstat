package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	stats "github.com/fukata/golang-stats-api-handler"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/stats", stats.Handler)
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		log.Print("GET: /hello")
		fmt.Fprintf(w, "hello")
		return
	})
	mux.HandleFunc("/sleep", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 1)
		log.Print("GET: /sleep")
		fmt.Fprintf(w, "just woke up")
		return
	})
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
