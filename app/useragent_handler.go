package main

import(
  "fmt"
  "net"
  "net/http"
)

func UserAgentHandler(request *http.Request, conn net.Conn) {
  userAgent := request.UserAgent()
  response_str := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgent), userAgent)
  conn.Write([]byte(response_str))
}
