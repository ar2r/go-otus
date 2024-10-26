package main

import (
	"errors"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address    string
	timeout    time.Duration
	in         io.ReadCloser
	out        io.Writer
	connection net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (client *telnetClient) Connect() (err error) {
	client.connection, err = net.DialTimeout("tcp", client.address, client.timeout)
	return
}

func (client *telnetClient) Close() (err error) {
	return client.connection.Close()
}

func (client *telnetClient) Send() (err error) {
	buffer := make([]byte, 1024)
	n, err := client.in.Read(buffer)
	if errors.Is(err, io.EOF) {
		return
	}
	if err != nil {
		return
	}
	_, err = client.connection.Write(buffer[:n])
	return
}

func (client *telnetClient) Receive() (err error) {
	buffer := make([]byte, 1024)
	n, err := client.connection.Read(buffer)
	if err != nil {
		return
	}
	_, err = client.out.Write(buffer[:n])
	return
}
