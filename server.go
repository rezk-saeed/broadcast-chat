package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

type Client struct {
	ID   string
	Conn net.Conn
	Chan chan string
}

var (
	clients = make(map[string]Client)
	mutex   = sync.Mutex{}
)

func broadcast(senderID, message string) {
	mutex.Lock()
	defer mutex.Unlock()

	for id, client := range clients {
		if id != senderID { // no self echo
			client.Chan <- message
		}
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	id, _ := reader.ReadString('\n')
	id = id[:len(id)-1]

	client := Client{
		ID:   id,
		Conn: conn,
		Chan: make(chan string),
	}

	mutex.Lock()
	clients[id] = client
	mutex.Unlock()

	// Notify others
	broadcast(id, "User "+id+" joined")

	// Goroutine for sending messages
	go func() {
		for msg := range client.Chan {
			fmt.Fprintln(conn, msg)
		}
	}()

	// Receive messages
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		broadcast(id, id+": "+msg)
	}

	mutex.Lock()
	delete(clients, id)
	mutex.Unlock()
}

func main() {
	listener, _ := net.Listen("tcp", ":8080")
	fmt.Println("Server running on port 8080...")

	for {
		conn, _ := listener.Accept()
		go handleClient(conn)
	}
}
