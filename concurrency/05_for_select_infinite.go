package main

import (
	"fmt"
	"time"
)

func main() {

	// Infinitely running goroutine
	// FIX: in next example using readonly 'done' channel
	go func() {
		for {
			select {
			default:
				fmt.Println("DOING WORK>>>")
			}
		}
	}()

	// time.Sleep(time.Second * 10)
	time.Sleep(time.Hour * 299)
}
