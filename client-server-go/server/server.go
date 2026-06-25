package main

import (
	"flag"
	"fmt"
	"net"
)

func server_init(port *int) {
	conn_port := fmt.Sprintf(":%v", *port)
	listener, err := net.Listen("tcp", conn_port)
	if err != nil {
		fmt.Printf("Cant start a server on port %v", *port)
		return
	}
	defer listener.Close()

	fmt.Printf("Server started at port %v\n", *port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting conenctions:", err)
			continue
		}

		handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
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

func main() {
	port := flag.Int("port", 12399, "Port to start the server on")

	flag.Parse()

	server_init(port)
}
