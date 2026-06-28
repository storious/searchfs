package client

import "errors"

var (
	ErrNilBlockReader = errors.New("nil block reader")
	ErrNilBlockWriter = errors.New("nil block writer")
)
