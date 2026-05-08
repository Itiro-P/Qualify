package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"main/pkg"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JWTManager handles JWT token generation and validation
type JWTManager struct {
	secretKey     []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// jwtCustomClaims defines the JWT payload structure
type jwtCustomClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

// NewJWTManager creates a new JWT manager instance
func NewJWTManager() *JWTManager {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "your-secret-key-change-in-production"
	}

	return &JWTManager{
		secretKey:     []byte(secretKey),
		accessExpiry:  15 * time.Minute,
		refreshExpiry: 7 * 24 * time.Hour,
	}
}

// GenerateAccessToken creates a signed JWT access token
func (jm *JWTManager) GenerateAccessToken(user *pkg.User) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(jm.accessExpiry)

	claims := jwtCustomClaims{
		UserID: user.Id,
		Email:  user.Email,
		Name:   user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(jm.secretKey)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return signed, expiresAt, nil
}

// GenerateRefreshToken creates a cryptographically random refresh token
func (jm *JWTManager) GenerateRefreshToken() (string, time.Time, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate token: %w", err)
	}

	tokenStr := base64.URLEncoding.EncodeToString(token)
	expiresAt := time.Now().Add(jm.refreshExpiry)

	return tokenStr, expiresAt, nil
}

// ValidateAccessToken parses and validates a JWT access token
func (jm *JWTManager) ValidateAccessToken(tokenStr string) (*pkg.JWTClaims, error) {
	if tokenStr == "" {
		return nil, fmt.Errorf("token is required")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &jwtCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Enforce the expected signing method to prevent algorithm confusion attacks
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jm.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*jwtCustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return &pkg.JWTClaims{
		UserID:    claims.UserID,
		Email:     claims.Email,
		Name:      claims.Name,
		ExpiresAt: claims.ExpiresAt.Time,
		IssuedAt:  claims.IssuedAt.Time,
		NotBefore: claims.NotBefore.Time,
	}, nil
}

// HashToken creates a SHA-256 hash of a token for safe storage
func HashToken(token string) string {
	sha := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sha[:])
}

// ValidatePassword checks if the provided password matches the bcrypt hash
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
