package goerrors

import "errors"

var (
	ErrPatientAlreadyExists = errors.New("patient already exists")
	ErrPatientNotFound      = errors.New("patient not found")
	ErrInvalidInput         = errors.New("invalid input")
	ErrInvalidPatientID     = errors.New("invalid patient ID")
	ErrDatabase             = errors.New("database error occurred")
)
