package commands

type AppError struct {
	Error error
}

func ThrowError(err error) {
	panic(&AppError{err})
}
