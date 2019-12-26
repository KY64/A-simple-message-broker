package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var LoginSession = make(chan string)
var RegisterSession = make(chan string)

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

	if DBStatus < 1 {
		fmt.Println("DBStatus", DBStatus)
	}

	err = hmset(conn, string(User.Username), User)

	if err != nil {
		log.Fatalln(err)
	}

	err = expire(conn, string(User.Username), "70")

	if err != nil {
		log.Fatalln(err)
	}

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
