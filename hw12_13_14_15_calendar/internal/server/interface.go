package server

//go:generate mockgen -source=interface.go -destination=mocks/server_mock.go -package=server

// IServer Контракт сервера.
type IServer interface {
	Run() error
	Stop() error
}
