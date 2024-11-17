package adapters

import "errors"

var (
	ErrNotFound       = errors.New("event not found")
	ErrDateBusy       = errors.New("intersecting events")
	ErrNotImplemented = errors.New("method not implemented")
	ErrNoUserID       = errors.New("no user id")
)
