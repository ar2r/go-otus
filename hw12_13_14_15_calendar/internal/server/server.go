package server

// Interface Контракт сервера.
type Interface interface {
	Run() error
	Stop() error
}
