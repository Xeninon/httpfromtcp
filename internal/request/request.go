package request

import (
	"errors"
	"io"
	"strings"
	"unicode"
)

const bufferSize = 8

type State int

const (
	Initialized State = iota
	Done
)

type Request struct {
	RequestLine RequestLine
	State       State
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := &Request{State: Initialized}
	buf := make([]byte, bufferSize)
	readToIndex := 0
	for request.State != Done {
		if readToIndex == cap(buf) {
			tmp := make([]byte, 2*cap(buf))
			copy(tmp, buf)
			buf = tmp
		}

		bytesRead, err := reader.Read(buf[readToIndex:])
		if err == io.EOF {
			request.State = Done
			break
		}

		readToIndex += bytesRead

		bytesParsed, err := request.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		readToIndex -= bytesParsed
		tmp := make([]byte, max(readToIndex, bufferSize))
		copy(tmp, buf[bytesParsed:])
		buf = tmp
	}
	return request, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.State {
	case Initialized:
		bytesParsed, requestLine, err := parseRequestLine(data)
		if bytesParsed != 0 {
			r.RequestLine = requestLine
			r.State = Done
		}

		return bytesParsed, err
	case Done:
		return 0, errors.New("parsed done request")
	default:
		return 0, errors.New("unknown request state")
	}
}

func parseRequestLine(data []byte) (int, RequestLine, error) {
	request, _, found := strings.Cut(string(data), "\r\n")
	if !found {
		return 0, RequestLine{}, nil
	}

	parts := strings.Split(request, " ")
	if len(parts) != 3 {
		return 0, RequestLine{}, errors.New("Invalid number of parts in request line")
	}

	method := parts[0]
	target := parts[1]
	version, found := strings.CutPrefix(parts[2], "HTTP/")

	for _, r := range method {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return 0, RequestLine{}, errors.New("Invalid Method")
		}
	}

	if version != "1.1" || !found {
		return 0, RequestLine{}, errors.New("Invalid Version")
	}

	return len([]byte(request + "\r\n")), RequestLine{HttpVersion: version, RequestTarget: target, Method: method}, nil
}
