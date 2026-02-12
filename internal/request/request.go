package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	State       parserState
}
type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var ERROR_START_LINE = fmt.Errorf(" ERROR")

var SEPERATOR = "\r\n"

type parserState string

func (r *Request) done() bool {
	return r.State == StateDone
}

const (
	StateInit parserState = "init"
	StateDone parserState = "done"
)

func newRequest() *Request {
	return &Request{
		RequestLine: RequestLine{},
		State:       StateInit,
	}
}

func (r *RequestLine) ValidateHttp() bool {
	return r.HttpVersion == "1.1"
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	request := newRequest()
	buff := make([]byte, 1024)
	buffLen := 0
	for !request.done() {
		n, err := reader.Read(buff[buffLen:])
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		buffLen += n
		readn, err := request.parse(buff[:buffLen])
		if err != nil {
			return nil, err
		}
		copy(buff, buff[readn:buffLen])
		buffLen -= readn

	}

	return request, nil

}

func ParseRequestLine(req []byte) (*RequestLine, int, error) {
	idx := strings.Index(string(req), SEPERATOR)
	if idx == -1 {
		return nil, 0, nil
	}

	requestLine := req[:idx]
	read := idx + len(SEPERATOR)
	arr := strings.Split(string(requestLine), " ")

	if len(arr) != 3 {
		return nil, 0, ERROR_START_LINE
	}
	http := strings.Split(arr[2], "/")
	if len(http) != 2 || http[0] != "HTTP" || http[1] != "1.1" {
		return nil, 0, ERROR_START_LINE
	}

	rl := &RequestLine{
		Method:        arr[0],
		RequestTarget: arr[1],
		HttpVersion:   http[1],
	}
	if !rl.ValidateHttp() {
		return nil, 0, ERROR_START_LINE

	}
	return rl, read, nil
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {

		switch r.State {
		case StateInit:
			rl, n, err := ParseRequestLine(data[read:])
			if err != nil {
				return read, err
			}
			if n == 0 {
				break outer
			}
			read += n
			r.RequestLine = *rl
			r.State = StateDone

		case StateDone:
			break outer
		}
	}
	return read, nil
}
