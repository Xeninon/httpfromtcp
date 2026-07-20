package headers

import (
	"errors"
	"strings"
	"unicode"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	header, _, found := strings.Cut(string(data), crlf)
	if !found {
		return 0, false, nil
	}

	if header == "" {
		return len([]byte(crlf)), true, nil
	}

	key, value, found := strings.Cut(header, ":")
	if !found {
		return 0, false, errors.New("No key:value pair")
	}

	if key != strings.TrimSpace(key) {
		return 0, false, errors.New("Whitespace around field name")
	}

	for _, char := range key {
		if !unicode.IsDigit(char) && !unicode.IsLetter(char) && !strings.Contains("!#$%&'*+-.^_`|~", string(char)) {
			return 0, false, errors.New("invalid character in field name")
		}
	}

	key = strings.ToLower(key)
	value = strings.TrimSpace(value)
	if previousValue, exists := h[key]; exists {
		h[key] = previousValue + ", " + value
	} else {
		h[key] = value
	}

	return len([]byte(header + crlf)), false, nil
}

func (h Headers) Get(key string) (string, bool) {
	value, ok := h[strings.ToLower(key)]
	return value, ok
}

func (h Headers) Set(key, value string){
	h[strings.ToLower(key)] = value 
}
