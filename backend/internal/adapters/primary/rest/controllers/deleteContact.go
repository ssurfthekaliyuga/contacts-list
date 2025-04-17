package controllers

import (
	"context"
	"github.com/gofiber/fiber/v2"
)

type contactDeleter interface {
	DeleteContact(context.Context, int64) error
}

func NewDeleteContact(deleter contactDeleter) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := struct {
			ID int64 `json:"ID"`
		}{}

		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrUnprocessableEntity
		}

		err := deleter.DeleteContact(c.Context(), req.ID)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
