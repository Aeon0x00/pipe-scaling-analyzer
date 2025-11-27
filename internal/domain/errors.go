package domain

import "errors"

var (
	ErrPipeNotFound = errors.New("pipe not found")
	ErrAlertNotFound = errors.New("alert not found")
	
	ErrEmptyPoints = errors.New("sensor points list is empty")
	ErrInvalidPipeData = errors.New("invalid pipe data")
	
)
