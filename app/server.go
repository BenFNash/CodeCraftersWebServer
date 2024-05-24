package main

import (
	"fmt"
	"net"
	"os"
  "net/http"
  "bufio"
  "strings"
  "io/ioutil"
)

func filePost(conn net.Conn, request *http.Request, directory string, filename string) {
  body, err := ioutil.ReadAll(request.Body)
  if err != nil {
    error_string := fmt.Sprintf("HTTP/1.1 500 Internal Server Error\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len("Error reading request body"), "Error reading request body")
    conn.Write([]byte(error_string))
  }

  defer request.Body.Close()

  err = ioutil.WriteFile(directory + filename, []byte(body), 0644)
  if err != nil {
    error_string := fmt.Sprintf("HTTP/1.1 500 Internal Server Error\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len("Error writing body to file"), "Error writing body to file")
    conn.Write([]byte(error_string))
  }

  conn.Write([]byte("HTTP/1.1 201 Created\r\n\r\n"))
}

func fileGet(conn net.Conn, directory string, filename string) {

  fileContent, err := os.ReadFile(directory + filename)
  
  if err != nil {
    conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
    return
  }

  response_string := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(fileContent), fileContent)
  conn.Write([]byte(response_string))
}


func filesHandler(request *http.Request, conn net.Conn) {
  argsDir := os.Args[2]
  filename := strings.Split(request.URL.Path, "/")[2]

  var directory string
  if argsDir[len(argsDir)-1] != '/' {
    directory = argsDir + "/"
  } else {
    directory = argsDir
  }

  switch request.Method {
  case http.MethodGet:
    fileGet(conn, directory, filename)
  case http.MethodPost:
    filePost(conn, request, directory, filename)
  }


}


func userAgentHandler(request *http.Request, conn net.Conn) {
  userAgent := request.UserAgent()
  response_str := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgent), userAgent)
  conn.Write([]byte(response_str))
}


func echoHandler(request *http.Request, conn net.Conn) {
  body := request.URL.Path[6:]
  encodingHeader := request.Header.Get("Accept-Encoding")

  var response_str string
  if encodingHeader == "gzip" {
  response_str = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Encoding: gzip\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(body), body)
  } else {
    response_str = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(body), body)
  }

  conn.Write([]byte(response_str))
}

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
      echoHandler(request, conn)
    case user_agent:
      userAgentHandler(request, conn)
  case files:
      filesHandler(request, conn)
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
