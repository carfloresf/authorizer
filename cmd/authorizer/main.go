package main

import (
	"fmt"
	"os"

	cmd "authorizer/cmd/root"
	"authorizer/internal/app/service"
	"authorizer/internal/app/storage"
	"authorizer/internal/common/logfile"
)

func main() {
	logfile.Init()

	if len(os.Args) > 1 {
		switch os.Args[1:][0] {
		case "version":
			fmt.Println("v1.0")

		case "help":
			fmt.Println("send file with transactions to stdin")
		}

		os.Exit(0)
	}

	// Initialize DB
	db := storage.InMemory{}

	// Initialize service
	svc := service.New(&db)

	// Get input from stdin
	stdin := os.Stdin
	// Return output to stdout
	stdout := os.Stdout

	// Execute application
	cmd.Execute(svc, stdin, stdout)
}
