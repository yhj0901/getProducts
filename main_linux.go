//go:build linux

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"getProducts/pkg/linux"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <output_file>\n", os.Args[0])
		os.Exit(1)
	}

	startTime := time.Now()
	outputFile := os.Args[1]

	results, err := linux.ScanSystem(defaultPaths["linux"])
	if err != nil {
		log.Fatalf("Error scanning directories: %v", err)
	}

	saveResults(results, outputFile, startTime)
}
