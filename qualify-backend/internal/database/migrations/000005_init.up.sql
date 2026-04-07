CREATE TABLE IF NOT EXISTS user_skill (
    user_id  int REFERENCES "user"(id) ON DELETE CASCADE,
    skill_id int REFERENCES skill(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, skill_id)
);