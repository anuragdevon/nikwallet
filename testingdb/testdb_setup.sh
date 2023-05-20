#!/bin/bash
#-- Create the custom domain for currency
# CREATE DOMAIN currency_enum AS TEXT CHECK (VALUE IN ('USD', 'EUR', 'INR'));

DB_NAME="testdb"

psql -U postgres -d ${DB_NAME} << EOF
DROP TABLE IF EXISTS wallets;
DROP TABLE IF EXISTS users;
EOF

psql -U postgres -d ${DB_NAME} << EOF

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email_id TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  deleted_at TIMESTAMP DEFAULT now()
);

CREATE TABLE wallets (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    amount TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

EOF

echo "Tables created successfully."

    # amount NUMERIC(15, 2) NOT NULL,
    # currency TEXT NOT NULL,
