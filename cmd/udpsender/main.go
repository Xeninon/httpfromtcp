package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		input, err := reader.ReadString(byte('\n'))
		if err != nil {
			fmt.Println("Read Error:", err)
			continue
		}

		if _, err = conn.Write([]byte(input)); err != nil {
			fmt.Println("Write Error:", err)
			continue
		}
	}
}
