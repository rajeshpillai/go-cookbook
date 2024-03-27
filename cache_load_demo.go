package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"
)
import "sync/atomic"

// Client operation simulates interactions with the cache server.
func clientOperation(clientID int, wg *sync.WaitGroup, serverAddress string, ops *int64) {
	defer wg.Done()

	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Printf("Client %d: Error connecting to server: %v\n", clientID, err)
		return
	}
	defer conn.Close()

	for i := 0; i < 10; i++ {
		start := time.Now()

		// Simulate a SET operation
		setCommand := fmt.Sprintf("%d SET key%d value%d\n", clientID, i, i)
		fmt.Fprintf(conn, setCommand)

		// Wait for the response
		reader := bufio.NewReader(conn)
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Client %d: Error reading response: %v\n", clientID, err)
			return
		}
		fmt.Printf("Client %d: %s", clientID, response)

		// Simulate a GET operation
		getCommand := fmt.Sprintf("%d GET key%d\n", clientID, i)
		fmt.Fprintf(conn, getCommand)

		// Wait for the response
		response, err = reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Client %d: Error reading response: %v\n", clientID, err)
			return
		}
		fmt.Printf("Client %d: %s", clientID, response)

		end := time.Now()
		duration := end.Sub(start)
		fmt.Printf("Client %d: Operation %d took %v\n", clientID, i, duration)

		// Increment the operations counter
		atomic.AddInt64(ops, 2) // Each loop iteration sends a SET and a GET, totaling 2 operations.

		// Sleep to simulate think time
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	serverAddress := "localhost:7070"
	var wg sync.WaitGroup
	var ops int64 // Total operations performed

	clientCount := 10 // Number of concurrent clients
	testStart := time.Now()

	for i := 0; i < clientCount; i++ {
		wg.Add(1)
		go clientOperation(i, &wg, serverAddress, &ops)
	}

	wg.Wait()

	testEnd := time.Now()
	testDuration := testEnd.Sub(testStart)
	fmt.Printf("Load test completed in %v.\n", testDuration)

	totalOps := atomic.LoadInt64(&ops)
	fmt.Printf("Total operations: %d\n", totalOps)
	rps := float64(totalOps) / testDuration.Seconds()
	fmt.Printf("Requests per second (RPS): %.2f\n", rps)
}
