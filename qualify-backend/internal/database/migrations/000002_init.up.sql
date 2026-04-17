ALTER TABLE IF EXISTS "user"
    ADD COLUMN IF NOT EXISTS time_created TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    ADD COLUMN IF NOT EXISTS country_code CHAR(2)       NOT NULL DEFAULT 'BR',
    ADD COLUMN IF NOT EXISTS country_name VARCHAR(100)  NOT NULL DEFAULT 'Brazil',
    ADD COLUMN IF NOT EXISTS city         VARCHAR(100)  NOT NULL DEFAULT 'Campo Mourão',
    ADD COLUMN IF NOT EXISTS timezone     VARCHAR(50)  NOT NULL DEFAULT 'America/Sao_Paulo';


CREATE TABLE IF NOT EXISTS analyst (
    user_id INTEGER NOT NULL UNIQUE,
    hourly_rate FLOAT NOT NULL DEFAULT 0 CHECK (hourly_rate >= 0),
    total_reviews INTEGER NOT NULL DEFAULT 0 CHECK (total_reviews >= 0),
    mean_rating INTEGER NOT NULL DEFAULT 0 CHECK (mean_rating >= 0 AND mean_rating <= 5),
    FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS client (
    user_id INTEGER NOT NULL UNIQUE,
    proposed_budget FLOAT NOT NULL DEFAULT 0 CHECK (proposed_budget >= 0),
    FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS analyst_skill (
    analyst_id INTEGER NOT NULL,
    skill_id INTEGER NOT NULL,
    PRIMARY KEY (analyst_id, skill_id),
    FOREIGN KEY (analyst_id) REFERENCES analyst (user_id) ON DELETE CASCADE,
    FOREIGN KEY (skill_id) REFERENCES skill (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS review (
    id SERIAL PRIMARY KEY,
    analyst_id INTEGER NOT NULL,
    client_id INTEGER NOT NULL,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    time_created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (analyst_id) REFERENCES analyst (user_id) ON DELETE CASCADE,
    FOREIGN KEY (client_id) REFERENCES client (user_id) ON DELETE CASCADE
);
