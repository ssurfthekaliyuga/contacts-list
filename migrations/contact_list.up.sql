CREATE TABLE contacts
(
    id           BIGSERIAL PRIMARY KEY,
    full_name    VARCHAR(255),
    phone_number VARCHAR(20),
    note         TEXT,

    CONSTRAINT phone_format CHECK (
        phone_number ~ '^\+?\d{1,4}[-\s]?\(?\d{1,3}\)?[-\s]?\d{1,4}[-\s]?\d{1,4}$'
)
    );
