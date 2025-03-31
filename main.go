package main

import (
	"fmt"
	//"io"
	"net"
	//"os"
)

func main() {

	//stdout letting user know where server is running
	fmt.Println("Listening on port 6379")

	//create a new server

	//create a listener object from net.listen that listens on port 6379
	//specify a tcp connection
	l, err := net.Listen("tcp", ":6379")

	//error handling, print error and leave
	if err != nil {
		fmt.Println(err)
		return
	}

	//create a connection from the listener
	conn, err := l.Accept()
	//error checking again
	if err != nil {
		fmt.Println(err)
		return
	}
	//ensure closing connection at exit
	defer conn.Close()

}
