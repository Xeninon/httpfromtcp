package server

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/Xeninon/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	isClosed atomic.Bool
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	server := &Server{listener: listener, isClosed: atomic.Bool{}}
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
			break
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
	err := response.WriteStatusLine(conn, response.StatusCode(200))
	if err != nil {
		fmt.Printf("HandlingError: %v\n", err)
	}

	headers := response.GetDefaultHeaders(0)
	err = response.WriteHeaders(conn, headers)
	if err != nil {
		fmt.Printf("HandlingError: %v\n", err)
	}
}
