package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

var (
	peersMutex sync.Mutex
	peers      = make(map[string]net.Conn)
)

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func StartP2PServer(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("P2P listen error:", err)
		return
	}
	fmt.Println("P2P server listening on", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	addr := conn.RemoteAddr().String()
	peersMutex.Lock()
	peers[addr] = conn
	peersMutex.Unlock()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		var msg Message
		_ = json.Unmarshal(scanner.Bytes(), &msg)
		switch msg.Type {
		case "Handshake":
			fmt.Println("Handshake from", addr)
		case "NewBlock":
			fmt.Println("Received NewBlock from", addr)
		case "NewTransaction":
			fmt.Println("Received NewTransaction from", addr)
		default:
			fmt.Println("Unknown message from", addr, msg.Type)
		}
	}
	peersMutex.Lock()
	delete(peers, addr)
	peersMutex.Unlock()
}

func Broadcast(msg Message) {
	// connect to localhost server to simulate broadcast
	conn, err := net.Dial("tcp", ":40404")
	if err != nil {
		return
	}
	defer conn.Close()
	b, _ := json.Marshal(msg)
	fmt.Fprintln(conn, string(b))
}

func NetworkStatus() string {
	peersMutex.Lock()
	defer peersMutex.Unlock()
	return fmt.Sprintf("peers=%d", len(peers))
}
