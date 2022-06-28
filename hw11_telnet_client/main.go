package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	host := args[0]
	port := args[1]
	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Fatalf("Connection error: %v", err)
	}
	defer func(client TelnetHandler) {
		err := client.Close()
		if err != nil {

		}
	}(client)

	fmt.Fprintf(os.Stderr, "...Connected to %s\n", address)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		err := client.Receive()
		if err != nil {
			log.Fatalf("Receive error: %v", err)
		}
		fmt.Fprintf(os.Stderr, "...Connection was closed by peer\n")
	}()

	go func() {
		defer wg.Done()
		err := client.Send()
		if err != nil {
			log.Fatalf("Send error: %v", err)
		}
		fmt.Fprintf(os.Stderr, "...EOF\n")
	}()

	wg.Wait()
}
