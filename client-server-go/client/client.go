package main

import (
	"flag"
	"fmt"
	"io"
	"net"
)

func send_req(site string) {
	// site = fmt.Sprintf("%s:80", site)
	conn, err := net.Dial("tcp", site)
	if err != nil {
		fmt.Printf("Error connecting to %v: %v", site, err)
		return
	}
	defer conn.Close()

	req := fmt.Sprintf("GET / HTTP/1.1\r\nHost: %s\r\nConnection: close\r\n\r\n", site)
	_, err = fmt.Fprintf(conn, req)
	if err != nil {
		fmt.Printf("Error sending data: %s", err)
	}

	response, err := io.ReadAll(conn)
	if err != nil {
		fmt.Printf("Error reading response: %v", err)
		return
	}

	fmt.Print(string(response))
}

func main() {
	name := flag.String("site", "", "Pass in the site to hit")
	port := flag.Int("port", 80, "Pass the port number for the server to be initialized at")
	flag.Parse()

	url := fmt.Sprintf("%s:%d", *name, *port)

	send_req(url)
}
