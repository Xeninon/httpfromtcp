package main

import (
	"net"
	"io"
	"fmt"
	"strings"
)

func main()  {
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

		for line := range getLinesChannel(conn) {
			fmt.Println(line)
		}
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	current := ""
	b := make([]byte, 8)
	go func(){
		for {
			if _, err := f.Read(b); err == io.EOF {
				break
			}
			parts := strings.Split(string(b), "\n")
			current += parts[0]
			if len(parts) > 1 {
				ch <- current
				current = parts[1]
			}
		}
		close(ch)
	}()
	return ch
}
