CREATE TABLE IF NOT EXISTS conversation (
    id SERIAL PRIMARY KEY,

    service_id INTEGER NULL,
    proposal_id INTEGER NULL,

    client_id INTEGER NOT NULL REFERENCES client(id),
    analyst_id INTEGER NOT NULL REFERENCES analyst(id),

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS message (
    id SERIAL PRIMARY KEY,

    conversation_id INTEGER NOT NULL REFERENCES conversation(id),

    sender_id INTEGER NOT NULL,

    content TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    read_at TIMESTAMPTZ NULL
);