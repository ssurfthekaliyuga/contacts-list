package controllers

import (
	"contacts-list/internal/domain/ents"
	"contacts-list/internal/primary/rest"
	"contacts-list/internal/primary/rest/middlewares/auth"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type contactsUseCases interface {
	Get(ctx context.Context, in ents.GetContactsIn) ([]ents.Contact, error)
	Create(ctx context.Context, in ents.CreateContactIn) (*ents.Contact, error)
	Update(ctx context.Context, userID uuid.UUID, in ents.UpdateContactIn) (*ents.Contact, error)
	Delete(ctx context.Context, userID uuid.UUID, contactID uuid.UUID) error
}

type Contacts struct {
	contacts contactsUseCases
}

func NewContacts(contacts contactsUseCases) *Contacts {
	return &Contacts{
		contacts: contacts,
	}
}

func (o *Contacts) Get(c *fiber.Ctx) error {
	in := ents.GetContactsIn{
		CreatorID: auth.Extract(c),
		Page:      c.QueryInt("page", 0),
		Size:      c.QueryInt("size", 10),
	}

	contacts, err := o.contacts.Get(c.UserContext(), in)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"contacts": o.convertContacts(contacts),
	})
}

func (o *Contacts) Create(c *fiber.Ctx) error {
	var body createContactJSON
	if err := c.BodyParser(&body); err != nil {
		return rest.NewUnmarshalError(err)
	}

	in := ents.CreateContactIn{
		CreatorID:   auth.Extract(c),
		FullName:    body.FullName,
		PhoneNumber: body.PhoneNumber,
		Note:        body.Note,
	}

	contact, err := o.contacts.Create(c.UserContext(), in)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"contact": o.convertContact(*contact),
	})
}

func (o *Contacts) Update(c *fiber.Ctx) error {
	var body updateContactJSON
	if err := c.BodyParser(&body); err != nil {
		return rest.NewUnmarshalError(err)
	}

	in := ents.UpdateContactIn{
		ContactID:   o.contactID(c),
		FullName:    body.FullName,
		PhoneNumber: body.PhoneNumber,
		Note:        body.Note,
	}

	userID := auth.Extract(c)

	contact, err := o.contacts.Update(c.UserContext(), userID, in)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"contact": o.convertContact(*contact),
	})
}

func (o *Contacts) Delete(c *fiber.Ctx) error {
	contactID := o.contactID(c)
	userID := auth.Extract(c)

	err := o.contacts.Delete(c.UserContext(), userID, contactID)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (o *Contacts) contactID(c *fiber.Ctx) uuid.UUID {
	const paramName = "contactID"

	id, err := uuid.Parse(c.Params(paramName))
	if err != nil {
		return uuid.Nil
	}

	return id
}

func (o *Contacts) convertContact(in ents.Contact) contactJSON {
	return contactJSON{
		ID:          in.ID,
		CreatedBy:   in.CreatedBy,
		FullName:    in.FullName,
		PhoneNumber: in.PhoneNumber,
		Note:        in.Note,
	}
}

func (o *Contacts) convertContacts(in []ents.Contact) []contactJSON {
	return lo.Map(in, func(c ents.Contact, _ int) contactJSON {
		return o.convertContact(c)
	})
}
