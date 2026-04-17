<<<<<<< HEAD
DROP TABLE IF EXISTS certification;

DROP TABLE IF EXISTS analyst_certification;

DROP TABLE IF EXISTS proposal_letter;

DROP TABLE IF EXISTS "service";
    
DROP TABLE IF EXISTS review;

ALTER TABLE IF EXISTS review
    DROP COLUMN IF EXISTS service_id,
    ADD COLUMN IF NOT EXISTS analyst_id,
    ADD COLUMN IF NOT EXISTS client_id,
    ADD FOREIGN KEY (analyst_id) REFERENCES analyst (user_id) ON DELETE CASCADE;
    ADD FOREIGN KEY (client_id) REFERENCES client (user_id) ON DELETE CASCADE;
=======
ALTER TABLE "user" 
DROP CONSTRAINT email_format;
>>>>>>> feat-criar-migrations-e-tabelas
