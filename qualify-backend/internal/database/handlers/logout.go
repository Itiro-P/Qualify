package handlers

import (
	"context"
	"main/internal/utils"
	"main/pkg"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Logout godoc
// @Summary Logout e invalidação de tokens
// @Description Revoga todos os refresh tokens do usuário autenticado, efetivamente desconectando todas as sessões ativas. Se um refresh_token específico for fornecido como query parameter, apenas aquele token será revogado (logout do dispositivo atual). Caso contrário, todos os tokens do usuário são revogados (logout global).
// @Description
// @Description **Comportamentos:**
// @Description - Sem parâmetro: Revoga TODOS os refresh tokens (logout de todos os dispositivos)
// @Description - Com refresh_token: Revoga APENAS aquele token específico (logout do dispositivo atual)
// @Tags Autenticação
// @Accept json
// @Produce json
// @Param refresh_token query string false "Refresh token específico para revogar (opcional). Se não fornecido, revoga todos os tokens do usuário"
// @Success 200 {object} pkg.SuccessResponse "Logout realizado com sucesso"
// @Failure 400 {object} pkg.ErrorResponse "Requisição inválida"
// @Failure 401 {object} pkg.ErrorResponse "Não autorizado - Token inválido ou expirado"
// @Failure 500 {object} pkg.ErrorResponse "Erro interno do servidor"
// @Security     BearerAuth
// @Router /auth/logout [post]
// @Security Bearer
func Logout(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request pkg.LogoutRequest

		// Optional: bind logout request (may contain refresh token)
		_ = c.ShouldBindJSON(&request)

		// Extract user ID from authenticated context
		userIDInterface, exists := c.Get("user_id")

		if !exists {
			// Try to extract from token (simplified - use proper JWT parsing)
			c.JSON(http.StatusUnauthorized, pkg.ErrorResponse{
				Error:     "invalid_token",
				Message:   "Cannot extract user information from token",
				Code:      "INVALID_TOKEN",
				Timestamp: time.Now(),
			})
			return
		}

		userID, ok := userIDInterface.(int)
		if !ok {
			c.JSON(http.StatusUnauthorized, pkg.ErrorResponse{
				Error:     "invalid_token",
				Message:   "Invalid user information in token",
				Code:      "INVALID_TOKEN",
				Timestamp: time.Now(),
			})
			return
		}

		now := time.Now()
		revokeCount := 0

		if request.RefreshToken != "" {
			// Logout só deste dispositivo
			tokenHash := utils.HashToken(request.RefreshToken)
			result, err := conn.Exec(c.Request.Context(),
				`UPDATE refresh_token SET revoked_at = $1
				WHERE user_id = $2 AND token_hash = $3 AND revoked_at IS NULL`,
				now, userID, tokenHash)
			if err != nil {
				c.JSON(http.StatusInternalServerError, pkg.ErrorResponse{
					Error:     "revocation_failed",
					Message:   "Failed to revoke refresh token",
					Code:      "TOKEN_ERROR",
					Timestamp: time.Now(),
				})
				return
			}
			revokeCount = int(result.RowsAffected())
		} else {
			// Logout global — revoga todos
			result, err := conn.Exec(c.Request.Context(),
				`UPDATE refresh_token SET revoked_at = $1
				WHERE user_id = $2 AND revoked_at IS NULL`,
				now, userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, pkg.ErrorResponse{
					Error:     "revocation_failed",
					Message:   "Failed to revoke refresh token",
					Code:      "TOKEN_ERROR",
					Timestamp: time.Now(),
				})
				return
			}
			revokeCount = int(result.RowsAffected())
		}

		// Log logout event (optional)
		// You can implement audit logging here
		_ = logLogoutEvent(c.Request.Context(), conn, userID, c.ClientIP())

		// Return success response
		c.JSON(http.StatusOK, pkg.SuccessResponse{
			Success: true,
			Message: "Successfully logged out from all devices",
			Data: gin.H{
				"tokens_revoked": revokeCount,
				"logout_time":    now,
			},
		})
	}
}

// Helper function to log logout events (optional audit logging)
func logLogoutEvent(ctx context.Context, conn *pgxpool.Pool, userID int, ipAddress string) error {
	// This is a placeholder for audit logging
	// In production, you might want to log logout events for security purposes
	// Example: store in a log table or send to a logging service
	return nil
}

// RevokeAllTokens revokes all tokens for a user (used for password changes, security events)
// This is an internal function that can be called from other handlers
func RevokeAllTokens(ctx context.Context, conn *pgxpool.Pool, userID int) error {
	now := time.Now()
	_, err := conn.Exec(ctx,
		`UPDATE refresh_token
		 SET revoked_at = $1
		 WHERE user_id = $2 AND revoked_at IS NULL`,
		now, userID)

	return err
}

// RevokeTokensExcept revokes all tokens except the specified ones (for token rotation)
func RevokeTokensExcept(ctx context.Context, conn *pgxpool.Pool, userID int, exceptTokenHash string) error {
	now := time.Now()
	_, err := conn.Exec(ctx,
		`UPDATE refresh_token
		 SET revoked_at = $1
		 WHERE user_id = $2 AND token_hash != $3 AND revoked_at IS NULL`,
		now, userID, exceptTokenHash)

	return err
}
