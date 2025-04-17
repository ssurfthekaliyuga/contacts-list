package controllers

import (
	"contacts-list/internal/domain"
	"context"
	"github.com/gofiber/fiber/v2"
)

type contactCreator interface {
	CreateContact(context.Context, domain.CreateContactIn) (*domain.Contact, error)
}

func NewCreateContact(creator contactCreator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in domain.CreateContactIn

		if err := c.BodyParser(&in); err != nil {
			return fiber.ErrUnprocessableEntity
		}

		contact, err := creator.CreateContact(c.Context(), in)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(fiber.Map{
			"contact": contact,
		})
	}
}
