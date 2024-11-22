package server

//go:generate mockgen -source=interface.go -destination=mocks/interface.go

// IServer Контракт сервера.
type IServer interface {
	Run() error
	Stop() error
}
