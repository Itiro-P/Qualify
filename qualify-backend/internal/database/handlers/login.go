package handlers

import (
	"main/internal/utils"
	"main/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Login godoc
// @Summary Login de usuário com geração de JWT
// @Description Autentica usuário com email e senha, retornando tokens JWT de acesso e refresh.
// @Tags Autenticação
// @Accept json
// @Produce json
// @Param email query string false "Email do usuário"
// @Param password query string false "Senha do usuário"
// @Success 200 {object} pkg.LoginResponse "Autenticado com sucesso e tokens retornados"
// @Failure 400 {object} pkg.ErrorResponse "Entrada inválida ou erro de validação"
// @Failure 401 {object} pkg.ErrorResponse "Credenciais inválidas"
// @Failure 429 {object} pkg.ErrorResponse "Muitas tentativas de login"
// @Failure 500 {object} pkg.ErrorResponse "Erro interno do servidor"
// @Router /auth/login [post]
func Login(conn *pgxpool.Pool) gin.HandlerFunc {
	jwtManager := utils.NewJWTManager()

	return func(c *gin.Context) {
		var credentials pkg.UserLogin

		// Validate input binding
		if err := c.ShouldBindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid_input",
				"message": "Email and password are required",
				"code":    "MISSING_CREDENTIALS",
			})
			return
		}

		// Validate email format
		credentials.Email = utils.SanitizeEmail(credentials.Email)
		if err := utils.ValidateEmail(credentials.Email); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid_email",
				"message": err.Error(),
				"code":    "INVALID_EMAIL_FORMAT",
			})
			return
		}

		// Fetch user from database
		var user pkg.User
		var storedHash string
		var accountLocked bool
		var failedLoginAttempts int

		err := conn.QueryRow(c.Request.Context(),
			`SELECT id, name, email, password_hash, account_locked, failed_login_attempts 
			 FROM "user" WHERE email = $1`,
			credentials.Email).Scan(&user.Id, &user.Name, &user.Email, &storedHash, &accountLocked, &failedLoginAttempts)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":   "invalid_credentials",
					"message": "Invalid email or password",
					"code":    "AUTHENTICATION_FAILED",
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "database_error",
				"message": "An error occurred while processing your request",
				"code":    "INTERNAL_ERROR",
			})
			return
		}

		// Check if account is locked
		if accountLocked {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "account_locked",
				"message": "Account is locked due to too many failed login attempts. Please contact support.",
				"code":    "ACCOUNT_LOCKED",
			})
			return
		}

		// Validate password
		err = utils.ValidatePassword(credentials.Password, storedHash)
		if err != nil {
			// Increment failed login attempts
			failedLoginAttempts++
			maxAttempts := 5

			if failedLoginAttempts >= maxAttempts {
				// Lock account
				updateErr := conn.QueryRow(c.Request.Context(),
					`UPDATE "user" SET account_locked = true, failed_login_attempts = $1 WHERE id = $2`,
					failedLoginAttempts, user.Id).Scan()
				if updateErr != nil && updateErr != pgx.ErrNoRows {
					// Log error but don't fail the login attempt response
				}

				c.JSON(http.StatusUnauthorized, gin.H{
					"error":   "account_locked",
					"message": "Account locked due to too many failed attempts",
					"code":    "ACCOUNT_LOCKED_MAX_ATTEMPTS",
				})
				return
			}

			// Update failed attempts
			conn.QueryRow(c.Request.Context(),
				`UPDATE "user" SET failed_login_attempts = $1 WHERE id = $2`,
				failedLoginAttempts, user.Id).Scan()

			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid_credentials",
				"message": "Invalid email or password",
				"code":    "AUTHENTICATION_FAILED",
			})
			return
		}

		// Reset failed login attempts on successful login
		conn.QueryRow(c.Request.Context(),
			`UPDATE "user" SET failed_login_attempts = 0, last_login_at = NOW() WHERE id = $1`,
			user.Id).Scan()

		// Generate tokens
		accessToken, accessExpiresAt, err := jwtManager.GenerateAccessToken(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "token_generation_failed",
				"message": "Failed to generate access token",
				"code":    "TOKEN_ERROR",
			})
			return
		}

		refreshToken, refreshExpiresAt, err := jwtManager.GenerateRefreshToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "token_generation_failed",
				"message": "Failed to generate refresh token",
				"code":    "TOKEN_ERROR",
			})
			return
		}

		// Hash and store refresh token
		tokenHash := utils.HashToken(refreshToken)
		clientIP := c.ClientIP()
		userAgent := c.Request.Header.Get("User-Agent")

		_, err = conn.Exec(c.Request.Context(),
			`INSERT INTO refresh_token (user_id, token_hash, expires_at, ip_address, user_agent)
			 VALUES ($1, $2, $3, $4, $5)`,
			user.Id, tokenHash, refreshExpiresAt, clientIP, userAgent)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "token_storage_failed",
				"message": "Failed to store refresh token",
				"code":    "TOKEN_ERROR",
			})
			return
		}

		// Return login response
		c.JSON(http.StatusOK, pkg.LoginResponse{
			User:         user,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    int64(jwtManager.GetTokenExpiryDuration().Seconds()),
			ExpiresAt:    accessExpiresAt,
		})
	}
}
