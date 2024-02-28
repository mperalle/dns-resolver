package main

import (
	"fmt"
	"net"
	"os"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: ./lookup hostname")
		return
	}

	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:3000")
	if err != nil {
		fmt.Println(err)
		return
	}

	c, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Listen failed:", err)
	}

	defer c.Close()

	_, err = c.Write([]byte(os.Args[1]))
	if err != nil {
		fmt.Println("Error in writing:", err)
	}

	buffer := make([]byte, 1024)
	n, err := c.Read(buffer)
	if err != nil {
		fmt.Println("Read failed:", err)
	}
	fmt.Println("The IP address for", os.Args[1], "is:", string(buffer[:n]))

}
