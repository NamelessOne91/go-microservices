-- Create the users table
CREATE TABLE users (
    id        SERIAL PRIMARY KEY,
    email     VARCHAR(255) NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    user_active    INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Insert test users
INSERT INTO users (email, first_name, last_name, password, user_active, created_at, updated_at)
VALUES ('test1@example.com', 'John', 'Doe', 'password1', 1, NOW(), NOW()),
       ('test2@example.com', 'Jane', 'Smith', 'password2', 1, NOW(), NOW()),
       ('test3@example.com', 'Michael', 'Johnson', 'password3', 1, NOW(), NOW());