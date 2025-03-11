package errors

type AuthError struct {
	Msg string
}

func (ae AuthError) Error() string {
	return ae.Msg
}
