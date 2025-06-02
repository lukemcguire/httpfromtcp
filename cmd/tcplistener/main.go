package main

import (
	"fmt"
	"log"
	"net"

	"github.com/lukemcguire/httpfromtcp/internal/request"
)

const PORT = ":42069"

func main() {
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("could not open tcp connection on port %s: %s\n", PORT, err.Error())
	}
	defer listener.Close()

	fmt.Println("Listening for TCP traffic on", PORT)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Accepted connection from", conn.RemoteAddr())

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error getting request: %s\n", err.Error())
		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for key, value := range req.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}

		fmt.Println("Connection to", conn.RemoteAddr(), "closed")
	}
}

/*
func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		defer f.Close()
		defer close(lines)
		out := make([]byte, 8)
		line := ""
		for i := 0; ; i += 8 {
			n, err := f.Read(out)
			if err != nil {
				if errors.Is(err, io.EOF) {
					if line != "" {
						lines <- line
					}
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				return
			}
			parts := strings.Split(string(out[:n]), "\n")
			for i, part := range parts {
				if i > 0 {
					lines <- line
					line = ""
				}
				line += part
			}
		}
	}()

	return lines
}
*/
