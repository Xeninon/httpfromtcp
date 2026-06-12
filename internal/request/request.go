package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)


type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\r\n")
	requestLine, err := parseRequestLine(lines[0])
	if err != nil {
		return nil, fmt.Errorf("Error: %v", err)
	}

	return &Request{RequestLine: requestLine}, nil
}

func parseRequestLine(requestString string) (RequestLine, error) {
	parts := strings.Split(requestString, " ")
	if len(parts) != 3 {
        	return RequestLine{}, errors.New("Invalid number of parts in request line")
	}

	method := parts[0]
	target := parts[1]
	version, found := strings.CutPrefix(parts[2], "HTTP/")

	for _, r := range method {
        	if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return RequestLine{}, errors.New("Invalid Method")
        	}
	}

	if version != "1.1" || !found {
		return RequestLine{}, errors.New("Invalid Version")
	}

	return RequestLine{HttpVersion: version, RequestTarget: target, Method: method}, nil
}
