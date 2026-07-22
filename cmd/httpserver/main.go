package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
	path := r.RequestLine.RequestTarget
	switch {
	case strings.HasPrefix(path, "/yourproblem"):
		handlerYourProblem(w, r)
	case strings.HasPrefix(path, "/myproblem"):
		handlerMyProblem(w, r)
	case strings.HasPrefix(path, "/httpbin/"):
		handlerHTTPBin(w, r)
	default:
		handlerDefault(w, r)
	}
}

func handlerYourProblem(w *response.Writer, r *request.Request) {
	headers := response.GetDefaultHeaders(0)
	body := "<html><head><title>400 Bad Request</title></head><body><h1>Bad Request</h1><p>Your request honestly kinda sucked.</p></body></html>"
	w.WriteStatusLine(response.StatusCodeBadRequest)
	headers.Set("content-length", fmt.Sprintf("%v", len([]byte(body))))
	headers.Set("content-type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody([]byte(body))
}

func handlerMyProblem(w *response.Writer, r *request.Request) {
	headers := response.GetDefaultHeaders(0)
	body := "<html><head><title>500 Internal Server Error</title></head><body><h1>Internal Server Error</h1><p>Okay, you know what? This one is on me.</p></body></html>"
	w.WriteStatusLine(response.StatusCodeInternalServerError)
	headers.Set("content-length", fmt.Sprintf("%v", len([]byte(body))))
	headers.Set("content-type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody([]byte(body))
}

func handlerDefault(w *response.Writer, r *request.Request) {
	headers := response.GetDefaultHeaders(0)
	body := "<html><head><title>200 OK</title></head><body><h1>Success!</h1><p>Your request was an absolute banger.</p></body></html>"
	w.WriteStatusLine(response.StatusCodeOK)
	headers.Set("content-length", fmt.Sprintf("%v", len([]byte(body))))
	headers.Set("content-type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody([]byte(body))
}

func handlerHTTPBin(w *response.Writer, r *request.Request) {
	headers := response.GetDefaultHeaders(0)
	w.WriteStatusLine(response.StatusCodeOK)
	headers.Delete("content-length")
	headers.Set("Transfer-Encoding", "chunked")
	w.WriteHeaders(headers)
	proxyPath := "https://httpbin.org/" + strings.TrimPrefix(r.RequestLine.RequestTarget, "/httpbin/")
	resp, err := http.Get(proxyPath)
	if err != nil {
		fmt.Println(err)
		w.WriteChunkedBodyDone()
	}
	defer resp.Body.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			w.WriteChunkedBody(buffer[:n])
		}
		if err != nil {
			break
		}
	}
	w.WriteChunkedBodyDone()
}
