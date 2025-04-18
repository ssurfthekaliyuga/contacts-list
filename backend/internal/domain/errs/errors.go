package errs

import (
	"log/slog"
)

func NewInternal(err error) *AppError {
	return &AppError{
		Underlying: err,
		Code:       CodeInternal,
		Message:    "internal server error",
		Public:     false,
		Level:      slog.LevelError,
	}
}

func NewNotFound(err error, msg string) *AppError {
	return &AppError{
		Underlying: err,
		Code:       CodeNotFound,
		Message:    msg,
		Public:     true,
		Level:      slog.LevelInfo,
	}
}
