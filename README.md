# TCP Network Connection Tool
## Port-scanner

## Description

This tool is a TCP network scanner written in Go. It scans specified ports on given targets to check if they are open and optionally grabs banners. It supports concurrent scanning with configurable worker count and timeout settings.

## Features

- Scan a range of ports or specific ports on a target.
- Concurrent scanning with configurable worker count.
- Banner grabbing.
- JSON output format.

## How to Build and Run

### Building the Tool

1. Clone the repository:
    ```sh
    git clone https://github.com/Juana-Tut/port-scanner.git

2. Build the tool:
    ```sh
    go run main.go -targets=scanme.nmap.org,google.com -ports=22,80,100,150 -workers=100 -timeout=5 -json

### Sample Output

Scanning port 8/8 (100% complete)
[
  {
    "target": "scanme.nmap.org",
    "port": 22,
    "status": "open",
    "banner": "SSH-2.0-OpenSSH_6.6.1p1 Ubuntu-2ubuntu2.13\r\n"
  },
  {
    "target": "google.com",
    "port": 80,
    "status": "open"
  },
  {
    "target": "scanme.nmap.org",
    "port": 80,
    "status": "open"
  },
  
### Link to Video
(https://youtu.be/aOaX338akYg)
