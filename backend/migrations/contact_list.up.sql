CREATE TABLE contacts
(
    id           UUID PRIMARY KEY,
    created_by   UUID, ---todo add index
    full_name    VARCHAR(255),
    phone_number VARCHAR(20),
    note         TEXT,
    created_at   timestamptz DEFAULT now() ---todo проверить нужен ли реально z при пуше одного времени в разных зонах в го
);
