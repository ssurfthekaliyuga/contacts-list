package rest

import (
	"contacts-list/internal/domain/errs"
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

type appError struct {
	StatusCode int            `json:"-"`
	ID         string         `json:"id,omitempty"`
	Message    string         `json:"message,omitempty"`
	Code       string         `json:"code,omitempty"`
	Additional map[string]any `json:"additional,omitempty"`
}

func (e *appError) send(c *fiber.Ctx) error {
	return c.Status(e.StatusCode).JSON(e)
}

func NewUnmarshalError(err error) *errs.AppError { //todo
	return &errs.AppError{
		Underlying: err,
		Message:    "cannot unmarshal request",
		Code:       errs.CodeValidation,
		Level:      slog.LevelInfo,
		Additional: nil,
	}
}
