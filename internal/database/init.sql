\c database;

CREATE TABLE refresh_tokens (
    rawRefreshToken TEXT UNIQUE PRIMARY KEY,
    token_hash TEXT NOT NULL,                 
    user_guid UUID NOT NULL UNIQUE,                 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);