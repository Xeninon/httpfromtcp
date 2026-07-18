package request

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/Xeninon/httpfromtcp/internal/headers"
)

const crlf = "\r\n"
const bufferSize = 8

type State int

const (
	requestStateInitialized State = iota
	requestStateParsingHeaders
	requestStateParsingBody
	requestStateDone
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	state       State
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := &Request{state: requestStateInitialized, Headers: headers.NewHeaders()}
	buf := make([]byte, bufferSize)
	readToIndex := 0
	for request.state != requestStateDone {
		if readToIndex == cap(buf) {
			tmp := make([]byte, 2*cap(buf))
			copy(tmp, buf)
			buf = tmp
		}

		bytesRead, err := reader.Read(buf[readToIndex:])
		EOF := false
		if err == io.EOF {
			EOF = true
		}
		if err != nil && err != io.EOF {
			return nil, err
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

		if EOF {
			break
		}
	}
	if request.state != requestStateDone {
		return nil, errors.New("Parsing didn't finish")
	}
	return request, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
		bytesParsed, requestLine, err := parseRequestLine(data)
		if bytesParsed != 0 {
			r.RequestLine = requestLine
			r.state = requestStateParsingHeaders
		}

		return bytesParsed, err
	case requestStateParsingHeaders:
		totalBytesParsed := 0
		for r.state != requestStateParsingBody {
			n, err := r.parseSingle(data[totalBytesParsed:])
			if err != nil {
				return 0, err
			}
			if n == 0 {
				break
			}
			totalBytesParsed += n
		}
		return totalBytesParsed, nil
	case requestStateParsingBody:
		lengthStr, exists := r.Headers.Get("Content-Length")
		if !exists {
			r.state = requestStateDone
			return 0, nil
		}

		length, err := strconv.Atoi(lengthStr)
		if err != nil || length < 0 {
			return 0, errors.New("Improper Content-Length value")
		}

		r.Body = append(r.Body, data...)
		if len(r.Body) > length {
			return 0, errors.New("Body longer than Content-Length")
		}

		if len(r.Body) == length {
			r.state = requestStateDone
		}

		return len(data), nil
	case requestStateDone:
		return 0, errors.New("Parsed done request")
	default:
		return 0, errors.New("Unknown request state")
	}
}

func (r *Request) parseSingle(data []byte) (int, error) {
	n, done, err := r.Headers.Parse(data)
	if done {
		r.state = requestStateParsingBody
	}

	return n, err
}

func parseRequestLine(data []byte) (int, RequestLine, error) {
	request, _, found := strings.Cut(string(data), crlf)
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

	return len([]byte(request + crlf)), RequestLine{HttpVersion: version, RequestTarget: target, Method: method}, nil
}
