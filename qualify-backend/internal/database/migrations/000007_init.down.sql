-- Revert mean_rating to NOT NULL with default 0
ALTER TABLE analyst ALTER COLUMN mean_rating SET DEFAULT 0;
ALTER TABLE analyst ALTER COLUMN mean_rating SET NOT NULL;