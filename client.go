package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var DBStatus int
var Respond, Request string

func client() {
	arguments := os.Args

	if len(arguments) == 2 {
		log.Fatalln("Please provide CLIENT host:port")
		return
	}

	c, err := net.Dial("tcp", arguments[2])

	if err != nil {
		log.Println(err)
		DBStatus = 0
		go client()
		return
	}

	DBStatus = 1
	log.Println("Connected!")

	defer c.Close()

	for {
		select {
		case <-time.After(60 * time.Millisecond):
			fmt.Println("")

		default:
			if len(Request) > 0 {
				data, error := bufio.NewReader(c).ReadString('\n')
				log.Println(Request)

				if error != nil {
					log.Println(err)
					DBStatus = 0
					go client()
					return
				}

				Respond = data
				Request = ""
			} else {
				_, err = bufio.NewReader(c).ReadString('\n')

				if err != nil {
					log.Println(err)
					DBStatus = 0
					go client()
					return
				}
			}
			// log.Println(data)
		}
	}
}
