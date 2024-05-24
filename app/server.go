package main

import (
	"fmt"
	"net"
	"os"
  "net/http"
  "bufio"
  "strings"
)




func handler(conn net.Conn) {
  defer conn.Close()
  request, err := http.ReadRequest(bufio.NewReader(conn))

  if err != nil {
    fmt.Println("Error reading request. ", err.Error())
    return
  }

  root := ""
  echo := "echo"
  user_agent := "user-agent"
  files := "files"

  path := strings.Split(request.URL.Path, "/")[1:]
  switch path[0] {
    case root: 
      conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
    case echo:
      EchoHandler(request, conn)
    case user_agent:
      UserAgentHandler(request, conn)
  case files:
      FilesHandler(request, conn)
    default:
      conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
  }
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

  for {
    conn, err := ln.Accept()
	  if err != nil {
		  fmt.Println("Error accepting connection: ", err.Error())
		  os.Exit(1)
	  }

    go handler(conn)
  }
}
