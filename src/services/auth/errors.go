package auth

type ErrorWrongPassword struct{}

func (e ErrorWrongPassword) Error() string {
	return "Wrong password"
}
