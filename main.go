package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}
var scores = map[string]int64{}

func main() {
	address := os.Getenv("ADDRESS")
	if address == "" {
		address = "0.0.0.0:3444"
	}

	server, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	go listen(server)
	idle()
}

func listen(server net.Listener) {
	for {
		conn, err := server.Accept()
		go handle(conn, err)
	}
}

func handle(conn net.Conn, err error) error {
	if err != nil {
		return err
	}

	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return err
	}

	message = strings.TrimRight(message, "\n")

	if strings.HasPrefix(message, "+") {
		user := strings.TrimLeft(message, "+")
		if user == "" {
			return conn.Close()
		}

		mutex.Lock()
		scores[user] = scores[user] + 1
		mutex.Unlock()
	} else {
		err := json.NewEncoder(conn).Encode(scores)
		if err != nil {
			return err
		}
	}

	return conn.Close()
}

func idle() {
	for {
		time.Sleep(1 * time.Second)
	}
}
