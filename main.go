package main

import (
	"flag"
	"fmt"
	"net"
	"net/netip"
	"sort"
	"sync"
	"time"
)

func scanport(ip netip.Addr, port int, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, 3*time.Second)
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
	flag.Parse()

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
			go scanport(addr, port, &wg, results)
		}
	}
	wg.Wait()
	close(results)

	sort.Strings(openPorts)
	for _, port := range openPorts {
		fmt.Printf("[+] Port %s is open\n", port)

	}
	fmt.Printf("\nScan complete. Found %d open ports.\n", len(openPorts))
}
