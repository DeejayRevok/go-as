-- +goose Up
-- +goose StatementBegin
CREATE TABLE roles (
    name TEXT PRIMARY KEY
);
CREATE TABLE role_permission (
    role_name TEXT NOT NULL,
    permission_name TEXT NOT NULL,
    FOREIGN KEY (role_name) REFERENCES roles(name), 
    FOREIGN KEY (permission_name) REFERENCES permissions(name),
    UNIQUE (role_name, permission_name)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE role_permission;
DROP TABLE roles;
-- +goose StatementEnd
