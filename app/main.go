package main

import (
	"fmt"
	"log"

	"github.com/lmalunin/toolkit"
)

func main() {
	// Create a new instance of the Tools struct
	tools := toolkit.Tools{}

	// Generate a random string of length 10
	randomStr := tools.RandomString(10)
	fmt.Printf("Generated random string: %s\n", randomStr)

	// Example of using the random string for something like a token
	token := generateToken()
	fmt.Printf("Generated token: %s\n", token)

	log.Println("Application completed successfully")
}

// generateToken creates a secure token using the toolkit's RandomString function
func generateToken() string {
	tools := toolkit.Tools{}
	// Generate a longer string for security
	return tools.RandomString(32)
}
