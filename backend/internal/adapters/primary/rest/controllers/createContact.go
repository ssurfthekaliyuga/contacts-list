package controllers

import (
	"contacts-list/internal/domain/ents"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type createContactBody struct {
	FullName    string `json:"fullName"`
	PhoneNumber string `json:"phoneNumber"`
	Note        string `json:"note"`
}

type contactCreator interface {
	CreateContact(context.Context, ents.CreateContactIn) (*ents.Contact, error)
}

func NewCreateContact(creator contactCreator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body createContactBody
		if err := c.BodyParser(&body); err != nil {
			return fiber.ErrUnprocessableEntity //todo
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
