package proto

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	RespStatus = '+'
	RespString = '$'
	RespInt    = ':'
	RespError  = '-'
	RespArray  = '*'
)

type Reader struct {
	rd *bufio.Reader
}

func NewReader(rd io.Reader) *Reader {
	return &Reader{
		rd: bufio.NewReader(rd),
	}
}

func (r *Reader) readLine() ([]byte, error) {
	b, err := r.rd.ReadSlice('\n')

	if err != nil {
		if err != bufio.ErrBufferFull {
			return nil, err
		}

		full := make([]byte, len(b))
		copy(full, b)

		b, err = r.rd.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		full = append(full, b...) //nolint:makezero
		b = full
	}

	if len(b) <= 2 || b[len(b)-1] != '\n' || b[len(b)-2] != '\r' {
		return nil, fmt.Errorf("redis: invalid reply: %q", b)
	}

	return b[:len(b)-2], nil
}

func (r *Reader) readSlice(line []byte) ([]interface{}, error) {
	n, err := replyLen(line)
	if err != nil {
		return nil, err
	}

	val := make([]interface{}, n)
	for i := 0; i < len(val); i++ {
		v, err := r.ReadReply()
		if err != nil {
			if err == Nil {
				val[i] = nil
				continue
			}
			if err, ok := err.(RedisError); ok {
				val[i] = err
				continue
			}
			return nil, err
		}
		val[i] = v
	}
	return val, nil
}

func replyLen(line []byte) (n int, err error) {
	n, err = strconv.Atoi(string(line[1:]))
	if err != nil {
		return 0, err
	}

	if n < -1 {
		return 0, fmt.Errorf("redis: invalid reply: %q", line)
	}

	switch line[0] {
	case RespString, RespArray:
		if n == -1 {
			return 0, Nil
		}
	}
	return n, nil
}

func (r *Reader) readStringReply(line []byte) (string, error) {
	n, err := replyLen(line)
	if err != nil {
		return "", err
	}

	b := make([]byte, n+2)
	_, err = io.ReadFull(r.rd, b)
	if err != nil {
		return "", err
	}

	return string(b[:n]), nil
}

func (r *Reader) ReadLine() ([]byte, error) {
	line, err := r.readLine()

	if err != nil {
		return nil, err
	}
	switch line[0] {
	case RespError:
		return nil, ParseErrorReply(line)
	}

	return line, nil
}

func (r *Reader) ReadReply() (interface{}, error) {
	line, err := r.ReadLine()

	if err != nil {
		return nil, err
	}

	switch line[0] {
	case RespStatus:
		return string(line[1:]), nil
	case RespInt:
		return strconv.ParseInt(string(line[1:]), 10, 64)
	case RespString:
		return r.readStringReply(line)
	case RespArray:
		return r.readSlice(line)
	}
	return nil, fmt.Errorf("redis: can't parse %.100q", line)
}

// ---------------------------------------------------- //

type RedisError string

func (e RedisError) Error() string { return string(e) }

func (RedisError) RedisError() {}

func ParseErrorReply(line []byte) error {
	return RedisError(line[1:])
}

const Nil = RedisError("redis: nil")
