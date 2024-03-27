package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

var (
	dataStore = make(map[string]string)
	lock      = sync.RWMutex{}
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			break
		}
		processCommand(conn, message)
	}
}

func processCommand(conn net.Conn, commandStr string) {
	commandStr = strings.TrimSpace(commandStr)
	parts := strings.SplitN(commandStr, " ", 4)
	if len(parts) < 3 {
		fmt.Fprintf(conn, "ERROR\n")
		return
	}

	id, command, key := parts[0], parts[1], parts[2]
	value := ""
	if len(parts) == 4 {
		value = parts[3]
	}

	fmt.Println("Command:",  command);
	

	switch command {
		case "SET":
			lock.Lock()
			dataStore[key] = value
			lock.Unlock()
			fmt.Fprintf(conn, "%s OK\n", id)
		case "GET":
			lock.RLock()
			val, ok := dataStore[key]
			lock.RUnlock()
			if !ok {
				fmt.Fprintf(conn, "%s ERROR NOT FOUND\n", id)
				return
			}
			fmt.Fprintf(conn, "%s OK %s\n", id, val)
		default:
			fmt.Fprintf(conn, "%s ERROR UNKNOWN COMMAND\n", id)
		}
}

func main() {
	ln, err := net.Listen("tcp", ":7070")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer ln.Close()
	fmt.Println("Cache server running on port 7070")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			return
		}
		go handleConnection(conn)
	}
}
