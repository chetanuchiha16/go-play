-- +goose Up
CREATE TABLE IF NOT EXISTS permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

-- Pivot table for users and roles
CREATE TABLE IF NOT EXISTS user_role (
    user_id INT REFERENCES users(id),
    perm_id INT REFERENCES permissions(id),
    PRIMARY KEY (user_id, perm_id)
);

-- Pivot table for roles and permissions
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id INT REFERENCES roles(id),
    perm_id INT REFERENCES permissions(id),
    PRIMARY KEY (role_id, perm_id)
);

-- +goose Down
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS user_role;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS permissions;
