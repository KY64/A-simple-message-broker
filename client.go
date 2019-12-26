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

func client() {
	arguments := os.Args

	if len(arguments) == 2 {
		log.Fatalln("Please provide CLIENT host:port")
		return
	}

	c, err := net.Dial("tcp", arguments[2])

	if err != nil {
		// log.Println(err)
		DBStatus = 0
		go client()
		return
	}

	DBStatus = 1
	log.Println("Connected!")

	defer c.Close()

	for {
		select {
		case <-time.After(1 * time.Second):
			fmt.Println("")

		default:
			// data, error := bufio.NewReader(c).ReadString('\n')
			_, error := bufio.NewReader(c).ReadString('\n')
			if error != nil {
				// log.Println(err)
				DBStatus = 0
				go client()
				return
			}
			// log.Println(data)
		}
	}
}
