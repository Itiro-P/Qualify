package main

import (
	"context"
	"fmt"
	"main/internal/routes"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "main/docs"
)

func main() {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	router := gin.Default()

	routes.SetupRoutes(router, pool)
	router.Run(":8001")
}
