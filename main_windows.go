//go:build windows

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"getProducts/pkg/windows"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <output_file>\n", os.Args[0])
		os.Exit(1)
	}

	startTime := time.Now()
	outputFile := os.Args[1]

	results, err := windows.ScanSystem(defaultPaths["windows"])
	if err != nil {
		log.Fatalf("Error scanning directories: %v", err)
	}

	saveResults(results, outputFile, startTime)
}
