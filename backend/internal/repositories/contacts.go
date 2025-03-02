package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"os-lab-3-1/internal/domain"
)

type ContactsRepository struct {
	db *pgxpool.Pool
}

func NewContactsRepository(db *pgxpool.Pool) *ContactsRepository {
	return &ContactsRepository{
		db: db,
	}
}

func (r *ContactsRepository) GetContacts(ctx context.Context, in domain.GetContactsIn) ([]domain.Contact, error) {
	query := `
		SELECT id, full_name, phone_number, note 
		FROM contacts LIMIT $1 OFFSET $2
	`

	limit := in.Size
	offset := in.Page * in.Size

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	collect := func(row pgx.CollectableRow) (domain.Contact, error) {
		contact := domain.Contact{}
		err := row.Scan(
			&contact.ID,
			&contact.FullName,
			&contact.PhoneNumber,
			&contact.Note,
		)
		return contact, err
	}

	contacts, err := pgx.CollectRows(rows, collect)
	if err != nil {
		return nil, fmt.Errorf("cannot collect row: %w", err)
	}

	return contacts, nil
}

func (r *ContactsRepository) CreateContact(ctx context.Context, in domain.CreateContactIn) (*domain.Contact, error) {
	const query = `
		INSERT INTO contacts(full_name, phone_number, note)
		VALUES($1, $2, $3)
		RETURNING id, full_name, phone_number, note
	`

	var contact domain.Contact

	err := r.db.
		QueryRow(ctx, query, in.FullName, in.PhoneNumber, in.Note).
		Scan(&contact.ID, &contact.FullName, &contact.PhoneNumber, &contact.Note)

	if err != nil {
		return nil, fmt.Errorf("cannot exec query: %w", err)
	}

	return &contact, nil
}

func (r *ContactsRepository) UpdateContact(ctx context.Context, in domain.UpdateContactIn) (*domain.Contact, error) {
	const query = `
		UPDATE contacts 
		SET full_name = COALESCE(NULLIF($2, ''), full_name), 
			phone_number = COALESCE(NULLIF($3, ''), phone_number), 
			note = COALESCE(NULLIF($4, ''), note)
		WHERE id = $1 
		RETURNING id, full_name, phone_number, note
	`

	var contact domain.Contact

	err := r.db.
		QueryRow(ctx, query, in.ID, in.FullName, in.PhoneNumber, in.Note).
		Scan(&contact.ID, &contact.FullName, &contact.PhoneNumber, &contact.Note)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrContactNotExist
	}
	if err != nil {
		return nil, fmt.Errorf("cannot update contact: %w", err)
	}

	return &contact, nil
}

func (r *ContactsRepository) DeleteContact(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM contacts WHERE id = $1
	`

	if _, err := r.db.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("cannot exec query: %w", err)
	}

	return nil
}
