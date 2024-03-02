ALTER TABLE
    users
ADD
    CONSTRAINT unique_username UNIQUE (username);

ALTER TABLE
    users
ADD
    CONSTRAINT check_role CHECK (role IN ('ADMIN', 'USER'));

ALTER TABLE
    urls
ADD
    CONSTRAINT url_format_check CHECK(
        original_url LIKE 'http://%'
        OR original_url LIKE 'https://%'
    );