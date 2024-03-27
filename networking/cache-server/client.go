package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"
)

func sendCommand(conn net.Conn, command string) (string, error) {
	fmt.Fprintf(conn, command+"\n")
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return "", err
	}
	return response, nil
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(i int) {
					defer wg.Done()

					conn, err := net.Dial("tcp", "localhost:7070")
					if err != nil {
							fmt.Println("Error connecting:", err.Error())
							return
					}
					defer conn.Close()

					key := fmt.Sprintf("key%d", i)
					value := fmt.Sprintf("value%d", i)
					setCmd := fmt.Sprintf("%d SET %s %s", i, key, value)
					getCmd := fmt.Sprintf("%d GET %s", i, key)

					if response, err := sendCommand(conn, setCmd); err != nil {
							fmt.Println("Error:", err)
					} else {
							fmt.Print("Set Response: ", response)
					}
					
					time.Sleep(100 * time.Millisecond) // Simulate delay

					if response, err := sendCommand(conn, getCmd); err != nil {
							fmt.Println("Error:", err)
					} else {
							fmt.Print("Get Response: ", response)
					}
			}(i)
	}
	wg.Wait()
}
