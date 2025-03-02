package domain

import (
	"errors"
)

var ErrContactNotExist = errors.New("contact does not exist")

type Contact struct {
	ID          int64  `json:"ID"`
	FullName    string `json:"fullName"`
	PhoneNumber string `json:"phoneNumber"`
	Note        string `json:"note"`
}

type GetContactsIn struct {
	ForUpdate bool
	Page      int
	Size      int
}

type CreateContactIn struct {
	FullName    string `json:"fullName"`
	PhoneNumber string `json:"phoneNumber"`
	Note        string `json:"note"`
}

type UpdateContactIn struct {
	ID          int64  `json:"id"`
	FullName    string `json:"fullName"`
	PhoneNumber string `json:"phoneNumber"`
	Note        string `json:"note"`
}
