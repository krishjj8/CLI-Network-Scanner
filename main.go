package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	target := "scanme.nmap.org:80"

	conn, err := net.DialTimeout("tcp", target, 3*time.Second)

	if err != nil {
		fmt.Printf("[-] Connection to %s failed: %v\n", target, err)
		return
	}

	defer conn.Close()

	fmt.Printf("[+] Connection to %s succesful!", target)
}
