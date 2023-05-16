#!/bin/bash

# Set the database name
DB_NAME="testdb"

# Delete the existing tables
psql -U postgres -d ${DB_NAME} << EOF
DROP TABLE IF EXISTS wallet;
DROP TABLE IF EXISTS users;
EOF

# Recreate the tables
psql -U postgres -d ${DB_NAME} << EOF
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL
);

CREATE TABLE wallet (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    amount BIGINT NOT NULL,
    currency TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);
EOF

echo "Tables created successfully."
