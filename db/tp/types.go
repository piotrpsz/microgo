package tp

import (
	"errors"
)

// Database
type (
	// DatabaseEngine custom type for database description.
	DatabaseEngine uint8
	// Row custom type identical with map
	Row = map[string]any
	// TableName every table name must be regular text.
	TableName = string
	// ID every id value must be int64.
	ID = int64
)

var (
	// ErrInvalidRowData when something is wrong in data
	ErrInvalidRowData = errors.New("invalid row data")
	// ErrRowNotFound not found requested row
	ErrRowNotFound = errors.New("row with this ID not found")
)
