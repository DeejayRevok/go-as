-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    email text PRIMARY KEY,
    password text,
    superuser BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE TABLE user_role (
    user_email VARCHAR(36) NOT NULL,
    role_name TEXT NOT NULL,
    FOREIGN KEY (user_email) REFERENCES users(email),
    FOREIGN KEY (role_name) REFERENCES roles(name), 
    UNIQUE (user_email, role_name)
);
CREATE TABLE user_permission (
    user_email VARCHAR(36) NOT NULL,
    permission_name text NOT NULL,
    FOREIGN KEY (user_email) REFERENCES users(email),
    FOREIGN KEY (permission_name) REFERENCES permissions(name), 
    UNIQUE (user_email, permission_name)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_role;
DROP TABLE user_permission;
DROP TABLE users;
-- +goose StatementEnd
