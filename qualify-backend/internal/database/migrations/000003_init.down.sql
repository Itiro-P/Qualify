DROP TABLE IF EXISTS certification;

DROP TABLE IF EXISTS analyst_certification;

DROP TABLE IF EXISTS proposal_letter;

DROP TABLE IF EXISTS "service";
    
DROP TABLE IF EXISTS review;

ALTER TABLE review DROP CONSTRAINT IF EXISTS fk_review_service;
ALTER TABLE review DROP COLUMN IF EXISTS service_id;
