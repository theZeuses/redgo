package main

import (
	"bufio"
	"io"
	"net"
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
	default:
		return EmptyValue{}, nil
	}
}

func (p *Reader) parseBulkStr() (Value, error) {
	val := BulkStringValue{}

	len, _, err := p.readInteger()
	if err != nil {
		return val, err
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

type Writer struct {
	writer io.Writer
}

func NewWriter(conn net.Conn) *Writer {
	return &Writer{
		writer: io.Writer(conn),
	}
}

func (v StringValue) Marshal() (bytes []byte) {
	bytes = append(bytes, STRING)
	bytes = append(bytes, []byte(v.Val)...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v ErrorValue) Marshal() (bytes []byte) {
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.Val...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v IntegerValue) Marshal() (bytes []byte) {
	bytes = append(bytes, INTEGER)
	bytes = append(bytes, strconv.FormatInt(int64(v.Val), 10)...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v BulkStringValue) Marshal() (bytes []byte) {
	bytes = append(bytes, BULK_STRING)
	bytes = append(bytes, strconv.Itoa(len(v.Val))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.Val...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v ArrayValue) Marshal() (bytes []byte) {
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len(v.Val))...)
	bytes = append(bytes, '\r', '\n')

	for _, value := range v.Val {
		bytes = append(bytes, value.Marshal()...)
	}

	return bytes
}

func (v NullValue) Marshal() (bytes []byte) {
	bytes = append(bytes, NULL)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v EmptyValue) Marshal() (bytes []byte) {
	return []byte{}
}

func (w *Writer) WriteAsRespString(value Value) error {
	bytes := value.Marshal()
	println("Writing bytes:", string(bytes))
	_, err := w.writer.Write(bytes)

	if err != nil {
		return err
	}

	return nil
}
