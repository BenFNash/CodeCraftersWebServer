package main

import (
	"fmt"
	"net"
  "net/http"
  "strings"
  "compress/gzip"
  "bytes"
)

func gzipResponse(conn net.Conn, body string) {
  var buf bytes.Buffer

  gzipWriter := gzip.NewWriter(&buf)

  _, err := gzipWriter.Write([]byte(body))
  if err != nil {
    error_string := GetErrorString("Error compressing string")
    conn.Write([]byte(error_string))
  }

  err = gzipWriter.Close()
  if err != nil {
    error_string := GetErrorString("Error compressing string")
    conn.Write([]byte(error_string))
  }

  compressedBody := buf.Bytes()
  response_str := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Encoding: gzip\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(compressedBody), compressedBody)
  conn.Write([]byte(response_str))
}

func EchoHandler(request *http.Request, conn net.Conn) {
  body := request.URL.Path[6:]
  encodingHeader := request.Header.Get("Accept-Encoding")

  check_contains := strings.Contains(encodingHeader, ", gzip, ")
  check_prefix := strings.HasPrefix(encodingHeader, "gzip, ")
  check_suffix := strings.HasSuffix(encodingHeader, ", gzip")
  check_equal := (encodingHeader == "gzip")

  if check_equal || check_prefix || check_suffix || check_contains {
    gzipResponse(conn, body)
    return
  }

  response_str := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(body), body)
  conn.Write([]byte(response_str))
}

