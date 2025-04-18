package errs

import (
	"log/slog"
)

type AppError struct {
	Underlying error
	Message    string
	Code       Code
	Public     bool
	Level      slog.Level
	Additional map[string]any
}

func (e *AppError) Error() string {
	return e.Message
}
