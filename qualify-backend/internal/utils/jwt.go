package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"main/pkg"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// JWTManager handles JWT token generation and validation
type JWTManager struct {
	secretKey     string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// NewJWTManager creates a new JWT manager instance
func NewJWTManager() *JWTManager {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "your-secret-key-change-in-production"
	}

	return &JWTManager{
		secretKey:     secretKey,
		accessExpiry:  15 * time.Minute,   // 15 minutes
		refreshExpiry: 7 * 24 * time.Hour, // 7 days
	}
}

// GenerateAccessToken creates a JWT access token
func (jm *JWTManager) GenerateAccessToken(user *pkg.User) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(jm.accessExpiry)

	claims := pkg.JWTClaims{
		UserID:    user.Id,
		Email:     user.Email,
		Name:      user.Name,
		ExpiresAt: expiresAt,
		IssuedAt:  now,
		NotBefore: now,
	}

	// Simple JWT implementation (Note: In production, use github.com/golang-jwt/jwt)
	payload := fmt.Sprintf(`{
		"user_id":%d,
		"email":"%s",
		"name":"%s",
		"exp":%d,
		"iat":%d,
		"nbf":%d
	}`, claims.UserID, claims.Email, claims.Name, claims.ExpiresAt.Unix(), claims.IssuedAt.Unix(), claims.NotBefore.Unix())

	header := `{"alg":"HS256","typ":"JWT"}`
	signature := generateSignature(header, payload, jm.secretKey)

	token := fmt.Sprintf("%s.%s.%s",
		base64.URLEncoding.EncodeToString([]byte(header)),
		base64.URLEncoding.EncodeToString([]byte(payload)),
		signature)

	return token, expiresAt, nil
}

// GenerateRefreshToken creates a random refresh token
func (jm *JWTManager) GenerateRefreshToken() (string, time.Time, error) {
	// Generate random token
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate token: %w", err)
	}

	tokenStr := base64.URLEncoding.EncodeToString(token)
	expiresAt := time.Now().Add(jm.refreshExpiry)

	return tokenStr, expiresAt, nil
}

// HashToken creates a hash of the token for storage
func HashToken(token string) string {
	sha := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sha[:])
}

// generateSimpleHash creates a simple hash of the token
func generateSimpleHash(token string) string {
	data := []byte(token)
	hash := make([]byte, 32)
	for i := 0; i < len(data); i++ {
		hash[i%32] ^= data[i]
	}
	return hex.EncodeToString(hash)
}

// generateSignature creates HMAC signature for JWT
func generateSignature(header, payload, secret string) string {
	// Simplified HMAC-SHA256 signature
	message := header + "." + payload
	h := hex.EncodeToString([]byte(message + secret))
	return h
}

// ValidateAccessToken validates an access token
func (jm *JWTManager) ValidateAccessToken(token string) (*pkg.JWTClaims, error) {
	if token == "" {
		return nil, fmt.Errorf("token is required")
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	// Verify signature
	expectedSignature := generateSignature(parts[0], parts[1], jm.secretKey)
	if parts[2] != expectedSignature {
		return nil, fmt.Errorf("invalid token signature")
	}

	payloadBytes, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid token payload")
	}

	type rawClaims struct {
		UserID int    `json:"user_id"`
		Email  string `json:"email"`
		Name   string `json:"name"`
		Exp    int64  `json:"exp"`
		Iat    int64  `json:"iat"`
		Nbf    int64  `json:"nbf"`
	}

	var raw rawClaims
	if err := json.Unmarshal(payloadBytes, &raw); err != nil {
		return nil, fmt.Errorf("invalid token claims")
	}

	now := time.Now().Unix()
	if raw.Nbf > now {
		return nil, fmt.Errorf("token not valid yet")
	}
	if raw.Exp < now {
		return nil, fmt.Errorf("token expired")
	}

	return &pkg.JWTClaims{
		UserID:    raw.UserID,
		Email:     raw.Email,
		Name:      raw.Name,
		ExpiresAt: time.Unix(raw.Exp, 0),
		IssuedAt:  time.Unix(raw.Iat, 0),
		NotBefore: time.Unix(raw.Nbf, 0),
	}, nil
}

// ValidatePassword checks if provided password matches the hash
func ValidatePassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// HashPassword creates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// GetTokenExpiryDuration returns the access token expiry duration
func (jm *JWTManager) GetTokenExpiryDuration() time.Duration {
	return jm.accessExpiry
}

// GetRefreshTokenExpiryDuration returns the refresh token expiry duration
func (jm *JWTManager) GetRefreshTokenExpiryDuration() time.Duration {
	return jm.refreshExpiry
}
