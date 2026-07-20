package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Xeninon/httpfromtcp/internal/request"
	"github.com/Xeninon/httpfromtcp/internal/response"
	"github.com/Xeninon/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, r *request.Request) {
	headers := response.GetDefaultHeaders(0)
	if r.RequestLine.RequestTarget == "/yourproblem" {
		body := "<html><head><title>400 Bad Request</title></head><body><h1>Bad Request</h1><p>Your request honestly kinda sucked.</p></body></html>"
		w.WriteStatusLine(response.StatusCodeBadRequest)
		headers.Set("content-length", fmt.Sprintf("%v", len([]byte(body))))
		headers.Set("content-type", "text/html")
		w.WriteHeaders(headers)
		w.WriteBody([]byte(body))
		return
	}

	if r.RequestLine.RequestTarget == "/myproblem" {
		body := "<html><head><title>500 Internal Server Error</title></head><body><h1>Internal Server Error</h1><p>Okay, you know what? This one is on me.</p></body></html>"
		w.WriteStatusLine(response.StatusCodeInternalServerError)
		headers.Set("content-length", fmt.Sprintf("%v", len([]byte(body))))
		headers.Set("content-type", "text/html")
		w.WriteHeaders(headers)
		w.WriteBody([]byte(body))
		return
	}

	body := "<html><head><title>200 OK</title></head><body><h1>Success!</h1><p>Your request was an absolute banger.</p></body></html>"
	w.WriteStatusLine(response.StatusCodeOK)
	headers.Set("content-length", fmt.Sprintf("%v", len([]byte(body))))
	headers.Set("content-type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody([]byte(body))
}
