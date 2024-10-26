package main

import (
	"context"
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
	ctx        context.Context
	cancel     context.CancelFunc
	address    string
	timeout    time.Duration
	in         io.ReadCloser
	out        io.Writer
	connection net.Conn
}

func NewTelnetClient(
	ctx context.Context, address string, timeout time.Duration, in io.ReadCloser, out io.Writer,
) TelnetClient {
	ctx, cancel := context.WithCancel(ctx)
	return &telnetClient{
		ctx:     ctx,
		cancel:  cancel,
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (client *telnetClient) Connect() (err error) {
	ctx, cancel := context.WithTimeout(client.ctx, client.timeout)
	defer cancel()

	client.connection, err = net.DialTimeout("tcp", client.address, client.timeout)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		if err := client.Close(); err != nil {
			return err
		}
		return ctx.Err()
	default:
		return nil
	}
}

func (client *telnetClient) Close() (err error) {
	if client.connection == nil {
		return errors.New("client not connected")
	}
	if err := client.connection.Close(); err != nil {
		return errors.New("connection close failed")
	}
	return nil
}

func (client *telnetClient) Send() (err error) {
	buffer := make([]byte, 1024)
	for {
		select {
		case <-client.ctx.Done():
			return client.ctx.Err()
		default:
			n, err := client.in.Read(buffer)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}
				return err
			}
			_, err = client.connection.Write(buffer[:n])
			if err != nil {
				return err
			}
		}
	}
}

func (client *telnetClient) Receive() (err error) {
	buffer := make([]byte, 1024)
	for {
		select {
		case <-client.ctx.Done():
			return client.ctx.Err()
		default:
			n, err := client.connection.Read(buffer)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}
				return err
			}
			_, err = client.out.Write(buffer[:n])
			if err != nil {
				return err
			}
		}
	}
}
