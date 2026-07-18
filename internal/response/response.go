package response

import (
	"fmt"
	"io"

	"github.com/Xeninon/httpfromtcp/internal/headers"
)

type StatusCode int

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	status := ""
	switch statusCode {
	case 200:
		status = "HTTP/1.1 200 OK\r\n"
	case 400:
		status = "HTTP/1.1 400 Bad Request\r\n"
	case 500:
		status = "HTTP/1.1 500 Internal Server Error\r\n"
	default:
		status = fmt.Sprintf("HTTP/1.1 %v \r\n", statusCode)
	}
	_, err := w.Write([]byte(status))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers["content-length"] = fmt.Sprintf("%v", contentLen)
	headers["connection"] = "close"
	headers["content-type"] = "text/plain"
	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for field, value := range headers {
		_, err := fmt.Fprintf(w, "%v: %v\r\n", field, value)
		if err != nil {
			return err
		}
	}

	_, err := w.Write([]byte("\r\n"))
	return err
}
