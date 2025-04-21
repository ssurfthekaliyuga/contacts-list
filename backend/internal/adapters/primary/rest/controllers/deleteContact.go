package controllers

import (
	"contacts-list/internal/adapters/primary/rest"
	"contacts-list/internal/adapters/primary/rest/middlewares/auth"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type contactDeleter interface {
	Delete(ctx context.Context, userID uuid.UUID, contactID uuid.UUID) error
}

func NewDeleteContact(deleter contactDeleter) fiber.Handler {
	return func(c *fiber.Ctx) error {
		contactID, err := extractContactID(c)
		if err != nil {
			return rest.NewUnmarshalError(err) //todo param error
		}

		userID := auth.Extract(c)

		if err = deleter.Delete(c.UserContext(), userID, contactID); err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
