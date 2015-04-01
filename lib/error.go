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
