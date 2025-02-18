package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

var defaultPaths = map[string][]string{
	"windows": {
		`C:\Program Files`,
		`C:\Program Files (x86)`,
		`C:\Windows\System32`,
		`C:\Windows\SysWOW64`,
	},
	"linux": {
		"/usr/bin",
		"/usr/local/bin",
		"/lib",
		"/usr/lib",
		"/usr/local/lib",
		"/opt",
	},
}

func saveResults(results interface{}, outputFile string, startTime time.Time) {
	f, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(results); err != nil {
		log.Fatalf("Failed to encode JSON: %v", err)
	}

	totalTime := time.Since(startTime)
	fmt.Printf("Scan complete. Information saved to %s (Total time: %v)\n", outputFile, totalTime)
}
