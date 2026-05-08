package pkg

import "time"

// Token response with access and refresh tokens
type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int64     `json:"expires_in"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// Login response with user and tokens
type LoginResponse struct {
	User         User      `json:"user"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int64     `json:"expires_in"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// Refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Password change request
type PasswordChangeRequest struct {
	CurrentPassword  string `json:"current_password" binding:"required,min=8"`
	NewPassword      string `json:"new_password" binding:"required,min=12"`
	SendNotification bool   `json:"send_notification,omitempty"`
}

// Password change response
type PasswordChangeResponse struct {
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
	ChangedAt time.Time `json:"changed_at"`
}

// Password reset request
type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// Password reset confirmation request
type PasswordResetConfirmRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=12"`
}

// Logout request
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token,omitempty"`
}

// JWT claims structure
type JWTClaims struct {
	UserID    int       `json:"user_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	ExpiresAt time.Time `json:"exp"`
	IssuedAt  time.Time `json:"iat"`
	NotBefore time.Time `json:"nbf"`
}

// Refresh token database model
type RefreshToken struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	TokenHash string     `json:"-"`
	ExpiresAt time.Time  `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at"`
	IPAddress string     `json:"ip_address"`
	UserAgent string     `json:"user_agent"`
}

// Error response structure
type ErrorResponse struct {
	Error            string            `json:"error"`
	Message          string            `json:"message"`
	Code             string            `json:"code,omitempty"`
	Timestamp        time.Time         `json:"timestamp"`
	Path             string            `json:"path,omitempty"`
	ValidationErrors map[string]string `json:"validation_errors,omitempty"`
}

// Success response structure
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
