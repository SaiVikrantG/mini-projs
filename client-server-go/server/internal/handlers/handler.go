package handlers

import (
	"fmt"
	"net"
)

func HandleConn(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Error reading: %v\n", err)
		return
	}

	fmt.Printf("Recieved: \n%s", string(buffer[:n]))

	res := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 6\r\nConnection: close\r\n\r\nHello!\r\n"
	fmt.Fprintf(conn, res)
}
