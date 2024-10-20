CREATE TABLE IF NOT EXISTS users (
    user_id VARCHAR(60) PRIMARY KEY,
    username VARCHAR(60) NOT NULL,
    first_name VARCHAR(60),
    last_name VARCHAR(60),
    role VARCHAR(30),
    email VARCHAR(60) NOT NULL UNIQUE, -- Adding UNIQUE constraint for emails
    password VARCHAR(150) NOT NULL,
    created_at BIGINT DEFAULT 0, -- Use BIGINT instead of INT UNSIGNED
    updated_at BIGINT DEFAULT 0 -- Use BIGINT instead of INT UNSIGNED
);
