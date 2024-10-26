package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	timeout := flag.Duration("timeout", time.Second*10, "timeout")
	flag.Parse()

	if flag.NArg() < 2 {
		log.Fatalln("Usage: go-telnet --timeout=10 <host> <port>")
	}

	host := flag.Arg(0)
	port := flag.Arg(1)
	address := host + ":" + port

	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		fmt.Printf("Cannot connect to %s: %v\n", address, err)
	}

	defer func() {
		if err := client.Close(); err != nil {
			fmt.Printf("Cannot close connection: %v\n", err)
		}
	}()

	log.Printf("...Connected to %s\n", address)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		defer cancel()
		select {
		case <-ctx.Done():
			return
		default:
			if err := client.Send(); err != nil {
				log.Fatalln(err)
			}
		}
	}()

	go func() {
		defer cancel()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := client.Receive(); err != nil {
					log.Fatalln(err)
				}
			}
		}
	}()

	<-ctx.Done()
}
