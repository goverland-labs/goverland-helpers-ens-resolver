package infura

import (
	"errors"
)

var (
	ErrUnknownAddress = errors.New("address not found")
	ErrUnknownDomain  = errors.New("domain not found")
)
