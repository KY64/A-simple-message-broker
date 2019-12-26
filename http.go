package main

import (
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

	go client()

	http.HandleFunc("/account/register", registerHandler)
	http.HandleFunc("/account/auth", loginHandler)
	http.HandleFunc("/api/stream/1", streamHandler)
	PORT := ":" + arg[1]
	log.Fatal(http.ListenAndServe(PORT, nil))
}
