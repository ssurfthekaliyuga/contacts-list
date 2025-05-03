package fiber

import (
	"github.com/gofiber/fiber/v2"
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
