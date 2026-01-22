CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE posts (
    post_id uuid DEFAULT uuid_generate_v4 (),
    title varchar NOT NULL,
    content text NOT NULL,
    PRIMARY KEY (post_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
