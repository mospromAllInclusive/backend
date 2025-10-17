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
