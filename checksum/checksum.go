package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const IP_FILE = "tcp_addrs_%d.txt"
const IP_DATA = "tcp_data_%d.dat"

func ip_addr_byte(ip_addr string) (string, bool) {
	parts := strings.Split(ip_addr, ".")

	var ipBytes []byte
	for _, ip := range parts {

		n, err := strconv.Atoi(ip)
		if err != nil {
			fmt.Printf("Cant convert IP Address octet %v to integer\n", ip)
			return "", false
		}
		ipBytes = append(ipBytes, byte(n))
	}

	return fmt.Sprintf("% x", ipBytes), true
}

func compute_ip_header(ip_file string) ([]byte, bool) {
	data, err := os.ReadFile(ip_file)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, false
	}

	ip_strings := string(data)
	ip_strings = strings.TrimSpace(ip_strings)
	ip_addresses := strings.Split(ip_strings, " ")

	source_ip_address, ok := ip_addr_byte(ip_addresses[0])
	if !ok {
		return nil, false
	}

	dest_ip_address, ok := ip_addr_byte(ip_addresses[1])
	if !ok {
		return nil, false
	}

	fmt.Println(source_ip_address, dest_ip_address)

	return nil, false
}

func main() {
	for i := range 10 {
		ip_file := fmt.Sprintf(IP_FILE, i)
		// ip_data := fmt.Sprintf(IP_DATA, i)

		compute_ip_header(ip_file)
	}
}
