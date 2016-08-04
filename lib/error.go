package elastigo

import (
	"errors"
)

var (
	// 404 Response.
	RecordNotFound = errors.New("record not found")
	// DocAlreadyExists - document exists during insert doc
	DocAlreadyExists = errors.New("document already exists")
)

const (
	// error code for underlying connection error
	connErrorCode int = -777
)

// IsRecordNotFound checks if the error of http 404 error
func IsRecordNotFound(err error) bool {
	if esErr, ok := err.(*ESError); ok {
		if esErr.Code == 404 {
			return true
		}
	}
	return false
}

// IsConnError checks if underlying connection error
func IsConnError(err error) bool {
	if esErr, ok := err.(*ESError); ok {
		if esErr.Code == connErrorCode {
			return true
		}
	}
	return false
}
