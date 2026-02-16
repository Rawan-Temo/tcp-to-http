package headers

import (
	"bytes"
	"errors"
	"fmt"
)

type Headers map[string]string

//TODO
// func (h Headers) Get(name string) string {
// 	return h[name]

// }
// func (h Headers) Set(name string, value string) {
// 	h[name] = value
// }

var rn = []byte("\r\n")

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", errors.New(" MALFORMED HEADER")
	}
	name := parts[0]
	value := bytes.TrimSpace(parts[1])
	whtiespace := string(fieldLine)[len(string(name))-1]
	if string(whtiespace) == " " {
		fmt.Print("whtiespace")
		return "", "", errors.New(" MALFORMED FIELD NAME")
	}

	return string(name), string(value), nil
}

func NewHeaders() Headers {
	return Headers{}
}
func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	read := 0

	for {
		header := data[read:]

		idx := bytes.Index(header, rn)

		if idx == -1 {
			break
		}
		if idx == 0 {
			done = true
			read += len(rn)
			break
		}
		fmt.Println("idx : ", idx)

		fmt.Println("1 : ", string(data[read:read+idx]))

		name, value, err := parseHeader(data[read : read+idx])
		if err != nil {
			return 0, false, err
		}
		read += idx + len(rn)
		h[name] = value

	}

	fmt.Println("done : ", done)
	return read, done, nil
}
