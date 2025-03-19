package main

import (
	"fmt" // Import fmt package for formatted I/O operations.
)

/*
NOTES:
Explanation:
Channel Creation:
myChannel := make(chan string) creates a new unbuffered channel for strings. Unbuffered channels require both sender and receiver to be ready for the communication to occur.

Goroutine Communication:
The anonymous goroutine sends a string into the channel. The main goroutine receives that string using <-myChannel.

Blocking Behavior:
The send operation (myChannel <- "my data") and the receive operation (msg := <-myChannel) are blocking, ensuring synchronization between goroutines.
*/

// In Go, channels allow goroutines to communicate in a FIFO (first-in, first-out) manner.
func main() { // main() is the entry point and runs in the main goroutine.
	// Create a channel of type string. This channel will be used for sending and receiving string data.
	myChannel := make(chan string)

	// Launch an anonymous goroutine that sends data into the channel.
	go func() {
		// Send "my data" into myChannel. This send operation will block until another goroutine is ready to receive.
		myChannel <- "my data"
	}()

	// Receive a value from myChannel.
	// This receive operation is blocking, meaning it will wait here until a value is available.
	msg := <-myChannel

	// Print the received message to the console.
	fmt.Println(msg)
}
