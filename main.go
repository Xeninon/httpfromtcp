package main

import (
	"io"
	"os"
	"fmt"
	"strings"
)

func main()  {
	file, err := os.Open("./messages.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	ch := getLinesChannel(file)
	for line := range ch {
		fmt.Printf("read: %s\n", line)
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
		f.Close()
	}()
	return ch
}
