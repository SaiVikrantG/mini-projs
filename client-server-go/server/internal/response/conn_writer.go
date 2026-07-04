package response

import (
	"fmt"
	"net"
	"strconv"
)

type responseWriter struct {
	conn    net.Conn
	headers map[string]string
}

func NewResponseWriter(conn net.Conn) ResponseWriter {
	return &responseWriter{
		conn:    conn,
		headers: make(map[string]string),
	}
}

func (rw *responseWriter) Header() map[string]string {
	return rw.headers
}

func (rw *responseWriter) Write(statusCode int, body []byte) (int, error) {
	fmt.Fprintf(rw.conn, "HTTP/1.1 %d %s\r\n", statusCode, statusText(statusCode))
	rw.headers["Content-Length"] = strconv.Itoa(len(body))
	for k, v := range rw.headers {
		fmt.Fprintf(rw.conn, "%s: %s\r\n", k, v)
	}
	fmt.Fprintf(rw.conn, "\r\n")
	return rw.conn.Write(body)
}

func statusText(code int) string {
	switch code {
	case 200:
		return "OK"
	case 400:
		return "Bad Request"
	case 404:
		return "Not Found"
	case 413:
		return "Content Too Large"
	case 500:
		return "Internal Server Error"
	case 505:
		return "HTTP Version Not Supported"
	default:
		return "Unknown"
	}
}
