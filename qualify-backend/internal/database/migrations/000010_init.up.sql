-- Create refresh_token table for token management
CREATE TABLE IF NOT EXISTS public.refresh_token (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES public."user"(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMP WITH TIME ZONE,
    ip_address TEXT,
    user_agent TEXT
);

-- Create index for efficient lookups
CREATE INDEX IF NOT EXISTS idx_refresh_token_user_id ON public.refresh_token(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_token_expires_at ON public.refresh_token(expires_at);
CREATE INDEX IF NOT EXISTS idx_refresh_token_token_hash ON public.refresh_token(token_hash);

-- Add account_locked column to user table for security
ALTER TABLE public."user" ADD COLUMN IF NOT EXISTS account_locked BOOLEAN DEFAULT FALSE;
ALTER TABLE public."user" ADD COLUMN IF NOT EXISTS failed_login_attempts INTEGER DEFAULT 0;
ALTER TABLE public."user" ADD COLUMN IF NOT EXISTS last_login_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE public."user" ADD COLUMN IF NOT EXISTS last_password_change TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
