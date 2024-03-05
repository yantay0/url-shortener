CREATE TABLE url (
    id bigserial PRIMARY KEY,
    original_url text NOT NULL,
    short_url text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
    -- user_id BIGINT NOT NULL,
    -- FOREIGN KEY (user_id) REFERENCES users(id)
);


-- CREATE TABLE IF NOT EXISTS users (
--     id bigserial PRIMARY KEY,
--     created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
--     name text NOT NULL,
--     email citext UNIQUE NOT NULL,
--     password_hash bytea NOT NULL,
--     activated bool NOT NULL,
--     version integer NOT NULL DEFAULT 1
-- );