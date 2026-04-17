CREATE TABLE IF NOT EXISTS user_profile (
    user_id INTEGER PRIMARY KEY REFERENCES "user"(id) ON DELETE CASCADE,
    biography TEXT
);

CREATE TABLE IF NOT EXISTS analyst_profile (
    analyst_id INTEGER PRIMARY KEY REFERENCES user_profile(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS client_profile (
    client_id INTEGER PRIMARY KEY REFERENCES user_profile(user_id) ON DELETE CASCADE
);
