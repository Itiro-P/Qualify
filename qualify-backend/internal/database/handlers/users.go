package handlers

import (
	"main/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func GetUser(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var user pkg.User
		err := conn.QueryRow(c.Request.Context(), "SELECT id, name FROM \"user\" WHERE id = $1", id).Scan(&user.Id, &user.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar usuário"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
