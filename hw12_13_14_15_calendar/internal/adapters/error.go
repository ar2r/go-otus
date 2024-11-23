package adapters

import "errors"

var (
	ErrNotFound = errors.New("event not found")
	ErrDateBusy = errors.New("intersecting events")
)
