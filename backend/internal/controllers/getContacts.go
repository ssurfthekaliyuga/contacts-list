package controllers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"os-lab-3-1/internal/domain"
)

type contactProvider interface {
	GetContacts(context.Context, domain.GetContactsIn) ([]domain.Contact, error)
}

func NewGetContacts(provider contactProvider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		in := domain.GetContactsIn{
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
