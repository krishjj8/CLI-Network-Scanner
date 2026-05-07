package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/netip"
	"sort"
	"strings"
	"sync"
	"time"
)

func scanport(ip netip.Addr, port int, wg *sync.WaitGroup, results chan<- string, semaphore chan struct{}) {
	defer wg.Done()
	semaphore <- struct{}{}
	defer func() { <-semaphore }()
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, 500*time.Millisecond)
	if err != nil {
		return
	}
	defer conn.Close()
	results <- fmt.Sprintf("%s:%d", ip, port)

}

func main() {
	//flags
	cidrinput := flag.String("cidr", "45.33.32.156/32", "Network Range to scan")
	startPort := flag.Int("start", 80, "Start port")
	endPort := flag.Int("end", 85, "end port")
	concurrency := flag.Int("concurrency", 500, "Max concurrent goroutines")
	jsonOutput := flag.Bool("json", false, "Output results in JSON format")
	flag.Parse()

	semaphore := make(chan struct{}, *concurrency)

	prefix, err := netip.ParsePrefix(*cidrinput)
	if err != nil {
		fmt.Printf("Invalid CIDR:%v", err)
		return

	}

	var wg sync.WaitGroup
	results := make(chan string)
	var openPorts []string

	go func() {
		for p := range results {
			openPorts = append(openPorts, p)
		}
	}()
	fmt.Printf("Scanning %s from %d to %d...\n", *cidrinput, *startPort, *endPort)

	for addr := prefix.Addr(); prefix.Contains(addr); addr = addr.Next() {
		for port := *startPort; port <= *endPort; port++ {
			wg.Add(1)
			go scanport(addr, port, &wg, results, semaphore)
		}
	}
	wg.Wait()
	close(results)

	sort.Strings(openPorts)
	if *jsonOutput {
		data, err := json.MarshalIndent(openPorts, "", "  ")
		if err != nil {
			fmt.Printf("Error generating JSON: %v\n", err)
			return
		}
		fmt.Println(string(data))
	} else {
		fmt.Println("\n--- Scan Results ---")
		for _, port := range openPorts {
			parts := strings.Split(port, ":")
			ip := parts[0]

			hostnames, err := net.LookupAddr(ip)
			if err != nil || len(hostnames) == 0 {
				fmt.Printf("[+] %s is open\n", port)
				continue
			}
			hostname := strings.TrimSuffix(hostnames[0], ".")
			fmt.Printf("[+] %s is open (%s)\n", port, hostname)
		}
	}

}
