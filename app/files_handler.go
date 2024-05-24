package main

import (
	"fmt"
	"net"
	"os"
  "net/http"
  "strings"
  "io/ioutil"
)


func filePost(conn net.Conn, request *http.Request, directory string, filename string) {
  body, err := ioutil.ReadAll(request.Body)
  if err != nil {
    error_string := GetErrorString("Error reading request body")
    conn.Write([]byte(error_string))
  }

  defer request.Body.Close()

  err = ioutil.WriteFile(directory + filename, []byte(body), 0644)
  if err != nil {
    error_string := GetErrorString("Error writing body to file")
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

func FilesHandler(request *http.Request, conn net.Conn) {
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

