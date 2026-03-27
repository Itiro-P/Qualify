package main

import (
	"context"
	"fmt"
	"main/internal/database/handlers"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	router := gin.Default()

	router.GET("/users", handlers.GetUsers(conn))
	router.GET("/users/:id", handlers.GetUser(conn))

	router.Run(":8001")
}
