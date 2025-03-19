// Allows main to wait on channel operations
package main

import "fmt"

/*
NOTES:
Explanation:
Buffered Channel Creation:
The channel charChannel is created with a capacity of 3, allowing asynchronous communication without immediate blocking until the buffer is full.

Sending Values:
The for loop iterates over the chars slice and uses a select statement to send each value into the channel. Although there's only one case in the select, it still performs the send operation in a non-blocking style if other cases were present.

Channel Closure and Reading:
The channel is closed after sending all the values. The subsequent for loop reads and prints each buffered value until the channel is empty.
*/

// In Go, channels are FIFO queues that enable communication between goroutines.
// Buffered channels allow asynchronous communication by holding a fixed number of values.
func main() { // main goroutine
	// Create a buffered channel for strings with a capacity of 3.
	// This means that up to 3 values can be sent without requiring an immediate receiver.
	charChannel := make(chan string, 3)

	// Define a slice of characters to be sent into the channel.
	chars := []string{"a", "b", "c"}

	// Iterate over the slice and send each character to the channel.
	// The select statement here, although containing a single case,
	// is used to perform a non-blocking send operation.
	for _, s := range chars {
		select {
		// Attempt to send the string into charChannel.
		// Since the channel is buffered, these sends will succeed until the buffer is full.
		case charChannel <- s:
		}
	}
	// After sending all values, close the channel to signal that no more data will be sent.
	close(charChannel)

	// Even though the channel is closed, any remaining buffered values can still be read.
	// The range loop will iterate over the channel until it is empty.
	for result := range charChannel {
		// Print each received character.
		fmt.Println(result)
	}
}
