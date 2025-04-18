package domain

import (
	"contacts-list/internal/domain/ents"
	"context"
	"github.com/google/uuid"
	"log/slog"
)

type contactsRepo interface {
	Get(ctx context.Context, in ents.GetContactsIn) ([]ents.Contact, error)
	Create(ctx context.Context, in ents.CreateContactIn) (*ents.Contact, error)
	Update(ctx context.Context, in ents.UpdateContactIn) (*ents.Contact, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Contacts struct {
	repo   contactsRepo
	logger *slog.Logger
}

func NewContacts(repo contactsRepo, logger *slog.Logger) *Contacts {
	return &Contacts{
		repo:   repo,
		logger: logger,
	}
}

func (c *Contacts) Get(ctx context.Context, userID uuid.UUID, in ents.GetContactsIn) ([]ents.Contact, error) { //todo проверять кто запрашивает
	return c.repo.Get(ctx, in)
}

func (c *Contacts) Create(ctx context.Context, in ents.CreateContactIn) (*ents.Contact, error) { //todo валидация
	return c.repo.Create(ctx, in)
}

func (c *Contacts) Update(ctx context.Context, userID uuid.UUID, in ents.UpdateContactIn) (*ents.Contact, error) { //todo проверять кто обновляет + валидация
	return c.repo.Update(ctx, in)
}

func (c *Contacts) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error { //todo проверять кто удаляет
	return c.repo.Delete(ctx, id)
}
