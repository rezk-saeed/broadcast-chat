package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, _ := net.Dial("tcp", "localhost:8080")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your ID: ")
	id, _ := reader.ReadString('\n')
	fmt.Fprint(conn, id)

	// Receive messages from server
	go func() {
		serverReader := bufio.NewReader(conn)
		for {
			msg, _ := serverReader.ReadString('\n')
			fmt.Print(msg)
		}
	}()

	// Send messages
	for {
		text, _ := reader.ReadString('\n')
		fmt.Fprint(conn, text)
	}
}
