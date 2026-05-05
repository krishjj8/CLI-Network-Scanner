package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func scanport(host string, port int, wg *sync.WaitGroup) {
	defer wg.Done()
	target := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", target, 3*time.Second)
	if err != nil {
		return
	}
	defer conn.Close()
	fmt.Printf("[+] Port %d is open\n", port)
}

func main() {
	host := "scanme.nmap.org"
	var wg sync.WaitGroup
	for port := 1; port <= 1024; port++ {
		wg.Add(1)
		go scanport(host, port, &wg)
	}
	wg.Wait()
	fmt.Println("Scan complete")
}
