package errorx

type UserError struct {
	Hint string
}

func (e UserError) Error() string {
	return e.Hint
}

func NewUserError(hint string) UserError {
	return UserError{
		Hint: hint,
	}
}
