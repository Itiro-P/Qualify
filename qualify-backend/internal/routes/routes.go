package routes

import (
	"main/internal/database/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func SetupRoutes(router *gin.Engine, conn *pgx.Conn) {
	router.GET("/user/:id", handlers.GetUser(conn))

	router.GET("/analyst/:id", handlers.GetAnalyst(conn))
	router.GET("/analysts", handlers.GetAnalysts(conn))

	router.GET("/client/:id", handlers.GetClient(conn))
	router.GET("/clients", handlers.GetClients(conn))

	router.GET("/review/:id", handlers.GetReview(conn))
	router.GET("/reviews", handlers.GetReviews(conn))
}
