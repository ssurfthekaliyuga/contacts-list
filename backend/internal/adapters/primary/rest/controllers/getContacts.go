package controllers

import (
	"contacts-list/internal/domain/ents"
	"context"
	"github.com/gofiber/fiber/v2"
)

type contactProvider interface {
	GetContacts(context.Context, ents.GetContactsIn) ([]ents.Contact, error)
}

func NewGetContacts(provider contactProvider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		in := ents.GetContactsIn{
			Page: c.QueryInt("page", 0),
			Size: c.QueryInt("size", 10),
		}

		contacts, err := provider.GetContacts(c.Context(), in)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(fiber.Map{
			"contacts": contacts,
		})
	}
}
