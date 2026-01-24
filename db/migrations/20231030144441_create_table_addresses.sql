-- +goose Up
CREATE TABLE addresses (
    id VARCHAR(36) PRIMARY KEY,
    contact_id VARCHAR(36) NOT NULL,
    street VARCHAR(255),
    city VARCHAR(255),
    province VARCHAR(255),
    postal_code VARCHAR(10),
    country VARCHAR(100),
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    CONSTRAINT fk_addresses_contact_id FOREIGN KEY (contact_id) REFERENCES contacts (id)
);

-- +goose Down
DROP TABLE addresses;
