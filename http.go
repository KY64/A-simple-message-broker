package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	arg := os.Args
	if len(arg) == 1 {
		fmt.Println("Please provide SERVER port number")
		return
	}

	flag.Parse()
	go h.run()
	go client()

	http.HandleFunc("/account/register", registerHandler)
	http.HandleFunc("/account/auth", loginHandler)
	go http.HandleFunc("/api/stream", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r)
	})

	PORT := ":" + arg[1]
	log.Fatal(http.ListenAndServe(PORT, nil))
}
