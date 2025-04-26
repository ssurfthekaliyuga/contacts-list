package postgres

import (
	"contacts-list/internal/domain/ents"
	"contacts-list/internal/domain/errs"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Contacts struct {
	db *pgxpool.Pool
}

func NewContacts(db *pgxpool.Pool) *Contacts {
	return &Contacts{
		db: db,
	}
}

func (r *Contacts) Get(ctx context.Context, in ents.GetContactsIn) ([]ents.Contact, error) { //todo pagination
	const query = `
		SELECT id, created_by, full_name, phone_number, note 
		FROM contacts
		WHERE created_by = $3
		LIMIT $1 OFFSET $2
	`

	limit := in.Size
	offset := in.Page * in.Size

	rows, err := r.db.Query(ctx, query, limit, offset, in.CreatedBy)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, collectContact)
}

func (r *Contacts) Create(ctx context.Context, in ents.CreateContactIn) (*ents.Contact, error) {
	const query = `
		INSERT INTO contacts(id, full_name, phone_number, note)
		VALUES($1, $2, $3)
		RETURNING id, created_by, full_name, phone_number, note
	`

	var contact ents.Contact

	err := r.db.
		QueryRow(ctx, query, uuid.New(), in.FullName, in.PhoneNumber, in.Note).
		Scan(&contact.ID, &contact.CreatedBy, &contact.FullName, &contact.PhoneNumber, &contact.Note)

	if err != nil {
		return nil, fmt.Errorf("cannot insert contact: %w", err)
	}

	return &contact, nil
}

func (r *Contacts) Update(ctx context.Context, in ents.UpdateContactIn) (*ents.Contact, error) {
	const query = `
		UPDATE contacts 
		SET full_name = COALESCE(NULLIF($2, ''), full_name), 
			phone_number = COALESCE(NULLIF($3, ''), phone_number), 
			note = COALESCE(NULLIF($4, ''), note)
		WHERE id = $1 
		RETURNING id, created_by, full_name, phone_number, note
	`

	var contact ents.Contact

	err := r.db.
		QueryRow(ctx, query, in.ContactID, in.FullName, in.PhoneNumber, in.Note).
		Scan(&contact.ID, &contact.CreatedBy, &contact.FullName, &contact.PhoneNumber, &contact.Note)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errs.NewNotFound(err, "updating contact was not found")
	}
	if err != nil {
		return nil, fmt.Errorf("cannot update contact: %w", err)
	}

	return &contact, nil
}

func (r *Contacts) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `
		DELETE FROM contacts WHERE id = $1
	`

	if _, err := r.db.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("cannot delete contact: %w", err)
	}

	return nil
}

func collectContact(row pgx.CollectableRow) (contact ents.Contact, err error) {
	err = row.Scan(
		&contact.ID,
		&contact.CreatedBy,
		&contact.FullName,
		&contact.PhoneNumber,
		&contact.Note,
	)
	return contact, err
}
