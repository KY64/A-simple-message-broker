package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func registerHandler(w http.ResponseWriter, req *http.Request) {
	pool := newPool()
	conn := pool.Get()
	defer conn.Close()

	raw, err := ioutil.ReadAll(req.Body)

	if err != nil {
		panic(err)
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

	err = hmset(conn, string(User.Username), User)

	if err != nil {
		log.Fatalln(err)
	}
}

func loginHandler(w http.ResponseWriter, req *http.Request) {

	pool := newPool()
	conn := pool.Get()
	defer conn.Close()

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

	err = hmset(conn, string(User.Username), User)

	if err != nil {
		log.Fatalln(err)
	}
	isDone := make(chan bool)
	go expire(conn, User.Username, "2", isDone)

	<-isDone

	w.Write([]byte("200"))
}

func streamHandler(w http.ResponseWriter, req *http.Request) {

	pool := newPool()
	conn := pool.Get()

	_, err := subscribe(conn, "makan")

	if err != nil {
		panic(err)
	}

}
