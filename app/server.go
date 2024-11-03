package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	// Create TCP listener on port 6379
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	// Accept incoming connections
	conn, err := l.Accept()
	for {
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go handleRequest(conn)
	}
}
func handleRequest(conn net.Conn) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	conn.Write([]byte("+PONG\r\n"))
}
