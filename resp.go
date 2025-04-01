package main

import (
	"bufio"
	"io"
	"strconv"
)

// basically when an input comes in, we can use this to exctract the type of input from the bytes
const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

// this will hold the type of input provided, and based off of that
// it will store the content in the correct type in the struct
type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

// create a response type here that is a bufio reader
type Resp struct {
	reader *bufio.Reader
}

// basically this will create a Resp type that you can then use for future functions
// this will convert the byte rd into a reader and place that as a value in Resp
func newResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

// reads the line from the resp and then returns it with the last to \r and \n chars removed
func (r *Resp) readLine() (line []byte, n int, err error) {

	//infite for loop
	for {
		//read one byte at a time
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		//increment number of bytes read
		n += 1

		//append bite to the line
		line = append(line, b)

		//break when certain characters are added to line
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}

	//return line with the last two characters removed
	return line[:len(line)-2], n, nil
}

// read the integer from the buffer indicating how long the input is
func (r *Resp) readInteger() (x int, n int, err error) {

	//call our read line function to get the line
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}

	//parse the int from our line after converting it into a string
	//the other two inputs are the base and the number of bits used to store the int
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	//return that int
	return int(i64), n, nil

}
