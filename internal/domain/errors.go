package domain

import "errors"

var (
	ValidationError  = errors.New("validate error")
	ErrorNoAuth      = errors.New("no such user")
	ErrorWhileFetch  = errors.New("error while fetch")
	ErrorWhileCreate = errors.New("error while create")
)
