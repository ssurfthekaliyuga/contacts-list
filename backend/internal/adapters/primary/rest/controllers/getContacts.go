package controllers

import (
	"contacts-list/internal/adapters/primary/rest/middlewares/auth"
	"contacts-list/internal/domain/ents"
	"context"
	"github.com/gofiber/fiber/v2"
)

type contactProvider interface {
	Get(ctx context.Context, in ents.GetContactsIn) ([]ents.Contact, error)
}

func NewGetContacts(provider contactProvider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		in := ents.GetContactsIn{
			CreatorID: auth.Extract(c),
			Page:      c.QueryInt("page", 0),
			Size:      c.QueryInt("size", 10),
		}

		contacts, err := provider.Get(c.UserContext(), in)
		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"contacts": contacts,
		})
	}
}
