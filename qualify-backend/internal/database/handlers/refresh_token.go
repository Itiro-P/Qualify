package handlers

import (
	"main/internal/utils"
	"main/pkg"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RefreshToken godoc
// @Summary Atualizar token de acesso
// @Description Obtém um novo token de acesso usando um refresh token válido.
// @Tags Autenticação
// @Accept json
// @Produce json
// @Param token query string false "Refresh token"
// @Success 200 {object} pkg.TokenResponse "Novo token de acesso gerado"
// @Failure 400 {object} pkg.ErrorResponse "Requisição inválida"
// @Failure 401 {object} pkg.ErrorResponse "Refresh token inválido ou expirado"
// @Failure 500 {object} pkg.ErrorResponse "Erro interno do servidor"
// @Router /auth/refresh [post]
func RefreshToken(conn *pgxpool.Pool) gin.HandlerFunc {
	jwtManager := utils.NewJWTManager()

	return func(c *gin.Context) {
		var request pkg.RefreshTokenRequest

		// Validate input
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, pkg.ErrorResponse{
				Error:     "invalid_request",
				Message:   "Refresh token is required",
				Code:      "MISSING_REFRESH_TOKEN",
				Timestamp: time.Now(),
			})
			return
		}

		// Validate refresh token format
		if err := utils.ValidateRefreshToken(request.RefreshToken); err != nil {
			c.JSON(http.StatusBadRequest, pkg.ErrorResponse{
				Error:     "invalid_token",
				Message:   err.Error(),
				Code:      "INVALID_TOKEN_FORMAT",
				Timestamp: time.Now(),
			})
			return
		}

		// Hash the token for database lookup
		tokenHash := utils.HashToken(request.RefreshToken)

		// Fetch refresh token from database
		var refreshToken pkg.RefreshToken
		var userID int
		var userName string
		var userEmail string

		err := conn.QueryRow(c.Request.Context(),
			`SELECT rt.id, rt.user_id, rt.token_hash, rt.expires_at, rt.revoked_at, 
			        u.name, u.email
			 FROM refresh_token rt
			 JOIN "user" u ON rt.user_id = u.id
			 WHERE rt.token_hash = $1`,
			tokenHash).Scan(&refreshToken.ID, &userID, &refreshToken.TokenHash,
			&refreshToken.ExpiresAt, &refreshToken.RevokedAt, &userName, &userEmail)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusUnauthorized, pkg.ErrorResponse{
					Error:     "invalid_token",
					Message:   "Refresh token is invalid or has been revoked",
					Code:      "REFRESH_TOKEN_NOT_FOUND",
					Timestamp: time.Now(),
				})
				return
			}

			c.JSON(http.StatusInternalServerError, pkg.ErrorResponse{
				Error:     "database_error",
				Message:   "Failed to verify refresh token",
				Code:      "INTERNAL_ERROR",
				Timestamp: time.Now(),
			})
			return
		}

		// Check if token is revoked
		if refreshToken.RevokedAt != nil {
			c.JSON(http.StatusUnauthorized, pkg.ErrorResponse{
				Error:     "token_revoked",
				Message:   "Refresh token has been revoked",
				Code:      "TOKEN_REVOKED",
				Timestamp: time.Now(),
			})
			return
		}

		// Check if token is expired
		if refreshToken.ExpiresAt.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, pkg.ErrorResponse{
				Error:     "token_expired",
				Message:   "Refresh token has expired",
				Code:      "TOKEN_EXPIRED",
				Timestamp: time.Now(),
			})
			return
		}

		// Generate new access token
		user := &pkg.User{
			Id:    userID,
			Name:  userName,
			Email: userEmail,
		}

		newAccessToken, expiresAt, err := jwtManager.GenerateAccessToken(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.ErrorResponse{
				Error:     "token_generation_failed",
				Message:   "Failed to generate new access token",
				Code:      "TOKEN_ERROR",
				Timestamp: time.Now(),
			})
			return
		}

		// Generate new refresh token (token rotation)
		newRefreshToken, newRefreshExpiresAt, err := jwtManager.GenerateRefreshToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.ErrorResponse{
				Error:     "token_generation_failed",
				Message:   "Failed to generate new refresh token",
				Code:      "TOKEN_ERROR",
				Timestamp: time.Now(),
			})
			return
		}

		// Revoke old refresh token
		now := time.Now()
		_, err = conn.Exec(c.Request.Context(),
			`UPDATE refresh_token SET revoked_at = $1 WHERE id = $2`,
			now, refreshToken.ID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.ErrorResponse{
				Error:     "token_rotation_failed",
				Message:   "Failed to rotate refresh token",
				Code:      "TOKEN_ERROR",
				Timestamp: time.Now(),
			})
			return
		}

		// Store new refresh token
		newTokenHash := utils.HashToken(newRefreshToken)
		clientIP := c.ClientIP()
		userAgent := c.Request.Header.Get("User-Agent")

		_, err = conn.Exec(c.Request.Context(),
			`INSERT INTO refresh_token (user_id, token_hash, expires_at, ip_address, user_agent)
			 VALUES ($1, $2, $3, $4, $5)`,
			userID, newTokenHash, newRefreshExpiresAt, clientIP, userAgent)

		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.ErrorResponse{
				Error:     "token_storage_failed",
				Message:   "Failed to store new refresh token",
				Code:      "TOKEN_ERROR",
				Timestamp: time.Now(),
			})
			return
		}

		// Return new tokens
		c.JSON(http.StatusOK, pkg.TokenResponse{
			AccessToken:  newAccessToken,
			RefreshToken: newRefreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    int64(jwtManager.GetTokenExpiryDuration().Seconds()),
			ExpiresAt:    expiresAt,
		})
	}
}
