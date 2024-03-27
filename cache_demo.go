package main

import (
	"fmt"
	"log"
	"go-cookbook/networking/cache_client" 
)

func main() {
	client, err := cache_client.NewCacheClient("localhost:7070")
	if err != nil {
		log.Fatal("Failed to create cache client:", err)
	}
	fmt.Println("Client connected to cache server");
	defer client.Close()

	// Set a value
	setResponse, err := client.Set("testKey", "Hello, World!")
	if err != nil {
		log.Fatal("Failed to set value:", err)
	}
	fmt.Println("Set response:", setResponse)

	// Get the value
	getResponse, err := client.Get("testKey")
	if err != nil {
		log.Fatal("Failed to get value:", err)
	}
	fmt.Println("Get response:", getResponse)

	// Delete the value
	delResponse, err := client.Del("testKey")
	if err != nil {
		log.Fatal("Failed to delete value:", err)
	}
	fmt.Println("Delete response:", delResponse)
}
