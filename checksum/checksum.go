package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const IP_FILE = "tcp_addrs_%d.txt"
const IP_DATA = "tcp_data_%d.dat"
const TCP_PROTOCOL = 6

func ip_addr_byte(ip_addr string) ([]byte, bool) {
	parts := strings.Split(ip_addr, ".")
	if len(parts) != 4 {
		fmt.Printf("IP address %v does not have 4 octets\n", ip_addr)
		return nil, false
	}

	var ipBytes []byte
	for _, ip := range parts {

		n, err := strconv.Atoi(ip)
		if err != nil || n < 0 || n > 255 {
			fmt.Printf("Cant convert IP Address octet %v to integer\n", ip)
			return nil, false
		}
		ipBytes = append(ipBytes, byte(n))
	}

	return ipBytes, true
}

// checksum computes the Internet checksum (RFC 793 / RFC 1071) over data.
func checksum(data []byte) uint16 {
	var sum uint32

	length := len(data)
	for i := 0; i+1 < length; i += 2 {
		sum += uint32(data[i])<<8 | uint32(data[i+1])
	}
	if length%2 == 1 {
		sum += uint32(data[length-1]) << 8
	}

	for sum>>16 > 0 {
		sum = (sum & 0xFFFF) + (sum >> 16)
	}

	return ^uint16(sum)
}

func compute_ip_header(ip_file string, ip_data string) ([]byte, bool) {
	data, err := os.ReadFile(ip_file)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, false
	}

	ip_strings := string(data)
	ip_strings = strings.TrimSpace(ip_strings)
	ip_addresses := strings.Split(ip_strings, " ")
	if len(ip_addresses) != 2 {
		fmt.Printf("Expected 2 IP addresses in %v, got %v\n", ip_file, len(ip_addresses))
		return nil, false
	}

	source_ip_address, ok := ip_addr_byte(ip_addresses[0]) // eg. c6 33 64 4d
	if !ok {
		return nil, false
	}

	dest_ip_address, ok := ip_addr_byte(ip_addresses[1])
	if !ok {
		return nil, false
	}

	tcp_data, err := os.ReadFile(ip_data)
	if err != nil {
		fmt.Println("Error in reading data file: ", err)
		return nil, false
	}
	len_data := len(tcp_data)
	if len_data < 20 {
		fmt.Printf("TCP segment in %v is too short (%v bytes)\n", ip_data, len_data)
		return nil, false
	}

	original_checksum := binary.BigEndian.Uint16(tcp_data[16:18])

	tcp_zero_chksum := make([]byte, len_data)
	copy(tcp_zero_chksum, tcp_data)
	tcp_zero_chksum[16] = 0
	tcp_zero_chksum[17] = 0

	tcp_length := make([]byte, 2)
	binary.BigEndian.PutUint16(tcp_length, uint16(len_data))

	pseudo_header := bytes.Join([][]byte{
		source_ip_address,
		dest_ip_address,
		{0},
		{TCP_PROTOCOL},
		tcp_length,
		tcp_zero_chksum,
	}, []byte{})

	computed_checksum := checksum(pseudo_header)
	valid := computed_checksum == original_checksum

	if valid {
		fmt.Printf("%v: checksum OK (0x%04x)\n", ip_data, original_checksum)
	} else {
		fmt.Printf("%v: checksum MISMATCH (header 0x%04x, computed 0x%04x)\n", ip_data, original_checksum, computed_checksum)
	}

	return tcp_data, valid
}

func main() {
	for i := range 10 {
		ip_file := fmt.Sprintf(IP_FILE, i)
		ip_data := fmt.Sprintf(IP_DATA, i)

		compute_ip_header(ip_file, ip_data)
	}
}
