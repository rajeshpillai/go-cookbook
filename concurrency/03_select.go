// Allows main to wait on multiple channels
package main

import (
	"fmt" // Import fmt package for printing output.
)

/*
NOTES:
Explanation:
Channel Creation:
Two unbuffered channels (myChannel and anotherChannel) are created for communication between goroutines.

Goroutines:
Two anonymous goroutines are launched. Each goroutine sends a specific string to its respective channel.

Select Statement:
The select statement waits for a message from either channel. It blocks until one of the channels receives a value. If both channels are ready, the selection is made randomly.
*/

// Channels in Go act like FIFO queues, allowing communication between goroutines.
// The main goroutine and child goroutines can communicate using channels.
func main() { // main function runs in the main goroutine.
	// Create two unbuffered channels for string data.
	myChannel := make(chan string)
	anotherChannel := make(chan string)

	// Launch a goroutine that sends "rabbit" to myChannel.
	go func() {
		// Send the string "rabbit" to myChannel.
		// This operation will block until the value is received.
		myChannel <- "rabbit"
	}()

	// Launch another goroutine that sends "cow" to anotherChannel.
	go func() {
		// Send the string "cow" to anotherChannel.
		// This operation will block until the value is received.
		anotherChannel <- "cow"
	}()

	// The select statement below blocks until it receives a message from one of the channels.
	// If messages from both channels are ready simultaneously, one is chosen at random.
	select {
	// If a message is received from myChannel, store it in msgFromMyChannel.
	case msgFromMyChannel := <-myChannel:
		// Print the message received from myChannel.
		fmt.Println(msgFromMyChannel)
	// If a message is received from anotherChannel, store it in msgFromAnotherChannel.
	case msgFromAnotherChannel := <-anotherChannel:
		// Print the message received from anotherChannel.
		fmt.Println(msgFromAnotherChannel)
	}
}
