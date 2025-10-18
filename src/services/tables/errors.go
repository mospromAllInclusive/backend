package tables

import (
	"errors"
	"fmt"
)

type ErrorTableNotFound struct{}

func (e ErrorTableNotFound) Error() string {
	return "Table not found"
}

func IsErrTableNotFound(err error) bool {
	target := ErrorTableNotFound{}
	return errors.As(err, &target)
}

type ErrorColumnNotFound struct{}

func (e ErrorColumnNotFound) Error() string {
	return "Column not found"
}

func IsErrColumnNotFound(err error) bool {
	target := ErrorColumnNotFound{}
	return errors.As(err, &target)
}

type ErrorInvalidColumnValue struct {
	Value *string
}

func (e *ErrorInvalidColumnValue) Error() string {
	val := "null"
	if e.Value != nil {
		val = *e.Value
	}
	return fmt.Sprintf("Invalid column value `%s`", val)
}
