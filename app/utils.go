package main

import(
  "fmt"
)
func GetErrorString(message string) (string) {
  status := "HTTP/1.1 500 Internal Server Error\n\n"
  headers := "Content-Type: text/plain\r\nContent-Length: %d\r\n\r\n"
  error_string := fmt.Sprintf(status + headers + "%s", len(message), message)
  return error_string
}
