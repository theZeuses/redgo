package main

import (
	"fmt"
	"net"
)

func main() {
	serverAddr := "127.0.0.1:7000"
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Failed to connect to Redgo:", err)
		return
	}
	defer conn.Close()

	cli := NewCLI(conn, serverAddr)
	cli.Run()
}
