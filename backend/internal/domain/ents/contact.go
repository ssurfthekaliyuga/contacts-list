package ents

import (
	"github.com/google/uuid"
)

type Contact struct {
	ID          uuid.UUID
	CreatedBy   uuid.UUID
	FullName    string
	PhoneNumber string
	Note        string
}

type GetContactsIn struct {
	CreatedBy uuid.UUID
	Page      int
	Size      int
	ForUpdate bool
}

type CreateContactIn struct {
	CreatorID   uuid.UUID
	FullName    string
	PhoneNumber string
	Note        string
}

type UpdateContactIn struct {
	ID          uuid.UUID
	FullName    string
	PhoneNumber string
	Note        string
}
