echo 'DB_NAME="testdb"

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

echo "Tables created successfully."' >> init.sh
