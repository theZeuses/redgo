package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type Reader struct {
	reader *bufio.Reader
}

func NewReader(ioReader io.Reader) *Reader {
	return &Reader{
		reader: bufio.NewReader(ioReader),
	}
}

func EncodeCommandAsRespString(args []string) []byte {
	var resp string
	resp += fmt.Sprintf("*%d\r\n", len(args))

	for _, arg := range args {
		resp += fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg)
	}

	return []byte(resp)
}

func (p *Reader) readLine() (line []byte, n int, err error) {
	for {
		b, err := p.reader.ReadByte()

		if err != nil {
			return nil, 0, err
		}

		if b == '\r' {
			p.reader.ReadByte()
			break
		} else {
			n += 1
			line = append(line, b)
		}
	}

	return line, n, nil
}

func (p *Reader) readInteger() (val int64, n int, err error) {
	line, _, err := p.readLine()

	if err != nil {
		return 0, 0, err
	}

	val, err = strconv.ParseInt(string(line), 10, 64)

	if err != nil {
		return 0, 0, err
	}

	return val, n, nil
}

func (p *Reader) parseBulkStr() (Value, error) {
	val := BulkStringValue{}

	len, _, err := p.readInteger()
	if err != nil {
		return val, err
	}

	if len == -1 {
		return NullValue{}, nil
	}

	bulk := make([]byte, len)
	_, errDuringRead := p.reader.Read(bulk)

	if errDuringRead != nil {
		return val, errDuringRead
	}

	p.readLine()

	val.Val = string(bulk)

	return val, nil
}

func (p *Reader) parseArray() (Value, error) {
	val := ArrayValue{}

	len, _, err := p.readInteger()
	if err != nil {
		return val, err
	}

	val.Val = make([]Value, len)

	for i := 0; i < int(len); i++ {
		v, err := p.ParseFromRespString()

		if err != nil {
			return val, err
		}

		val.Val[i] = v
	}

	return val, nil
}

func (p *Reader) parseError() (Value, error) {
	val := ErrorValue{}

	line, _, err := p.readLine()

	if err != nil {
		return val, err
	}

	val.Val = string(line)

	return val, nil
}

func (p *Reader) parseInteger() (Value, error) {
	val := IntegerValue{}

	line, _, err := p.readLine()

	if err != nil {
		return val, err
	}

	paredInt, err := strconv.ParseInt(string(line), 10, 64)

	if err != nil {
		return val, err
	}

	val.Val, err = int(paredInt), nil

	return val, nil
}

func (p *Reader) parseNull() (Value, error) {
	val := NullValue{}

	p.readLine()

	return val, nil
}

func (p *Reader) parseSting() (Value, error) {
	val := StringValue{}

	line, _, err := p.readLine()

	if err != nil {
		return val, err
	}

	val.Val = string(line)

	return val, nil
}

func (p *Reader) ParseFromRespString() (Value, error) {
	_type, err := p.reader.ReadByte()

	if err != nil {
		return EmptyValue{}, err
	}

	switch _type {
	case ARRAY:
		return p.parseArray()
	case BULK_STRING:
		return p.parseBulkStr()
	case STRING:
		return p.parseSting()
	case INTEGER:
		return p.parseInteger()
	case ERROR:
		return p.parseError()
	case NULL:
		return p.parseNull()
	default:
		return BulkStringValue{Val: "Unrecognized Response"}, nil
	}
}

func (p StringValue) WriteToConsole() {
	fmt.Println(p.Val)
}

func (p ErrorValue) WriteToConsole() {
	fmt.Println("(error) ", p.Val)
}

func (p IntegerValue) WriteToConsole() {
	fmt.Println("(integer) ", p.Val)
}

func (p BulkStringValue) WriteToConsole() {
	fmt.Printf("\"%s\"\n", p.Val)
}

func iterateArray(v ArrayValue) string {
	val := ""

	for i, v := range v.Val {

		switch v.Type() {
		case R_STRING:
			val = v.(StringValue).Val
		case R_ERROR:
			val = "(integer) " + v.(ErrorValue).Val
		case R_INTEGER:
			val = "(integer) " + strconv.Itoa(v.(IntegerValue).Val)
		case R_BULK_STRING:
			val = "\"" + v.(BulkStringValue).Val + "\""
		case R_NULL:
			val = "(nil)"
		case R_EMPTY:
			val = ""
		}

		val = val + fmt.Sprintf("%d) %s\n", i, val)
	}

	return val
}

func (p ArrayValue) WriteToConsole() {
	for i, v := range p.Val {
		val := ""

		switch v.Type() {
		case R_STRING:
			val = v.(StringValue).Val
		case R_ERROR:
			val = "(integer) " + v.(ErrorValue).Val
		case R_INTEGER:
			val = "(integer) " + strconv.Itoa(v.(IntegerValue).Val)
		case R_BULK_STRING:
			val = "\"" + v.(BulkStringValue).Val + "\""
		case R_ARRAY:
			val = iterateArray(v.(ArrayValue))
		case R_NULL:
			val = "(nil)"
		case R_EMPTY:
			val = ""
		}

		fmt.Printf("%d) %s\n", i+1, val)
	}
}

func (p NullValue) WriteToConsole() {
	fmt.Println("(nil)")
}

func (p EmptyValue) WriteToConsole() {
	fmt.Println("")
}
