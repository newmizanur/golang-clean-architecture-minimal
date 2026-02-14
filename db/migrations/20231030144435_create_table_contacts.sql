-- +goose Up
CREATE TABLE contacts (
    id VARCHAR(36) PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100),
    email VARCHAR(100),
    phone VARCHAR(100),
    user_id VARCHAR(36) NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    CONSTRAINT fk_contacts_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);

-- +goose Down
DROP TABLE contacts;
