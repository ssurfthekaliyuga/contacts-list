package controllers

import (
	"contacts-list/internal/adapters/primary/rest"
	"contacts-list/internal/adapters/primary/rest/middlewares/auth"
	"contacts-list/internal/domain/ents"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type updateContactIn struct {
	FullName    string `json:"fullName"`
	PhoneNumber string `json:"phoneNumber"`
	Note        string `json:"note"`
}

type contactUpdater interface {
	Update(ctx context.Context, userID uuid.UUID, in ents.UpdateContactIn) (*ents.Contact, error)
}

func NewUpdateContact(updater contactUpdater) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body updateContactIn
		if err := c.BodyParser(&body); err != nil {
			return rest.NewUnmarshalError(err)
		}

		contactID, err := extractContactID(c)
		if err != nil {
			return rest.NewUnmarshalError(err) //todo param error
		}

		in := ents.UpdateContactIn{
			ContactID:   contactID,
			FullName:    body.FullName,
			PhoneNumber: body.PhoneNumber,
			Note:        body.Note,
		}

		userID := auth.Extract(c)

		contact, err := updater.Update(c.UserContext(), userID, in)
		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"contact": contact,
		})
	}
}
