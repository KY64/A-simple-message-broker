package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

func client() int {
	arguments := os.Args

	if len(arguments) == 2 {
		log.Fatalln("Please provide CLIENT host:port")
		return 1
	}

	c, err := net.Dial("tcp", arguments[2])

	if err != nil {
		log.Println(err)
		return 1
	}

	defer c.Close()
	log.Println("Connected!")
	data, error := bufio.NewReader(c).ReadString('\n')

	if error != nil {
		log.Println(err)
		return 1
	}

	log.Println(data)

	return 0
}
