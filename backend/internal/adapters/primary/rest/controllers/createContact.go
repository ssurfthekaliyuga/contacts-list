package controllers

import (
	"contacts-list/internal/adapters/primary/rest"
	"contacts-list/internal/adapters/primary/rest/middlewares/auth"
	"contacts-list/internal/domain/ents"
	"context"
	"github.com/gofiber/fiber/v2"
)

type createContactBody struct {
	FullName    string `json:"fullName"`
	PhoneNumber string `json:"phoneNumber"`
	Note        string `json:"note"`
}

type contactCreator interface {
	Create(ctx context.Context, in ents.CreateContactIn) (*ents.Contact, error)
}

func NewCreateContact(creator contactCreator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body createContactBody
		if err := c.BodyParser(&body); err != nil {
			return rest.NewUnmarshalError(err)
		}

		in := ents.CreateContactIn{
			CreatorID:   auth.Extract(c),
			FullName:    body.FullName,
			PhoneNumber: body.PhoneNumber,
			Note:        body.Note,
		}

		contact, err := creator.Create(c.UserContext(), in)
		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"contact": contact,
		})
	}
}
