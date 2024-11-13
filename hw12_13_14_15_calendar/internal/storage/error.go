package storage

import "errors"

var (
	ErrNotFound       = errors.New("event not found")
	ErrDateBusy       = errors.New("intersecting events")
	ErrNotImplemented = errors.New("method not implemented")
)
