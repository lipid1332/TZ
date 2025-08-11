package storage

import (
	"errors"
)

var (
	ErrWalletNotFound = errors.New("this wallet id not found")
	ErrNegativeAmount = errors.New("amount cant be negative")
)
