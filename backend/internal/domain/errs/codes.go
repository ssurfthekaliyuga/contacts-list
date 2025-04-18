package errs

type Code string

const (
	CodeInternal   = "internal"
	CodeValidation = "validation"
	CodeNotFound   = "not_found"
)
