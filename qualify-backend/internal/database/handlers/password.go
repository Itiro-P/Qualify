package handlers

import (
	"context"
	"main/internal/utils"
	"main/pkg"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ChangePassword godoc
// @Summary Alterar senha do usuário
// @Description Permite que usuário autenticado altere sua senha com validação. Requer token de acesso válido
// @Tags Autenticação
// @Accept json
// @Produce json
// @Param send_email query boolean false "Enviar confirmação por email (padrão: true)"
// @Param request body pkg.PasswordChangeRequest true "Senha atual e nova senha"
// @Success 200 {object} pkg.PasswordChangeResponse "Senha alterada com sucesso"
// @Failure 400 {object} pkg.ErrorResponse "Requisição inválida ou senha fraca"
// @Failure 401 {object} pkg.ErrorResponse "Senha atual inválida"
// @Failure 500 {object} pkg.ErrorResponse "Erro interno do servidor"
// @Router /auth/change-password [post]
// @Security Bearer
func ChangePassword(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request pkg.PasswordChangeRequest

		// Validate input
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, pkg.ErrorResponse{
				Error:     "invalid_request",
				Message:   "Current password and new password are required",
				Code:      "MISSING_FIELDS",
				Timestamp: time.Now(),
			})
			return
		}

		// Extract user ID from context (set by auth middleware)
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, pkg.ErrorResponse{
				Error:     "unauthorized",
				Message:   "User information not found",
				Code:      "UNAUTHORIZED",
				Timestamp: time.Now(),
			})
			return
		}

		userID, ok := userIDInterface.(int)
		if !ok {
			c.JSON(http.StatusUnauthorized, pkg.ErrorResponse{
				Error:     "unauthorized",
				Message:   "Invalid user information",
				Code:      "UNAUTHORIZED",
				Timestamp: time.Now(),
			})
			return
		}

		// Validate current password is not empty
		if request.CurrentPassword == "" {
			c.JSON(http.StatusBadRequest, pkg.ErrorResponse{
				Error:     "invalid_password",
				Message:   "Current password is required",
				Code:      "MISSING_CURRENT_PASSWORD",
				Timestamp: time.Now(),
			})
			return
		}

		// Validate new password strength
		validationErrors := utils.ValidatePasswordStrength(request.NewPassword)
		if len(validationErrors) > 0 {
			errorMap := utils.BuildValidationErrorMap(validationErrors)
			c.JSON(http.StatusBadRequest, pkg.ErrorResponse{
				Error:            "weak_password",
				Message:          "Password does not meet strength requirements",
				Code:             "PASSWORD_WEAKNESS",
				Timestamp:        time.Now(),
				ValidationErrors: errorMap,
			})
			return
		}

		// Verify current password is different from new password
		if request.CurrentPassword == request.NewPassword {
			c.JSON(http.StatusBadRequest, pkg.ErrorResponse{
				Error:     "invalid_password",
				Message:   "New password must be different from current password",
				Code:      "PASSWORD_SAME",
				Timestamp: time.Now(),
			})
			return
		}

		// Fetch current password hash from database
		var storedHash string
		var userEmail string
		var userName string

		err := conn.QueryRow(c.Request.Context(),
			`SELECT password_hash, email, name FROM "user" WHERE id = $1`,
			userID).Scan(&storedHash, &userEmail, &userName)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusUnauthorized, pkg.ErrorResponse{
					Error:     "user_not_found",
					Message:   "User not found",
					Code:      "USER_NOT_FOUND",
					Timestamp: time.Now(),
				})
				return
			}

			c.JSON(http.StatusInternalServerError, pkg.ErrorResponse{
				Error:     "database_error",
				Message:   "Failed to retrieve user information",
				Code:      "INTERNAL_ERROR",
				Timestamp: time.Now(),
			})
			return
		}

		// Verify current password
		err = utils.ValidatePassword(request.CurrentPassword, storedHash)
		if err != nil {
			c.JSON(http.StatusUnauthorized, pkg.ErrorResponse{
				Error:     "invalid_password",
				Message:   "Current password is incorrect",
				Code:      "INVALID_CURRENT_PASSWORD",
				Timestamp: time.Now(),
			})
			return
		}

		// Hash new password
		newPasswordHash, err := utils.HashPassword(request.NewPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.ErrorResponse{
				Error:     "hashing_failed",
				Message:   "Failed to process password change",
				Code:      "INTERNAL_ERROR",
				Timestamp: time.Now(),
			})
			return
		}

		// Update password in database
		now := time.Now()
		_, err = conn.Exec(c.Request.Context(),
			`UPDATE "user" 
			 SET password_hash = $1, last_password_change = $2, failed_login_attempts = 0
			 WHERE id = $3`,
			newPasswordHash, now, userID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.ErrorResponse{
				Error:     "password_update_failed",
				Message:   "Failed to update password",
				Code:      "INTERNAL_ERROR",
				Timestamp: time.Now(),
			})
			return
		}

		// Revoke all refresh tokens (force re-login on all devices for security)
		_ = RevokeAllTokens(c.Request.Context(), conn, userID)

		// Send email notification if requested
		if request.SendNotification {
			// TODO: Implement email notification
			// sendPasswordChangeEmail(userEmail, userName)
		}

		c.JSON(http.StatusOK, pkg.PasswordChangeResponse{
			Success:   true,
			Message:   "Password successfully changed",
			ChangedAt: now,
		})
	}
}

// ResetPassword godoc
// @Summary Solicitar redefinição de senha
// @Description Inicia o processo de redefinição de senha enviando link por email
// @Tags Autenticação
// @Accept json
// @Produce json
// @Param email query string true "Endereço de email do usuário"
// @Success 200 {object} pkg.SuccessResponse "Email de redefinição enviado"
// @Failure 400 {object} pkg.ErrorResponse "Email inválido"
// @Failure 404 {object} pkg.ErrorResponse "Usuário não encontrado"
// @Failure 500 {object} pkg.ErrorResponse "Erro interno do servidor"
// @Router /auth/reset-password [post]
func ResetPassword(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Query("email")

		// Validate email
		email = utils.SanitizeEmail(email)
		if err := utils.ValidateEmail(email); err != nil {
			c.JSON(http.StatusBadRequest, pkg.ErrorResponse{
				Error:     "invalid_email",
				Message:   err.Error(),
				Code:      "INVALID_EMAIL",
				Timestamp: time.Now(),
			})
			return
		}

		// Check if user exists
		var userID int
		var userName string

		err := conn.QueryRow(c.Request.Context(),
			`SELECT id, name FROM "user" WHERE email = $1`,
			email).Scan(&userID, &userName)

		if err != nil {
			if err == pgx.ErrNoRows {
				// Don't reveal if user exists (security best practice)
				c.JSON(http.StatusOK, pkg.SuccessResponse{
					Success: true,
					Message: "If the email exists, a password reset link has been sent",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, pkg.ErrorResponse{
				Error:     "database_error",
				Message:   "Failed to process request",
				Code:      "INTERNAL_ERROR",
				Timestamp: time.Now(),
			})
			return
		}

		// TODO: Generate password reset token
		// Generate a unique reset token and store it with expiry
		// resetToken := generateResetToken()
		// storeResetToken(userID, resetToken, expiresAt)
		// sendResetEmail(email, resetToken)

		// For now, return success
		c.JSON(http.StatusOK, pkg.SuccessResponse{
			Success: true,
			Message: "Password reset link has been sent to your email",
			Data: gin.H{
				"email": email,
			},
		})
	}
}

// ConfirmPasswordReset godoc
// @Summary Confirmar redefinição de senha com token
// @Description Define nova senha usando token de redefinição enviado por email
// @Tags Autenticação
// @Accept json
// @Produce json
// @Param token query string true "Token de redefinição de senha"
// @Param new_password query string true "Nova senha (mínimo 8 caracteres, com letras maiúsculas, minúsculas, números e símbolos)"
// @Param send_email query boolean false "Enviar confirmação por email (padrão: true)"
// @Success 200 {object} pkg.SuccessResponse "Redefinição de senha bem-sucedida"
// @Failure 400 {object} pkg.ErrorResponse "Token inválido ou senha fraca"
// @Failure 500 {object} pkg.ErrorResponse "Erro interno do servidor"
// @Router /auth/reset-password/confirm [post]
// @Router /auth/reset-password/confirm [post]
func ConfirmPasswordReset(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request pkg.PasswordResetConfirmRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, pkg.ErrorResponse{
				Error:     "invalid_request",
				Message:   "Token and new password are required",
				Code:      "MISSING_FIELDS",
				Timestamp: time.Now(),
			})
			return
		}

		// Validate new password strength
		validationErrors := utils.ValidatePasswordStrength(request.NewPassword)
		if len(validationErrors) > 0 {
			errorMap := utils.BuildValidationErrorMap(validationErrors)
			c.JSON(http.StatusBadRequest, pkg.ErrorResponse{
				Error:            "weak_password",
				Message:          "Password does not meet strength requirements",
				Code:             "PASSWORD_WEAKNESS",
				Timestamp:        time.Now(),
				ValidationErrors: errorMap,
			})
			return
		}

		// TODO: Validate reset token
		// Check if token exists and is not expired
		// If valid, get the associated user ID
		// Update password
		// Invalidate token

		// Placeholder response
		c.JSON(http.StatusOK, pkg.SuccessResponse{
			Success: true,
			Message: "Password reset successful",
		})
	}
}

// Context helper function (add to context before calling these handlers)
func setUserIDInContext(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, "user_id", userID)
}
