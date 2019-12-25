package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var DBStatus int = 0

func main() {
	arg := os.Args
	if len(arg) == 1 {
		fmt.Println("Please provide SERVER port number")
		return
	}

	DBStatus := client()
	if DBStatus > 0 {
		log.Println("Server not found")
		check := time.Tick(100 * time.Millisecond)
		go func() {
			for {
				select {
				case <-check:
					DBStatus = client()
					if DBStatus < 1 {
						return
					}
				default:
				}
			}
		}()
	}

	http.HandleFunc("/account/register", registerHandler)
	http.HandleFunc("/account/auth", loginHandler)
	http.HandleFunc("/api/stream/1", streamHandler)
	PORT := ":" + arg[1]
	log.Fatal(http.ListenAndServe(PORT, nil))
}
