-- Drop refresh token table and remove new columns from user table
DROP INDEX IF EXISTS public.idx_refresh_token_token_hash;
DROP INDEX IF EXISTS public.idx_refresh_token_expires_at;
DROP INDEX IF EXISTS public.idx_refresh_token_user_id;
DROP TABLE IF EXISTS public.refresh_token;

ALTER TABLE public."user" DROP COLUMN IF EXISTS account_locked;
ALTER TABLE public."user" DROP COLUMN IF EXISTS failed_login_attempts;
ALTER TABLE public."user" DROP COLUMN IF EXISTS last_login_at;
ALTER TABLE public."user" DROP COLUMN IF EXISTS last_password_change;
