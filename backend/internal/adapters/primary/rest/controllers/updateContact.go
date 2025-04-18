package controllers

import (
	"contacts-list/internal/domain/ents"
	"context"
	"github.com/gofiber/fiber/v2"
)

type contactUpdater interface {
	UpdateContact(context.Context, ents.UpdateContactIn) (*ents.Contact, error)
}

func NewUpdateContact(updater contactUpdater) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in ents.UpdateContactIn

		if err := c.BodyParser(&in); err != nil {
			return fiber.ErrUnprocessableEntity
		}

		contact, err := updater.UpdateContact(c.Context(), in)
		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"contact": contact,
		})
	}
}
