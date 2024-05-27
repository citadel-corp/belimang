DROP TYPE IF EXISTS user_type;
CREATE TYPE user_type AS ENUM('Admin', 'User');

CREATE TABLE IF NOT EXISTS
users (
	id SERIAL PRIMARY KEY,
	uid CHAR(16) UNIQUE NOT NULL,
    username VARCHAR(30) NOT NULL UNIQUE,
    user_type user_type NOT NULL,
    hashed_password BYTEA NOT NULL,
    email VARCHAR NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT current_timestamp
);

CREATE INDEX IF NOT EXISTS users_user_type
	ON users USING HASH (user_type);
CREATE INDEX IF NOT EXISTS users_username 
	ON users USING HASH (lower(username));
CREATE INDEX IF NOT EXISTS users_email 
	ON users USING HASH (lower(email));
CREATE INDEX IF NOT EXISTS users_created_at_desc
	ON users(created_at DESC);
CREATE INDEX IF NOT EXISTS users_created_at_asc
	ON users(created_at ASC);
