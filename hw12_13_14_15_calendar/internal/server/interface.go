package server

// IServer Контракт сервера.
//
//go:generate go run github.com/vektra/mockery/v2@v2.49.0 --name=IServer
type IServer interface {
	Run() error
	Stop() error
}
