package test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func ToPtrFloat64(n float64) *float64 {
	f := n
	return &f
}

var TestPool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := postgres.Run(ctx, "postgres:17-alpine",
		postgres.WithDatabase("godb"),
		postgres.WithUsername("gouser"),
		postgres.WithPassword("gopassword"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	dsn, _ := container.ConnectionString(ctx, "sslmode=disable")

	runMigrations(dsn)

	TestPool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	TestPool.Close()
	container.Terminate(ctx)
	os.Exit(code)
}

func runMigrations(dsn string) {
	m, err := migrate.New("file://../internal/database/migrations", dsn)
	if err != nil {
		log.Fatalf("Erro ao iniciar migração: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Erro ao subir migração: %v", err)
	}
}
