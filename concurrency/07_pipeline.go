package main

import "fmt"

// In a pipeline, data flows through multiple stages where each stage performs an operation.
// Each stage communicates via channels and typically runs concurrently using goroutines.
//
// Note: An unbuffered channel has no capacity (buffer size 0). A send on an unbuffered channel
// blocks until a receiver is ready to receive the value.

// sliceToChannel converts a slice of integers into a receive-only channel.
// It launches a goroutine that sends each integer from the slice into the channel,
// then closes the channel to signal that no more data will be sent.
func sliceToChannel(nums []int) <-chan int {
	out := make(chan int) // unbuffered channel (synchronous communication)
	go func() {
		for _, n := range nums {
			out <- n // Blocks until a receiver is ready to read the value
		}
		close(out) // Close the channel to indicate completion
	}()
	return out
}

// sq reads integers from the input channel, squares each one,
// and sends the result to a new output channel. It then closes the output channel.
func sq(in <-chan int) <-chan int {
	out := make(chan int) // unbuffered channel for synchronous communication
	go func() {
		for n := range in {
			out <- n * n // Sends the squared value; blocks until the value is received
		}
		close(out) // Close the output channel once all input values are processed
	}()
	return out
}

func main() {
	// Input: A slice of numbers
	nums := []int{2, 3, 4, 5, 6}

	// Stage 1: Convert the slice to a channel.
	dataChannel := sliceToChannel(nums)

	// Stage 2: Process the numbers by squaring each one.
	// Note: We call sq only once here to create the pipeline stage.
	finalChannel := sq(dataChannel)

	// Stage 3: Read from the final channel and print each squared number.
	for n := range finalChannel {
		fmt.Println(n)
	}
}
