package tp

import (
	"errors"
)

type DatabasEngine uint8
type Row = map[string]any

type (
	// TableName every table name must be regular text.
	TableName = string
	// ID every id value must be int64.
	ID = int64
)

var (
	ErrInvalidRowData   = errors.New("invalid row data")
	ErrRowAlreadyExists = errors.New("row with this ID already exists")
	ErrRowNotFound      = errors.New("row with this ID not found")
)
