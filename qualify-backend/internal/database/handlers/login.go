package handlers

import (
	"main/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// Login godoc
// @Summary Fazer login
// @Description Faz login com os dados fornecidos
// @Tags Usuários
// @Accept json
// @Produce json
// @Param user body pkg.UserLogin true "Objeto do usuário"
// @Success 200 {object} pkg.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /register [post]
func Login(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var credentials pkg.UserLogin
		if err := c.ShouldBindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email e senha são obrigatórios"})
			return
		}

		var user pkg.User
		var storedHash string

		err := conn.QueryRow(c.Request.Context(),
			`SELECT id, name, email, password_hash FROM "user" WHERE email = $1`,
			credentials.Email).Scan(&user.Id, &user.Name, &user.Email, &storedHash)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário ou senha inválidos"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(credentials.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário ou senha inválidos"})
			return
		}
		c.JSON(http.StatusOK, pkg.UserResponse{User: user})
	}
}
