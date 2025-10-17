package tables

import "errors"

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
