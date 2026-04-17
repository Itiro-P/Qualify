CREATE TABLE IF NOT EXISTS certification (
    id SERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "year" INTEGER NOT NULL, 
    description TEXT
);

CREATE TABLE IF NOT EXISTS analyst_certification (
    analyst_id INTEGER NOT NULL,
    certification_id INTEGER NOT NULL,
    PRIMARY KEY (analyst_id, certification_id),
    FOREIGN KEY (analyst_id) REFERENCES analyst (id) ON DELETE CASCADE,
    FOREIGN KEY (certification_id) REFERENCES certification (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS proposal_letter (
    id SERIAL PRIMARY KEY,
    client_id INTEGER NOT NULL,
    analyst_id INTEGER NOT NULL,
    proposed_hourly_rate FLOAT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    time_created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (client_id) REFERENCES client (id) ON DELETE CASCADE,
    FOREIGN KEY (analyst_id) REFERENCES analyst (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "service" (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    proposal_letter_id INTEGER NOT NULL,
    hourly_rate FLOAT NOT NULL,
    status VARCHAR(50) NOT NULL,
    time_created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (proposal_letter_id) REFERENCES proposal_letter (id) ON DELETE CASCADE
);
    
ALTER TABLE IF EXISTS review
    DROP COLUMN IF EXISTS analyst_id,
    DROP COLUMN IF EXISTS client_id,
    ADD COLUMN service_id INTEGER,
    ADD CONSTRAINT fk_service_id FOREIGN KEY (service_id) REFERENCES "service" (id) ON DELETE CASCADE;
