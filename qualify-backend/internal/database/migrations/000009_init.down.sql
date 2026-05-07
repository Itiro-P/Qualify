ALTER IF EXISTS "user" (
    DROP COLUMN IF EXISTS password_hash;
);

ALTER TABLE IF EXISTS certification (
    DROP COLUMN IF EXISTS institution;
);