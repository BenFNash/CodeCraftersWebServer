package main

import (
	"fmt"
	"net"
	"os"
  "net/http"
  "bufio"
)

func handler(conn net.Conn) {
  defer conn.Close()
  request, err := http.ReadRequest(bufio.NewReader(conn))

  if err != nil {
    fmt.Println("Error reading request. ", err.Error())
    return
  }

  if request.URL.Path == "/" {
    conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
    return
  }

  conn.Write([]byte("HTTP/1.1 400 Not Found\r\n\r\n"))
}


func main() {
	fmt.Println("Logs from your program will appear here!")

	ln, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

  fmt.Println("Listening on port 4221")
  defer ln.Close()

  conn, err := ln.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

  handler(conn)
}
