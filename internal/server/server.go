package server

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"sync/atomic"

	"github.com/Xeninon/httpfromtcp/internal/request"
	"github.com/Xeninon/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	isClosed atomic.Bool
	handler  Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	server := &Server{listener: listener, isClosed: atomic.Bool{}, handler: handler}
	go server.listen()
	return server, nil
}

func (s *Server) Close() error {
	if err := s.listener.Close(); err != nil {
		return err
	}
	s.isClosed.Store(true)
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if s.isClosed.Load() {
			return
		}
		if err != nil {
			fmt.Printf("ConnectionError: %v\n", err)
			continue
		}

		fmt.Println("Connection Accepted")
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := &HandlerError{
			StatusCode: response.StatusCodeBadRequest,
			Message:    err.Error(),
		}
		hErr.Write(conn)
		return
	}

	var buffer bytes.Buffer
	hErr := s.handler(&buffer, req)
	if hErr != nil {
		hErr.Write(conn)
		return
	}

	headers := response.GetDefaultHeaders(buffer.Len())
	response.WriteStatusLine(conn, response.StatusCode(response.StatusCodeOK))
	response.WriteHeaders(conn, headers)
	conn.Write(buffer.Bytes())
}

func (h *HandlerError) Write(conn net.Conn) {
	err := response.WriteStatusLine(conn, h.StatusCode)
	if err != nil {
		fmt.Printf("HandlingError: %v\n", err)
	}

	headers := response.GetDefaultHeaders(len(h.Message))
	err = response.WriteHeaders(conn, headers)
	if err != nil {
		fmt.Printf("HandlingError: %v\n", err)
	}

	_, err = conn.Write([]byte(h.Message))
	if err != nil {
		fmt.Printf("HandlingError: %v\n", err)
	}
}
