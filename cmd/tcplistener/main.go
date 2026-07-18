package main

import (
	"fmt"
	"net"

	"github.com/Xeninon/httpfromtcp/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("ConnectionError: %v\n", err)
			continue
		}

		fmt.Println("Connection Accepted")

		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("RequestError: %v\n", err)
			continue
		}

		fmt.Println("Request line:")
		fmt.Println("- Method:", req.RequestLine.Method)
		fmt.Println("- Target:", req.RequestLine.RequestTarget)
		fmt.Println("- Version:", req.RequestLine.HttpVersion)

		fmt.Println("Headers:")
		for key, value := range req.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}

		fmt.Println("Body:")
		fmt.Println(string(req.Body))
	}
}
