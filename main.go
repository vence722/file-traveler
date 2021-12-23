package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fileTravelerServer()
		return
	}
	if len(os.Args) < 3 {
		fmt.Println("Usage: for server: 'file-traveler', for client: 'file-traveler <file-path> <target-host-name>'")
		os.Exit(1)
	}
	filePath := os.Args[1]
	targetHostName := os.Args[2]
	fileTravelerClient(filePath, targetHostName)
}
