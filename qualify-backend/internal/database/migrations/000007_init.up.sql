-- Make mean_rating nullable for analyst table
ALTER TABLE analyst ALTER COLUMN mean_rating DROP NOT NULL;
ALTER TABLE analyst ALTER COLUMN mean_rating DROP DEFAULT;