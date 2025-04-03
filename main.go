// Filename: main.go
// Purpose: This program demonstrates how to create a TCP network connection using Go

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ScanResult struct {
	Target string `json:"target"`
	Port   int    `json:"port"`
	Status string `json:"status"`
	Banner string `json:"banner,omitempty"`
}

func worker(wg *sync.WaitGroup, tasks chan string, dialer net.Dialer, openPorts *int, mu *sync.Mutex, progress *int, totalPorts int, results *[]ScanResult) {
	defer wg.Done()
	maxRetries := 3
    for addr := range tasks {
		var success bool
		var banner string
		
		for i := range maxRetries {      
		conn, err := dialer.Dial("tcp", addr)// Attempt to connect to the address
		if err == nil {
			defer conn.Close()
			fmt.Printf("Connection to %s was successful\n", addr)
			mu.Lock() // Lock the mutex to safely update openPorts
			*openPorts++ // Increment the open ports counter
			mu.Unlock() // Unlock the mutex
			success = true

			//Banner grabbing
			buffer := make([]byte, 1024) // Create a buffer to read data
			conn.SetReadDeadline(time.Now().Add(2 * time.Second)) // Set a read deadline
			n, err := conn.Read(buffer) // Read data from the connection
			if err == nil {
				banner = string(buffer[:n]) // Convert the buffer to a string
				fmt.Printf("Banner from %s: %s\n", addr, string(buffer[:n])) // Print the banner
			} else {
				fmt.Printf("Failed to read banner from %s: %v\n", addr,err) // Print the error
				
			}
			break
		}
		backoff := time.Duration(1<<i) * time.Second // Exponential backoff
		fmt.Printf("Attempt %d to %s failed. Waiting %v...\n", i+1,  addr, backoff)
		time.Sleep(backoff)
	    }
		if !success {
			fmt.Printf("Failed to connect to %s after %d attempts\n", addr, maxRetries)
		}

		// Update and print progress
        mu.Lock()
        *progress++
        host, port, _ := net.SplitHostPort(addr)
        portNum, _ := strconv.Atoi(port)
        status := "closed"
        if success {
            status = "open"
        }
        *results = append(*results, ScanResult{Target: host, Port: portNum, Status: status, Banner: banner})
		fmt.Printf("Scanning port %d/%d (%d%% complete)\n", *progress, totalPorts, (*progress*100)/totalPorts) // Print progress
		mu.Unlock() // Unlock the mutex
	}

}
func main() {

	var wg sync.WaitGroup
	tasks := make(chan string, 100)

	//Define and parse the target flag
    targets := flag.String("targets","scanme.nmap.org", "Specify the IP address or hostname to scan")
	startPort := flag.Int("start-port", 1, "Specify the starting port")
	endPort := flag.Int("end-port", 1024, "Specify the ending port")
	workers := flag.Int("workers", 100, "Specify the number of concurrent workers")
	timeout := flag.Int("timeout", 5, "Specify the connection timeout in seconds")
	jsonOutput := flag.Bool("json", false, "Output results in JSON format")

	flag.Parse() // Parse the command line flags

	dialer := net.Dialer { // Create a new dialer 
		Timeout: time.Duration(*timeout) * time.Second,
	}

	var openPorts int // Counter for open ports
	var progress int // Progress counter
	var mu sync.Mutex // Mutex for thread-safe access to openPorts
	var results []ScanResult // Slice to store scan results

	startTime := time.Now() // Record the start time
	totalPorts := *endPort - *startPort + 1 // Calculate the total number of ports scanned

	targetList := strings.Split(*targets, ",") // Split the target string into a list of targets

    for i := 1; i <= *workers; i++ {// Create worker goroutines
		wg.Add(1)
		go worker(&wg, tasks, dialer, &openPorts, &mu, &progress, totalPorts*len(targetList), &results) // Start a worker goroutine
	}

	// Loop through the specified port range and send addresses to the tasks channel
	for _, target := range targetList {	
		for p := *startPort; p <= *endPort; p++ {
			port := strconv.Itoa(p)
			address := net.JoinHostPort(target, port) // Combine IP and port
			tasks <- address// Send the address to the tasks channel
		}
	}
	// Close the tasks channel after all addresses have been sent
	close(tasks)
	wg.Wait() // Wait for all workers to finish

	duration := time.Since(startTime) // Calculate the duration

	if *jsonOutput {
        jsonData, err := json.MarshalIndent(results, "", "  ")
        if err != nil {
            fmt.Printf("Error generating JSON output: %v\n", err)
            return
        }
        fmt.Println(string(jsonData))
    } else {
        fmt.Printf("\nScan Summary:\n")
        fmt.Printf("Number of open ports: %d\n", openPorts)
        fmt.Printf("Time taken: %v\n", duration)
        fmt.Printf("Total ports scanned: %d\n", totalPorts*len(targetList))
    }
}