package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	// Create TCP listener on port 6379
	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 6379...")

	for {
		// Accept a new connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			continue
		}
		// Handle the connection in a new goroutine
		go handleRequest(conn)
	}
}

// handleRequest manages the connection, reading data and sending a response.
func handleRequest(conn net.Conn) {
	defer conn.Close()

	// Keep reading data until the client disconnects
	for {
		// Buffer to read data
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			return // Exit the goroutine if an error occurs
		}

		// Send hardcoded PONG response for PING
		conn.Write([]byte("+PONG\r\n"))
	}
}
