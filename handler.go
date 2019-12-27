package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func registerHandler(w http.ResponseWriter, req *http.Request) {
	pool := newPool()
	conn := pool.Get()
	defer conn.Close()

	if req.Method == http.MethodPost {
		raw, err := ioutil.ReadAll(req.Body)

		if err != nil {
			log.Fatalln(err)
		}

		var User struct {
			Username string `redis:"username"`
			Fullname string `redis:"fullname"`
			Email    string `redis:"email"`
			Password string `redis:"password"`
			Phone    int    `redis:"phone"`
			Address  string `redis:"address"`
		}

		json.Unmarshal(raw, &User)

		if User.Username == "" || User.Password == "" || User.Phone == 0 || User.Address == "" || User.Email == "" {
			w.WriteHeader(400)
			return
		}

		err = hmset(conn, string("reg:"+User.Username), User)

		if err != nil {
			log.Fatalln(err)
		}

		isDone := make(chan bool)
		go expire(conn, User.Username, "60", isDone)

		w.Write([]byte("200"))
		<-isDone
	}
}

func loginHandler(w http.ResponseWriter, req *http.Request) {

	pool := newPool()
	conn := pool.Get()
	defer conn.Close()

	if req.Method == http.MethodPost {
		raw, err := ioutil.ReadAll(req.Body)

		log.Println(string(raw))

		var User struct {
			Username string `redis:"username"`
			Password string `redis:"password"`
		}

		json.Unmarshal(raw, &User)

		if User.Username == "" || User.Password == "" {
			w.WriteHeader(400)
			return
		}

		err = hmset(conn, string("login:"+User.Username), User)

		if err != nil {
			log.Fatalln(err)
		}
		isDone := make(chan bool)
		go expire(conn, User.Username, "60", isDone)

		w.Write([]byte("200"))
		<-isDone
	}
}

func streamHandler(w http.ResponseWriter, req *http.Request) {

	pool := newPool()
	conn := pool.Get()

	if req.Method == http.MethodGet {
		Request = req.RequestURI
		w.Write([]byte("'Sup?"))
	}

	var data = make(chan string)
	go subscribe(conn, "makan", data)
	fmt.Println(<-data)
}
