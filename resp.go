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
func (r *Resp) readLine() ([]byte, int, error) {
	line, err := r.reader.ReadBytes('\n')
	if err != nil {
		return nil, 0, err
	}
	if len(line) < 2 || line[len(line)-2] != '\r' {
		return nil, 0, fmt.Errorf("invalid line ending: %q", line)
	}
	return line[:len(line)-2], len(line), nil
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
	// Read the first byte which indicates the type
	_type, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	fmt.Printf("RESP type byte: %q\n", _type)

	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		return Value{}, fmt.Errorf("unsupporsted RESP type: %q", _type)
	}
}

// reading the array type
func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.typ = "array"

	// read length of array
	length, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	// Initialize the array slice
	v.array = make([]Value, length)

	// For every element in the array, call Read to read it
	for i := 0; i < length; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}
		v.array[i] = val
	}

	return v, nil
}

// no recursion needed for bulk, since it will be non nested
func (r *Resp) readBulk() (Value, error) {
	v := Value{}
	v.typ = "bulk"

	// Read the bulk string length
	length, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	// Handle NULL bulk string ($-1\r\n)
	if length == -1 {
		v.bulk = ""
		return v, nil
	}

	// Read exactly `length` bytes
	bulk := make([]byte, length)
	_, err = io.ReadFull(r.reader, bulk)
	if err != nil {
		return v, err
	}
	v.bulk = string(bulk)

	// Then read the trailing \r\n (exactly two bytes)
	cr, err := r.reader.ReadByte()
	if err != nil {
		return v, err
	}
	lf, err := r.reader.ReadByte()
	if err != nil {
		return v, err
	}
	if cr != '\r' || lf != '\n' {
		return v, fmt.Errorf("expected CRLF after bulk string, got: %q%q", cr, lf)
	}

	return v, nil
}
