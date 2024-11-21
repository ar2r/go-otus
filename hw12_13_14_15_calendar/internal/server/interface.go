package server

// IServer Контракт сервера.
type IServer interface {
	Run() error
	Stop() error
}
