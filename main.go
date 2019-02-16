package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// helper constants for os.Exit
const (
	OKAY            = 0
	ERROR           = 1
	MaxPossiblePort = 65535
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("No host provided")
		os.Exit(ERROR)
	}
	host := os.Args[1]
	var wg sync.WaitGroup
	for i := 0; i < MaxPossiblePort; i++ {
		go func(i int) {
			wg.Add(1)
			defer wg.Done()
			scan(host, i)
		}(i)
	}
	wg.Wait()
	os.Exit(OKAY)
}

/* scan takes in a host and a port
 * returns true if the port is active
 */
func scan(host string, port int) bool {
	addr := fmt.Sprintf("%s:%d", host, port)
	active := false
	protocols := []string{"tcp", "udp"}
	for _, protocol := range protocols {
		conn, err := net.Dial(protocol, addr)
		if err == nil {
			active = true
			conn.Close()
			break
		}
		if err != nil && strings.Contains(err.Error(), "too many open files") {
			time.Sleep(500 * time.Millisecond)
			protocols = append(protocols, protocol) //retry protocol again better to rate limit probably
		}
	}
	if active {
		fmt.Println(host, ":", port, " is active")
	}
	return true
}
