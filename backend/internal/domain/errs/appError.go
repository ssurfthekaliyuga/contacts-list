package errs

import (
	"log/slog"
)

type AppError struct {
	Underlying error
	Message    string
	Code       Code
	Level      slog.Level
	Additional map[string]any
}

func (e *AppError) Error() string { //TODO
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Underlying
}
