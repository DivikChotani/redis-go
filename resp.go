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
func (r *Resp) readIntegerLength() (x int, n int, err error) {
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
	case STRING:
		return r.readString()
	case ERROR:
		return r.readError()
	case INTEGER:
		return r.readIntegerValue()
	default:
		// If we get a carriage return, it's likely part of a CRLF sequence
		// Try to read the next byte to see if it's a newline
		if _type == '\r' {
			next, err := r.reader.ReadByte()
			if err != nil {
				return Value{}, err
			}
			if next == '\n' {
				// This was a CRLF sequence, try reading the next type
				return r.Read()
			}
			// Put the bytes back if it wasn't a CRLF
			if err := r.reader.UnreadByte(); err != nil {
				return Value{}, err
			}
			if err := r.reader.UnreadByte(); err != nil {
				return Value{}, err
			}
		}
		return Value{}, fmt.Errorf("unsupported RESP type: %q", _type)
	}
}

// reading the array type
func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.typ = "array"

	// read length of array
	length, _, err := r.readIntegerLength()
	if err != nil {
		return v, err
	}

	// Initialize the array slice
	v.array = make([]Value, length)

	// For every element in the array, call Read to read it
	for i := 0; i < length; i++ {
		// Read the type byte for each element
		_, err := r.reader.ReadByte()
		if err != nil {
			return v, err
		}

		// Put the type byte back in the reader since Read() will need it
		if err := r.reader.UnreadByte(); err != nil {
			return v, err
		}

		// Read the element
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
	length, _, err := r.readIntegerLength()
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

func (r *Resp) readString() (Value, error) {
	v := Value{}
	v.typ = "string"

	line, _, err := r.readLine()
	if err != nil {
		return v, err
	}

	v.str = string(line)
	return v, nil
}

func (r *Resp) readError() (Value, error) {
	v := Value{}
	v.typ = "error"

	line, _, err := r.readLine()
	if err != nil {
		return v, err
	}

	v.str = string(line)
	return v, nil
}

func (r *Resp) readIntegerValue() (Value, error) {
	v := Value{}
	v.typ = "integer"

	line, _, err := r.readLine()
	if err != nil {
		return v, err
	}

	num, err := strconv.Atoi(string(line))
	if err != nil {
		return v, err
	}

	v.num = num
	return v, nil
}

// Marshal serializes a Value into a byte slice following a Redis-like protocol.
// This protocol distinguishes between different types like simple strings, bulk strings,
// arrays, nulls, and errors using specific prefix characters.
func (v Value) Marshal() []byte {
	switch v.typ {
	case "array":
		// Serialize an array of values
		return v.marshalArray()
	case "bulk":
		// Serialize a binary-safe string (bulk string)
		return v.marshalBulk()
	case "string":
		// Serialize a simple string
		return v.marshalString()
	case "null":
		// Serialize a null (RESP null bulk string)
		return v.marshallNull()
	case "error":
		// Serialize an error message
		return v.marshallError()
	default:
		// Unknown type: return empty byte slice
		return []byte{}
	}
}

// marshalString converts a simple string into a RESP format string.
// RESP format: +<string>\r\n
// Example: +OK\r\n
func (v Value) marshalString() []byte {
	var bytes []byte

	// Prefix for simple strings is '+'
	bytes = append(bytes, STRING)

	// Append the actual string characters
	bytes = append(bytes, v.str...)

	// End with CRLF
	bytes = append(bytes, '\r', '\n')

	return bytes
}

// marshalBulk converts a bulk string into a RESP format bulk string.
// RESP format: $<length>\r\n<contents>\r\n
// Example: $6\r\nfoobar\r\n
func (v Value) marshalBulk() []byte {
	var bytes []byte

	// Prefix for bulk strings is '$'
	bytes = append(bytes, BULK)

	// Append the length of the content in decimal form
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)

	// CRLF to end the length line
	bytes = append(bytes, '\r', '\n')

	// Append the actual bulk content (binary-safe)
	bytes = append(bytes, v.bulk...)

	// Final CRLF
	bytes = append(bytes, '\r', '\n')

	return bytes
}

// marshalArray converts an array of Value objects into a RESP array.
// RESP format: *<num elements>\r\n<element1><element2>...\r\n
// Each element is recursively marshaled with the appropriate type rules.
// Example: *2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n
func (v Value) marshalArray() []byte {
	length := len(v.array) // Total number of elements
	var bytes []byte

	// Prefix for arrays is '*'
	bytes = append(bytes, ARRAY)

	// Append the number of elements in the array
	bytes = append(bytes, strconv.Itoa(length)...)

	// End header with CRLF
	bytes = append(bytes, '\r', '\n')

	// Marshal each element in the array recursively
	for i := 0; i < length; i++ {
		bytes = append(bytes, v.array[i].Marshal()...)
	}

	return bytes
}

// marshallError converts an error string into RESP error format.
// RESP format: -<error message>\r\n
// Example: -ERR unknown command\r\n
func (v Value) marshallError() []byte {
	var bytes []byte

	// Prefix for errors is '-'
	bytes = append(bytes, ERROR)

	// Append the error message string
	bytes = append(bytes, v.str...)

	// End with CRLF
	bytes = append(bytes, '\r', '\n')

	return bytes
}

// marshallNull returns a RESP representation of a null bulk string.
// RESP format for null: $-1\r\n
// This is used to represent a null value or missing data.
func (v Value) marshallNull() []byte {
	return []byte("$-1\r\n")
}
