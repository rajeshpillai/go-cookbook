package main

import (
	"fmt"  // Import fmt for printing output to the console.
	"time" // Import time for using time-related functions like Sleep.
)

/*
NOTES:
Explanation:
Goroutines:
The go keyword is used to start new goroutines, allowing someFunc to run concurrently.
Synchronization:
time.Sleep is used to pause the main function so that the goroutines have time to complete their execution.
Functionality:
Each goroutine prints a different string, and after a delay, the main function prints "hi"ss
*/

// someFunc takes a string and prints it.
// This function is designed to be executed as a goroutine.
func someFunc(num string) {
	fmt.Println(num) // Print the passed number.
}

func main() {
	// Start three goroutines concurrently.
	// Each goroutine calls someFunc with a different string.
	go someFunc("1")
	go someFunc("2")
	go someFunc("3")

	// Sleep for 2 seconds to ensure that the goroutines have enough time
	// to execute and print their outputs before the main function exits.
	time.Sleep(time.Second * 2)

	// After waiting, print "hi" to the console.
	fmt.Println("hi")
}
