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
// @Param user body pkg.UserLogin true "Objeto login"
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

		if err := c.ShouldBindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Email and password are required"))
			return
		}

		credentials.Email = utils.SanitizeEmail(credentials.Email)
		if err := utils.ValidateEmail(credentials.Email); err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

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
				c.JSON(http.StatusUnauthorized, pkg.Unauthorized(c.FullPath(), "Invalid email or password"))
				return
			}
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}

		if accountLocked {
			c.JSON(http.StatusUnauthorized, pkg.Unauthorized(c.FullPath(), "Account is locked due to too many failed login attempts. Please contact support."))
			return
		}

		if err = utils.ValidatePassword(credentials.Password, storedHash); err != nil {
			failedLoginAttempts++

			if failedLoginAttempts >= 5 {
				conn.Exec(c.Request.Context(),
					`UPDATE "user" SET account_locked = true, failed_login_attempts = $1 WHERE id = $2`,
					failedLoginAttempts, user.Id)
				c.JSON(http.StatusUnauthorized, pkg.Unauthorized(c.FullPath(), "Account locked due to too many failed attempts"))
				return
			}

			conn.Exec(c.Request.Context(),
				`UPDATE "user" SET failed_login_attempts = $1 WHERE id = $2`,
				failedLoginAttempts, user.Id)
			c.JSON(http.StatusUnauthorized, pkg.Unauthorized(c.FullPath(), "Invalid email or password"))
			return
		}

		conn.Exec(c.Request.Context(),
			`UPDATE "user" SET failed_login_attempts = 0, last_login_at = NOW() WHERE id = $1`,
			user.Id)

		accessToken, accessExpiresAt, err := jwtManager.GenerateAccessToken(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), "Failed to generate access token"))
			return
		}

		refreshToken, refreshExpiresAt, err := jwtManager.GenerateRefreshToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), "Failed to generate refresh token"))
			return
		}

		tokenHash := utils.HashToken(refreshToken)
		_, err = conn.Exec(c.Request.Context(),
			`INSERT INTO refresh_token (user_id, token_hash, expires_at, ip_address, user_agent)
			 VALUES ($1, $2, $3, $4, $5)`,
			user.Id, tokenHash, refreshExpiresAt, c.ClientIP(), c.Request.Header.Get("User-Agent"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), "Failed to store refresh token"))
			return
		}

		c.JSON(http.StatusOK, pkg.LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    int64(jwtManager.GetTokenExpiryDuration().Seconds()),
			ExpiresAt:    accessExpiresAt,
		})
	}
}
