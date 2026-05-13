package services

import (
	"context"
	"main/pkg"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateUser godoc
// @Summary Criar usuário (serviço)
// @Description Insere um novo usuário no banco de dados. Função de serviço usada pelos handlers.
// @Tags Serviços
// @Accept json
// @Produce json
// @Param user body pkg.User true "Objeto do usuário"
// @Success 200 {object} pkg.UserRegister
// @Failure 500 {object} pkg.ErrorResponse
func CreateUser(ctx context.Context, conn *pgxpool.Pool, user *pkg.User) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// Adicionamos password_hash na lista de colunas e no VALUES ($9)
	query := `
        INSERT INTO "user" (
            name, email, phone, country_code, country_name, 
            country_state, city, timezone, password_hash
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id, time_created`

	err = tx.QueryRow(ctx, query,
		user.Name,
		user.Email,
		user.Phone,
		user.Country_code,
		user.Country_name,
		user.Country_state,
		user.City,
		user.Timezone,
		user.Password_hash,
	).Scan(&user.Id, &user.Time_created)

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// AssignAnalystRole godoc
// @Summary Atribuir papel de analista (serviço)
// @Description Atribui o papel de analista a um usuário existente e retorna o registro do analista.
// @Tags Serviços
// @Accept json
// @Produce json
// @Param userID path int true "ID do usuário"
// @Param hourlyRate query number true "Valor por hora"
// @Success 200 {object} pkg.Analyst
// @Failure 500 {object} pkg.ErrorResponse
func AssignAnalystRole(ctx context.Context, conn *pgxpool.Pool, userID int, hourlyRate float64) (*pkg.Analyst, error) {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var analyst pkg.Analyst
	if err := tx.QueryRow(ctx,
		`INSERT INTO analyst (id, hourly_rate, total_reviews, mean_rating)
		 SELECT id, $2, 0, NULL FROM "user" WHERE id = $1
		 RETURNING id`,
		userID, hourlyRate).Scan(&analyst.Id); err != nil {
		return nil, err
	}

	if err := tx.QueryRow(ctx,
		`SELECT u.id, u.name, u.email, u.phone, u.time_created,
		       u.country_code, u.country_name, u.country_state, u.city, u.timezone,
		       a.hourly_rate, a.total_reviews, a.mean_rating
		FROM "user" u
		JOIN analyst a ON a.id = u.id
		WHERE u.id = $1`,
		userID).Scan(
		&analyst.Id,
		&analyst.Name,
		&analyst.Email,
		&analyst.Phone,
		&analyst.Time_created,
		&analyst.Country_code,
		&analyst.Country_name,
		&analyst.Country_state,
		&analyst.City,
		&analyst.Timezone,
		&analyst.Hourly_rate,
		&analyst.Total_reviews,
		&analyst.Mean_rating,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &analyst, nil
}

// AssignClientRole godoc
// @Summary Atribuir papel de cliente (serviço)
// @Description Atribui o papel de cliente a um usuário existente e retorna o registro do cliente.
// @Tags Serviços
// @Accept json
// @Produce json
// @Param userID path int true "ID do usuário"
// @Param proposedBudget query number true "Orçamento proposto"
// @Success 200 {object} pkg.Client
// @Failure 500 {object} pkg.ErrorResponse
func AssignClientRole(ctx context.Context, conn *pgxpool.Pool, userID int, proposedBudget float64) (*pkg.Client, error) {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var client pkg.Client
	if err := tx.QueryRow(ctx,
		`INSERT INTO client (id, proposed_budget)
		 SELECT id, $2 FROM "user" WHERE id = $1
		 RETURNING id`,
		userID, proposedBudget).Scan(&client.Id); err != nil {
		return nil, err
	}

	if err := tx.QueryRow(ctx,
		`SELECT u.id, u.name, u.email, u.phone, u.time_created,
		       u.country_code, u.country_name, u.country_state, u.city, u.timezone,
		       c.proposed_budget
		FROM "user" u
		JOIN client c ON c.id = u.id
		WHERE u.id = $1`,
		userID).Scan(
		&client.Id,
		&client.Name,
		&client.Email,
		&client.Phone,
		&client.Time_created,
		&client.Country_code,
		&client.Country_name,
		&client.Country_state,
		&client.City,
		&client.Timezone,
		&client.Proposed_budget,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &client, nil
}
