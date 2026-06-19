CREATE TABLE IF NOT EXISTS user_profile (
    user_id INTEGER PRIMARY KEY REFERENCES "user"(id) ON DELETE CASCADE,
    biography TEXT
);