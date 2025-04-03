// Filename: main.go
// Purpose: This program demonstrates how to create a TCP network connection using Go

package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)


func worker(wg *sync.WaitGroup, tasks chan string, dialer net.Dialer) {
	defer wg.Done()
	maxRetries := 3
    for addr := range tasks {
		var success bool
		for i := range maxRetries {      
		conn, err := dialer.Dial("tcp", addr)
		if err == nil {
			conn.Close()
			fmt.Printf("Connection to %s was successful\n", addr)
			success = true
			break
		}
		backoff := time.Duration(1<<i) * time.Second
		fmt.Printf("Attempt %d to %s failed. Waiting %v...\n", i+1,  addr, backoff)
		time.Sleep(backoff)
	    }
		if !success {
			fmt.Printf("Failed to connect to %s after %d attempts\n", addr, maxRetries)
		}
	}
}
func main() {

	var wg sync.WaitGroup
	tasks := make(chan string, 100)

	//Define and parse the target flag
    target := flag.String("target","scanme.nmap.org", "Specify the IP address or hostname to scan")
	startPort := flag.Int("start-port", 1, "Specify the starting port")
	endPort := flag.Int("end-port", 1024, "Specify the ending port")
	
	flag.Parse() // Parse the command line flags

	dialer := net.Dialer { // Create a new dialer 
		Timeout: 5 * time.Second,
	}
  
	workers := 100 // Number of concurrent workers

    for i := 1; i <= workers; i++ {// Create worker goroutines
		wg.Add(1)
		go worker(&wg, tasks, dialer)// Start a worker goroutine
	}

	// Loop through the specified port range and send addresses to the tasks channel
	for p := *startPort; p <= *endPort; p++ {
		port := strconv.Itoa(p)
        address := net.JoinHostPort(*target, port) // Combine IP and port
		tasks <- address// Send the address to the tasks channel
	}
	close(tasks)
	wg.Wait()
}