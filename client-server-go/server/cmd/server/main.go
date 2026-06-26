// Parse that request header to get the file name.
// Strip the path off for security reasons.
// Read the data from the named file.
// Determine the type of data in the file, HTML or text.
// Build an HTTP response packet with the file data in the payload.
// Send that HTTP response back to the client.

package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/SaiVikrantG/server/internal/handlers"
)

func main() {
	port_num := flag.Int("port", 28333, "Port to start the server on")

	flag.Parse()
	port := *port_num

	conn_port := fmt.Sprintf(":%v", port)
	listener, err := net.Listen("tcp", conn_port)
	if err != nil {
		fmt.Printf("Cant start a server on port %v", port)
		return
	}
	defer listener.Close()

	fmt.Printf("Server started at port %v\n", port)

	for {
		conn, err := listener.Accept()
		fmt.Println(conn.RemoteAddr())

		if err != nil {
			fmt.Println("Error accepting conenctions:", err)
			continue
		}

		handlers.HandleConn(conn)
	}

}
