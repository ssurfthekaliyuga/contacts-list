package rest

import (
	"contacts-list/internal/domain/errs"
	"contacts-list/pkg/sl"
	"errors"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"net/http"
)

type errorHandler struct {
	logger *slog.Logger
}

func NewErrorHandler(logger *slog.Logger) fiber.ErrorHandler {
	handler := &errorHandler{
		logger: logger,
	}

	return handler.handleError
}

func (h *errorHandler) handleError(c *fiber.Ctx, inErr error) error {
	return h.convertError(c, inErr).send(c)
}

func (h *errorHandler) convertError(c *fiber.Ctx, err error) *appError {
	var (
		fiberErr *fiber.Error
		appErr   *errs.AppError
	)

	switch {
	case errors.As(err, &appErr):
		return h.convertAppError(c, appErr)
	case errors.As(err, &fiberErr):
		return h.convertFiberError(c, fiberErr)
	default:
		return h.convertInternalError(c, err)
	}
}

func (h *errorHandler) convertAppError(c *fiber.Ctx, appErr *errs.AppError) *appError {
	var status int
	switch appErr.Code {
	case errs.CodeNotFound:
		status = http.StatusNotFound
	case errs.CodeValidation:
		status = http.StatusBadRequest
	default:
		status = http.StatusInternalServerError
	}

	h.logger.Log(c.UserContext(), appErr.Level, "error",
		sl.Error(appErr),
		sl.Struct(*appErr),
	)

	return &appError{
		StatusCode: status,
		ID:         "", //todo sentry
		Message:    appErr.Message,
		Code:       string(appErr.Code),
		Additional: appErr.Additional,
	}
}

func (h *errorHandler) convertFiberError(c *fiber.Ctx, fiberErr *fiber.Error) *appError {
	var level slog.Level
	if fiberErr.Code >= 400 && fiberErr.Code < 500 {
		level = slog.LevelInfo
	} else {
		level = slog.LevelError
	}

	h.logger.Log(c.UserContext(), level, "error",
		sl.Error(fiberErr),
		sl.Struct(*fiberErr),
	)

	return &appError{
		StatusCode: fiberErr.Code,
		ID:         "", //todo sentry
		Message:    fiberErr.Message,
	}
}

func (h *errorHandler) convertInternalError(c *fiber.Ctx, err error) *appError {
	h.logger.ErrorContext(c.UserContext(), "internal server error",
		sl.Error(err),
	)

	return &appError{
		StatusCode: http.StatusInternalServerError,
		ID:         "", //todo
		Message:    "internal server error",
		Code:       "internal",
	}
}
