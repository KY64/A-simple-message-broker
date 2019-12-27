package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

func newPool() *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle:     5,
		IdleTimeout: 60 * time.Second,
		// max number of connections
		MaxActive: 10,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		// Other pool configuration not shown in this example.
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				log.Fatal(err.Error())
			}
			return c, err
		},
	}
}

func ping(c redis.Conn) error {
	// Send PING command to Redis
	pong, err := c.Do("PING")
	if err != nil {
		return err
	}

	s, err := redis.String(pong, err)
	if err != nil {
		return err
	}

	fmt.Printf("PING Response = %s\n", s)

	return nil
}

func set(c redis.Conn, key string, data string) error {
	_, err := c.Do("SET", key, data)
	return err
}

func hmset(c redis.Conn, hash string, Class interface{}) error {
	_, err := c.Do("HMSET", redis.Args{}.Add(hash).AddFlat(Class)...)
	return err
}

func expire(c redis.Conn, key string, expiredTime string, done chan bool) {
	t, _ := strconv.ParseInt(expiredTime, 0, 64)
	for {
		select {
		case <-time.After(time.Duration(t) * time.Second):
			_, err := c.Do("EXPIRE", key, expiredTime)
			if DBStatus > 0 {
				done <- true
				return
			}
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}

func subscribe(c redis.Conn, channel string, data chan string) {
	psc := redis.PubSubConn{c}

	psc.Subscribe(channel)

	err := c.Err()

	for err == nil {
		switch v := psc.Receive().(type) {
		case redis.Message:
			data <- string(v.Data)
		case error:
			log.Println(v.Error())
		}
	}

	psc.Unsubscribe()
	c.Close()
}
