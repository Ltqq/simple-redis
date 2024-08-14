package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	typ   string  //data type
	str   string  //string
	num   int     //int
	bulk  string  //strings
	array []Value //array
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}

func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.typ = "array"

	//读数组的长度
	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	// 对于每一行，解析并读取值
	v.array = make([]Value, 0)
	for i := 0; i < len; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}
		//将解析后的值追加到数组
		v.array = append(v.array, val)
	}

	return v, nil
}

func (r *Resp) readBulk() (Value, error) {
	v := Value{}
	v.typ = "bulk"

	//读取字符串的长度
	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	var bulk = make([]byte, len)
	r.reader.Read(bulk)
	v.bulk = string(bulk)
	// 在读取字符串后调用读取每个批量字符串后面的 '\r\n'。如果不这样，指针将留在 '\r' 处，Read 方法将无法正确读取下一个批量字符串。
	r.readLine()
	return v, nil
}

func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()

	if err != nil {
		return Value{}, err
	}

	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(_type))
		return Value{}, nil
	}
}

func convert(input string) (string, error) {
	reader := bufio.NewReader(strings.NewReader(input))

	b, err := reader.ReadByte()
	if err != nil {
		return "", err
	}
	if b != '$' {
		fmt.Println("Invalid type, expecting bulk strings only")
		return "", errors.New("invalid type, expecting bulk strings only")
	}
	size, err := reader.ReadByte()
	if err != nil {
		return "", err
	}
	strSize, _ := strconv.ParseInt(string(size), 10, 64)

	// consume /r/n
	reader.ReadByte()
	reader.ReadByte()

	name := make([]byte, strSize)
	_, err = reader.Read(name)
	if err != nil {
		return "", err
	}

	return string(name), nil
}
