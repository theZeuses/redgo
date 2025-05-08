package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Initiating RedGo....")

	listener, err := net.Listen("tcp", ":7000")

	if err != nil {
		fmt.Println(err)
		return
	}

	aof, err := InitAof()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()

	fmt.Println("Listening on port 7000...")

	for {
		connection, err := listener.Accept()

		if err != nil {
			fmt.Println(err)
			return
		}

		go handleConnection(connection, aof)
	}
}

func handleConnection(conn net.Conn, aof *Aof) {
	clientID := conn.RemoteAddr().String()
	reader := NewReader(conn)
	writer := NewWriter(conn)

	client := &Client{
		ID:            clientID,
		Conn:          conn,
		Subscriptions: make(map[string]*PubSubChannel),
		Reader:        reader,
		Writer:        writer,
	}

	ClientsMutex.Lock()
	Clients[clientID] = client
	ClientsMutex.Unlock()

	fmt.Println("New client connected:", clientID)

	defer func() {
		conn.Close()
		ClientsMutex.Lock()
		delete(Clients, clientID)
		ClientsMutex.Unlock()
		unsubscribeAll(client)
		fmt.Println("Client disconnected:", clientID)
	}()

	for {
		err := Handle(*client, aof)
		if err != nil {
			break
		}
	}
}
