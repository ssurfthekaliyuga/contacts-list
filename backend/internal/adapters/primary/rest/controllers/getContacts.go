package controllers

import (
	"contacts-list/internal/adapters/primary/rest/middlewares/auth"
	"contacts-list/internal/domain/ents"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type contactProvider interface {
	Get(ctx context.Context, userID uuid.UUID, in ents.GetContactsIn) ([]ents.Contact, error)
}

func NewGetContacts(provider contactProvider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		in := ents.GetContactsIn{
			Page: c.QueryInt("page", 0),
			Size: c.QueryInt("size", 10),
		}

		userID := auth.Extract(c)

		contacts, err := provider.Get(c.UserContext(), userID, in)
		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"contacts": contacts,
		})
	}
}
