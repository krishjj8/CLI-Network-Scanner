package main

import (
	"flag"
	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)

func scanport(host string, port int, wg *sync.WaitGroup, results chan int) {
	defer wg.Done()
	target := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", target, 3*time.Second)
	if err != nil {
		return
	}
	defer conn.Close()
	results <- port

}

func main() {
	//flags
	host := flag.String("host", "scanme.nmap.org", "Host to scan")
	startPort := flag.Int("start", 1, "Start of port range")
	endPort := flag.Int("end", 1024, "end of port range")
	flag.Parse()

	var wg sync.WaitGroup
	results := make(chan int)
	var openPorts []int

	go func() {
		for p := range results {
			openPorts = append(openPorts, p)
		}
	}()
	fmt.Printf("Scanning %s from %d to %d...\n", *host, *startPort, *endPort)

	for port := *startPort; port <= *endPort; port++ {
		wg.Add(1)
		go scanport(*host, port, &wg, results)
	}
	wg.Wait()
	close(results)

	sort.Ints(openPorts)
	for _, port := range openPorts {
		fmt.Printf("[+] Port %d is open\n", port)

	}
	fmt.Println("Scan complete")
}
