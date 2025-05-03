package controllers

import "github.com/google/uuid"

type createContactJSON struct {
	FullName    string `json:"fullName"`
	PhoneNumber string `json:"phoneNumber"`
	Note        string `json:"note"`
}

type updateContactJSON struct {
	FullName    string `json:"fullName"`
	PhoneNumber string `json:"phoneNumber"`
	Note        string `json:"note"`
}

type contactJSON struct {
	ID          uuid.UUID `json:"id"`
	CreatedBy   uuid.UUID `json:"createdBy"`
	FullName    string    `json:"fullName"`
	PhoneNumber string    `json:"phoneNumber"`
	Note        string    `json:"note"`
}
