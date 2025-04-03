package main

import (
	"fmt"
	"net"
)

func main() {
	//stdout letting user know where server is running
	fmt.Println("Listening on port 6379")

	//create a new server
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		//create a connection from the listener
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Handle each connection in a separate goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("New connection established")

	// Create a new RESP parser for this connection
	resp := newResp(conn)

	// Try to read one command
	value, err := resp.Read()
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}

	// Print the parsed command in a more readable format
	if value.typ == "array" {
		fmt.Printf("Received command array with %d elements:\n", len(value.array))
		for i, elem := range value.array {
			fmt.Printf("  [%d] Type: %q, Value: %q\n", i, elem.typ, elem.bulk)
		}
	} else {
		switch value.typ {
		case "string":
			fmt.Printf("Received string: %q\n", value.str)
		case "bulk":
			fmt.Printf("Received bulk string: %q\n", value.bulk)
		case "integer":
			fmt.Printf("Received integer: %d\n", value.num)
		case "error":
			fmt.Printf("Received error: %q\n", value.str)
		default:
			fmt.Printf("Received value of type: %s\n", value.typ)
		}
	}

	// Send OK response
	conn.Write([]byte("+OK\r\n"))
}
