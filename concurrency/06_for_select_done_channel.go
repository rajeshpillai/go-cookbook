package main

import (
	"fmt"
	"time"
)

/*
NOTES:
Explanation:
Channel for Cancellation:
The done channel is created in main and passed to the doWork function. It is used solely for signaling cancellation, which is why it is declared as a read-only channel (<-chan bool) in doWork.

Goroutine Execution:
The doWork function is run as a goroutine. It uses a select statement inside an infinite loop. The select listens for a signal on the done channel. When done is closed (or a value is sent), the goroutine exits.

Simulating Work:
The default case in the select statement allows the goroutine to continue doing its work (in this example, printing "DOING WORK>>>") when no cancellation signal is present.

Cancellation after Delay:
In the main function, after sleeping for 3 seconds, the done channel is closed. This signals the worker goroutine to stop, thus cleanly terminating its execution.
*/

// doWork continuously performs work until a signal is received on the done channel.
// The done channel is a read-only channel used to signal cancellation.
// When a value is received or the channel is closed, the goroutine stops working.
func doWork(done <-chan bool) {
	for {
		select {
		// If a signal is received on done (or the channel is closed), exit the function.
		case <-done:
			// Exit the goroutine immediately.
			return
		// The default case executes if no signal is received; here we simulate doing some work.
		default:
			fmt.Println("DOING WORK>>>")
		}
	}
}

func main() {
	// Create a channel to signal when the work should be stopped.
	done := make(chan bool)

	// Launch doWork in a new goroutine.
	go doWork(done)

	// Let the goroutine run for 3 seconds.
	time.Sleep(time.Second * 3)

	// Close the done channel to signal cancellation to the doWork goroutine.
	// Closing the channel notifies all receivers that no further values will be sent.
	close(done)
}
