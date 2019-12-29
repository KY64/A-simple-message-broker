package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
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
	defer conn.Close()

	Request = req.URL.RawQuery
	query := strings.Split(Request, "=")

	var upgrader = websocket.Upgrader{} // use default options

	c, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	var m = make(chan int)
	var msg = make(chan []byte)

	go func() {
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			m <- mt
			msg <- message
			if query[0] == "publish" {
				val, err := publish(conn, "kamar="+query[1], string(message))
				if err != nil {
					log.Println(err)
				}
				if val.(int64) < 1 {
					log.Println("Lonely")
				}
				log.Println("val is", val)
			}
			log.Printf("recv: %s mt: %d", message, mt)
		}
	}()

	for {
		mt := <-m
		if query[0] == "kamar" {
			<-msg
			go subscribe(conn, Request, msg)
		}
		err = c.WriteMessage(mt, <-msg)
		log.Println("yuuhuu")
		if err != nil {
			log.Println("write:", err)
			c.Close()
			break
		}
	}

}
