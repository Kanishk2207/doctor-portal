CREATE TABLE IF NOT EXISTS patients (
    patient_id VARCHAR(60) PRIMARY KEY,
    first_name VARCHAR(60),
    last_name VARCHAR(60),
    email VARCHAR(60) NOT NULL UNIQUE, -- Adding UNIQUE constraint for emails
    created_at BIGINT DEFAULT 0, -- Use BIGINT instead of INT UNSIGNED
    updated_at BIGINT DEFAULT 0 -- Use BIGINT instead of INT UNSIGNED
);
