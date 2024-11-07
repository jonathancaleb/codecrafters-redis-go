package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
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
	reader := bufio.NewReader(conn)

	for {
		// Parse the command from the client
		command, args, err := parseRESP(reader)
		if err != nil {
			fmt.Println("Error parsing command:", err)
			return
		}

		// Convert the command to uppercase to make it case-insensitive
		command = strings.ToUpper(command)

		// Handle PING and ECHO commands
		switch command {
		case "PING":
			conn.Write([]byte("+PONG\r\n"))
		case "ECHO":
			if len(args) > 0 {
				// Send the argument back as a bulk string
				response := formatBulkString(args[0])
				conn.Write([]byte(response))
			} else {
				// If no argument is provided, send an empty bulk string
				conn.Write([]byte("$0\r\n\r\n"))
			}
		default:
			// Handle unknown commands (for future extensibility)
			conn.Write([]byte("-ERR unknown command\r\n"))
		}
	}
}

// parseRESP parses the RESP protocol to extract the command and arguments.
func parseRESP(reader *bufio.Reader) (string, []string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", nil, err
	}
	line = strings.TrimSpace(line)

	// Check if the input is an array (starts with *)
	if len(line) == 0 || line[0] != '*' {
		return "", nil, fmt.Errorf("invalid RESP array")
	}

	// Read the number of elements in the array
	numElements, err := parseInteger(line[1:])
	if err != nil {
		return "", nil, err
	}

	// Read each element in the array
	elements := make([]string, numElements)
	for i := 0; i < numElements; i++ {
		// Read the bulk string header ($ followed by length)
		line, err = reader.ReadString('\n')
		if err != nil {
			return "", nil, err
		}
		if len(line) == 0 || line[0] != '$' {
			return "", nil, fmt.Errorf("expected bulk string")
		}

		// Read the length of the string
		strLen, err := parseInteger(line[1:])
		if err != nil {
			return "", nil, err
		}

		// Read the actual string
		buf := make([]byte, strLen+2) // +2 for \r\n
		_, err = reader.Read(buf)
		if err != nil {
			return "", nil, err
		}
		elements[i] = string(buf[:strLen])
	}

	// First element is the command, remaining are arguments
	return elements[0], elements[1:], nil
}

// parseInteger parses an integer from a string and returns it.
func parseInteger(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}

// formatBulkString formats a string as a RESP bulk string.
func formatBulkString(s string) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)
}
