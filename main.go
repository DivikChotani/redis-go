package main

import (
	"fmt"
	"io"
	"net"
	"os"
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

	//for loop for continous connections
	for {
		//create a 1024 byte buffer to store the respones for any request
		buf := make([]byte, 1024)

		// read message from client
		//any request that has come in is read here

		_, err = conn.Read(buf)
		//logging to confirm request, printing out byte version of request
		fmt.Println("REQUEST:", buf)

		//error handling, except for an EOF error
		if err != nil {
			if err == io.EOF {
				//logging to confirm eof
				fmt.Println("EOF: ", err.Error())
				break
			}
			//any other error is bad and should be logged
			fmt.Println("error reading from client: ", err.Error())
			os.Exit(1)
		}

		// ignore request and send back OK
		conn.Write([]byte("+OK\r\n"))
	}

}
