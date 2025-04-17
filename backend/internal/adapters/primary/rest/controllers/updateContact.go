package controllers

import (
	"contacts-list/internal/domain"
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
)

type contactUpdater interface {
	UpdateContact(context.Context, domain.UpdateContactIn) (*domain.Contact, error)
}

func NewUpdateContact(updater contactUpdater) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in domain.UpdateContactIn

		if err := c.BodyParser(&in); err != nil {
			return fiber.ErrUnprocessableEntity
		}

		contact, err := updater.UpdateContact(c.Context(), in)
		if errors.Is(err, domain.ErrContactNotExist) {
			return fiber.ErrNotFound
		}
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(fiber.Map{
			"contact": contact,
		})
	}
}
