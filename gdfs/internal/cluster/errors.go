package cluster

import "errors"

var (
	ErrEmptyDataNodeID   = errors.New("empty datanode id")
	ErrEmptyDataNodeAddr = errors.New("empty datanode address")
)
