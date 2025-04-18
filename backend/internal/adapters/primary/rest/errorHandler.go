package rest

import (
	"contacts-list/internal/domain/errs"
	"errors"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"net/http"
)

type appError struct {
	ID         string         `json:"id,omitempty"`
	Message    string         `json:"message,omitempty"`
	Code       string         `json:"code,omitempty"`
	Additional map[string]any `json:"additional,omitempty"`
}

func (e *appError) send(c *fiber.Ctx) error {
	return c.Status(e.statusCode()).JSON(e)
}

func (e *appError) statusCode() int {
	switch e.Code {
	case errs.CodeInternal:
		return http.StatusInternalServerError
	case errs.CodeNotFound:
		return http.StatusNotFound
	case errs.CodeValidation:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

type ErrorHandler struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

func (h *ErrorHandler) Handle(c *fiber.Ctx, inError error) error {
	var appErr *errs.AppError
	if errors.As(inError, &appErr) {
		appErr = errs.NewInternal(inError)
	}

	h.logger.Log(c.UserContext(), appErr.Level, "could not process request",
		slog.Any("error", appErr),
	)

	outErr := appError{
		ID:         "", //todo
		Message:    appErr.Message,
		Code:       string(appErr.Code),
		Additional: appErr.Additional,
	}

	return outErr.send(c)
}
