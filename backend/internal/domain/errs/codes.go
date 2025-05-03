package errs

type Code string

const (
	CodeValidation   = "validation"
	CodeBadInput     = "bad_input"
	CodeNotFound     = "not_found"
	CodeUnauthorized = "unauthorized"
	CodeInternal     = "internal"
)
