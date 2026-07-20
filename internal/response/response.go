package response

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/Xeninon/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusCodeOK StatusCode = 200
	StatusCodeBadRequest StatusCode = 400
	StatusCodeInternalServerError StatusCode = 500
)

type writerState int

const (
	writerStateInitialized writerState = iota
	writerStateWritingHeaders
	writerStateWritingBody
	writerStateDone
)

type Writer struct {
	Writer io.Writer
	writerState writerState
}

func NewWriter(conn net.Conn) *Writer {
	return &Writer{Writer: conn, writerState: writerStateInitialized}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != writerStateInitialized {
		return errors.New("response written out of order")
	}

	status := ""
	switch statusCode {
	case StatusCodeOK:
		status = "HTTP/1.1 200 OK\r\n"
	case StatusCodeBadRequest:
		status = "HTTP/1.1 400 Bad Request\r\n"
	case StatusCodeInternalServerError:
		status = "HTTP/1.1 500 Internal Server Error\r\n"
	default:
		status = fmt.Sprintf("HTTP/1.1 %v \r\n", statusCode)
	}
	_, err := w.Writer.Write([]byte(status))
	w.writerState = writerStateWritingHeaders
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers["content-length"] = fmt.Sprintf("%v", contentLen)
	headers["connection"] = "close"
	headers["content-type"] = "text/plain"
	return headers
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState != writerStateWritingHeaders {
		return errors.New("response written out of order")
	}

	for field, value := range headers {
		_, err := fmt.Fprintf(w.Writer, "%v: %v\r\n", field, value)
		if err != nil {
			return err
		}
	}

	_, err := w.Writer.Write([]byte("\r\n"))
	w.writerState = writerStateWritingBody
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != writerStateWritingBody {
		return 0, errors.New("response written out of order")
	}

	w.writerState = writerStateDone
	return w.Writer.Write(p)
}
