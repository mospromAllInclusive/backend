package users

import "errors"

type ErrorUserNotFound struct{}

func (e ErrorUserNotFound) Error() string {
	return "User not found"
}

func IsErrUserNotFound(err error) bool {
	target := ErrorUserNotFound{}
	return errors.As(err, &target)
}
