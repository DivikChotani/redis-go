package main

import (
	"bufio"
	"fmt"
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

// create an overarching read function that accepts array or bulk string inputs
func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()

	if err != nil {
		return Value{}, err
	}
	//based on the first byte, we can determine the type of input
	switch _type {
	case ARRAY:
		fmt.Printf("ARRAY TYPE")
		return r.readArray()
	case BULK:
		fmt.Printf("Bulk string type")
		return r.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(_type))
		return Value{}, nil
	}
}

// reading the array type
func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.typ = "array"
	//create a type array here

	// read length of array
	length, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	v.array = make([]Value, length)

	//for every element in the array, call read to read it and add to value
	for i := 0; i < length; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}
		//add after the recursive call here
		v.array[i] = val
	}
	return v, nil
}

// no recursion needed for bulk, since it will be non nested
func (r *Resp) readBulk() (Value, error) {
	v := Value{}

	//bulk
	v.typ = "bulk"

	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	//create space in memory store the string
	bulk := make([]byte, len)

	//fill the buffer
	r.reader.Read(bulk)

	//add to calue
	v.bulk = string(bulk)

	// Read the trailing CRLF
	r.readLine()

	return v, nil
}
