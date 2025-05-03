package errs

import (
	"log/slog"
)

//todo maybe generate sentry error in these functions

func NewNotFound(err error, msg string) *AppError {
	return &AppError{
		Underlying: err,
		Code:       CodeNotFound,
		Message:    msg,
		Level:      slog.LevelInfo,
	}
}

func NewUnauthorized(msg string) *AppError {
	return &AppError{
		Underlying: nil,
		Code:       CodeUnauthorized,
		Message:    msg,
		Level:      slog.LevelInfo,
	}
}
