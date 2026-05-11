-- +goose Up
CREATE TABLE permissions (
    id SERIAL INT PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE roles (
    id SERIAL INT,
    name VARCHAR(50) UNIQUE NOT NULL
);

-- Pivot table for users and roles
CREATE TABLE user_role (
    users_id INT REFERENCES users(id),
    perm_id INT REFERENCES permissions(id),
    PRIMARY KEY (user_id, perm_id)
);

-- Pivot table for roles and permissions
CREATE TABLE role_permissions (
    role_id INT REFERENCES roles(id),
    perm_id INT REFERENCES permissions(id),
    PRIMARY KEY (role_id, perm_id)
);

-- +goose Down
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS user_role;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS permissions;
