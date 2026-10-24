CREATE TYPE user_role AS ENUM ('admin', 'user');

CREATE SEQUENCE IF NOT EXISTS users_id_seq
    START WITH 10000000
    INCREMENT BY 1
    MINVALUE 10000000
    CACHE 1;

CREATE TABLE users (
    id BIGINT PRIMARY KEY
       DEFAULT nextval('users_id_seq'),

    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,

    role user_role NOT NULL DEFAULT 'user',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER SEQUENCE users_id_seq
    OWNED BY users.id;
