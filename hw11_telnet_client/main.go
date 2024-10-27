package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
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

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	client := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "...failed to connect to %s\n", address)
		return
	}
	fmt.Fprintf(os.Stderr, "...connected to %s\n", address)

	defer func() {
		if err := client.Close(); err != nil {
			fmt.Printf("Cannot close connection: %v\n", err)
		}
	}()

	go func() {
		defer cancel()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				err := client.Send()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Cannot send data: %v\n", err)
					return
				}
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
				err := client.Receive()
				if err != nil {
					switch {
					case errors.Is(err, syscall.ECONNRESET):
						fmt.Fprint(os.Stderr, "...Connection was reset by peer\n")
					case errors.Is(err, io.EOF):
						fmt.Fprint(os.Stderr, "...EOF\n")
					default:
						fmt.Fprintf(os.Stderr, "Cannot receive data: %v\n", err)
					}
					return
				}
			}
		}
	}()

	<-ctx.Done()
}
