package routes

import (
	"main/internal/database/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func SetupRoutes(router *gin.Engine, conn *pgx.Conn) {
	router.GET("/users", handlers.GetUsers(conn))
	router.GET("/users/:id", handlers.GetUser(conn))
}
