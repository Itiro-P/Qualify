package handlers

import (
	"fmt"
	"main/internal/database/services"
	"main/pkg"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GetAnalysts godoc
// @Summary Listar analistas
// @Description Retorna uma lista de analistas com filtros opcionais
// @Tags Analistas
// @Accept json
// @Produce json
// @Param name query string false "Nome parcial para busca"
// @Param email query string false "Email parcial para busca"
// @Param country query string false "País"
// @Param country_code query string false "Código do país"
// @Param country_state query string false "Estado"
// @Param city query string false "Cidade"
// @Param timezone query string false "Fuso horário"
// @Param min_hourly_rate query number false "Valor mínimo por hora"
// @Param max_hourly_rate query number false "Valor máximo por hora"
// @Param min_total_reviews query int false "Quantidade mínima de avaliações totais"
// @Param min_mean_rating query number false "Avaliação média mínima"
// @Param sort_by query string false "Campo para ordenar: name,country_name,country_state,city,timezone,hourly_rate,total_reviews,mean_rating,time_created"
// @Param order query string false "Direção: ASC ou DESC"
// @Success 200 {object} pkg.AnalystsResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /analysts [get]
func GetAnalysts(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `
			SELECT u.id, u.name, u.email, u.phone, u.time_created,
			       u.country_code, u.country_name, u.country_state, u.city, u.timezone,
			       a.hourly_rate, a.total_reviews, a.mean_rating
			FROM "user" u
			JOIN analyst a ON a.id = u.id
			WHERE 1=1`

		args := []interface{}{}
		argCounter := 1

		if name := c.Query("name"); name != "" {
			query += fmt.Sprintf(" AND u.name ILIKE $%d", argCounter)
			args = append(args, "%"+name+"%")
			argCounter++
		}

		if email := c.Query("email"); email != "" {
			query += fmt.Sprintf(" AND u.email ILIKE $%d", argCounter)
			args = append(args, "%"+email+"%")
			argCounter++
		}

		if country := c.Query("country"); country != "" {
			query += fmt.Sprintf(" AND u.country_name ILIKE $%d", argCounter)
			args = append(args, "%"+country+"%")
			argCounter++
		}

		if countryCode := c.Query("country_code"); countryCode != "" {
			query += fmt.Sprintf(" AND u.country_code ILIKE $%d", argCounter)
			args = append(args, "%"+countryCode+"%")
			argCounter++
		}

		if countryState := c.Query("country_state"); countryState != "" {
			query += fmt.Sprintf(" AND u.country_state ILIKE $%d", argCounter)
			args = append(args, "%"+countryState+"%")
			argCounter++
		}

		if city := c.Query("city"); city != "" {
			query += fmt.Sprintf(" AND u.city ILIKE $%d", argCounter)
			args = append(args, "%"+city+"%")
			argCounter++
		}

		if timezone := c.Query("timezone"); timezone != "" {
			query += fmt.Sprintf(" AND u.timezone ILIKE $%d", argCounter)
			args = append(args, "%"+timezone+"%")
			argCounter++
		}

		if minRate := c.Query("min_hourly_rate"); minRate != "" {
			if minRateVal, err := strconv.ParseFloat(minRate, 64); err == nil {
				query += fmt.Sprintf(" AND a.hourly_rate >= $%d", argCounter)
				args = append(args, minRateVal)
				argCounter++
			}
		}

		if maxRate := c.Query("max_hourly_rate"); maxRate != "" {
			if maxRateVal, err := strconv.ParseFloat(maxRate, 64); err == nil {
				query += fmt.Sprintf(" AND a.hourly_rate <= $%d", argCounter)
				args = append(args, maxRateVal)
				argCounter++
			}
		}

		if minReviews := c.Query("min_total_reviews"); minReviews != "" {
			if minReviewsVal, err := strconv.Atoi(minReviews); err == nil {
				query += fmt.Sprintf(" AND a.total_reviews >= $%d", argCounter)
				args = append(args, minReviewsVal)
				argCounter++
			}
		}

		if minRating := c.Query("min_mean_rating"); minRating != "" {
			if minRatingVal, err := strconv.ParseFloat(minRating, 64); err == nil {
				query += fmt.Sprintf(" AND a.mean_rating >= $%d", argCounter)
				args = append(args, minRatingVal)
				argCounter++
			}
		}

		allowedSortFields := map[string]bool{
			"name": true, "country_name": true, "city": true,
			"hourly_rate": true, "total_reviews": true, "mean_rating": true,
			"time_created": true,
		}
		if sortBy := c.Query("sort_by"); sortBy != "" {
			if allowedSortFields[sortBy] {
				order := c.DefaultQuery("order", "ASC")
				if order == "ASC" || order == "DESC" {
					query += fmt.Sprintf(" ORDER BY %s %s", sortBy, order)
				}
			}
		} else {
			query += " ORDER BY u.time_created DESC"
		}

		rows, err := conn.Query(c.Request.Context(), query, args...)
		if pkg.HandleErr(c, err) {
			return
		}
		defer rows.Close()

		var analysts []pkg.Analyst
		for rows.Next() {
			analyst, err := pkg.ScanAnalyst(rows)
			if pkg.HandleErr(c, err) {
				return
			}
			analysts = append(analysts, analyst)
		}

		if err = rows.Err(); pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusOK, pkg.AnalystsResponse{Analysts: analysts, Count: len(analysts)})
	}
}

// GetAnalyst godoc
// @Summary Obter analista
// @Description Retorna os detalhes de um analista pelo ID do usuário
// @Tags Analistas
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário"
// @Success 200 {object} pkg.AnalystResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /users/{id}/analyst [get]
func GetAnalyst(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		analyst, err := pkg.ScanAnalyst(
			conn.QueryRow(c.Request.Context(),
				`SELECT u.id, u.name, u.email, u.phone, u.time_created,
			       u.country_code, u.country_name, u.country_state, u.city, u.timezone,
			       a.hourly_rate, a.total_reviews, a.mean_rating
				FROM "user" u
				JOIN analyst a ON a.id = u.id
				WHERE u.id = $1`, id))

		if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusOK, pkg.AnalystResponse{Analyst: analyst})
	}
}

// UpdateAnalystPartial godoc
// @Summary Atualizar parcialmente um analista
// @Description Atualiza um ou mais campos do usuário/analista
// @Tags Analistas
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário"
// @Param analyst body pkg.AnalystUpdateRequest true "Campos opcionais para atualização"
// @Success 200 {object} pkg.AnalystResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /users/{id}/analyst [patch]
func UpdateAnalystPartial(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		var req pkg.AnalystUpdateRequest
		if err := c.BindJSON(&req); pkg.HandleErr(c, err) {
			return
		}

		userSet := []string{}
		userArgs := []interface{}{}
		i := 1
		if req.Name != nil {
			userSet = append(userSet, fmt.Sprintf("name = $%d", i))
			userArgs = append(userArgs, *req.Name)
			i++
		}
		if req.Email != nil {
			userSet = append(userSet, fmt.Sprintf("email = $%d", i))
			userArgs = append(userArgs, *req.Email)
			i++
		}
		if req.Phone != nil {
			userSet = append(userSet, fmt.Sprintf("phone = $%d", i))
			userArgs = append(userArgs, *req.Phone)
			i++
		}
		if req.Country_code != nil {
			userSet = append(userSet, fmt.Sprintf("country_code = $%d", i))
			userArgs = append(userArgs, *req.Country_code)
			i++
		}
		if req.Country_name != nil {
			userSet = append(userSet, fmt.Sprintf("country_name = $%d", i))
			userArgs = append(userArgs, *req.Country_name)
			i++
		}
		if req.Country_state != nil {
			userSet = append(userSet, fmt.Sprintf("country_state = $%d", i))
			userArgs = append(userArgs, *req.Country_state)
			i++
		}
		if req.City != nil {
			userSet = append(userSet, fmt.Sprintf("city = $%d", i))
			userArgs = append(userArgs, *req.City)
			i++
		}
		if req.Timezone != nil {
			userSet = append(userSet, fmt.Sprintf("timezone = $%d", i))
			userArgs = append(userArgs, *req.Timezone)
			i++
		}

		analystSet := []string{}
		analystArgs := []interface{}{}
		if req.Hourly_rate != nil {
			analystSet = append(analystSet, fmt.Sprintf("hourly_rate = $%d", i))
			analystArgs = append(analystArgs, *req.Hourly_rate)
			i++
		}
		if req.Total_reviews != nil {
			analystSet = append(analystSet, fmt.Sprintf("total_reviews = $%d", i))
			analystArgs = append(analystArgs, *req.Total_reviews)
			i++
		}
		if req.Mean_rating != nil {
			analystSet = append(analystSet, fmt.Sprintf("mean_rating = $%d", i))
			analystArgs = append(analystArgs, *req.Mean_rating)
			i++
		}

		if len(userSet) == 0 && len(analystSet) == 0 {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "No valid fields were given"))
			return
		}

		tx, err := conn.Begin(c.Request.Context())

		if pkg.HandleErr(c, err) {
			return
		}

		defer func() { _ = tx.Rollback(c.Request.Context()) }()

		// Update user if needed
		if len(userSet) > 0 {
			userArgs = append(userArgs, id)
			userQuery := fmt.Sprintf(`UPDATE "user" SET %s WHERE id = $%d`, strings.Join(userSet, ", "), len(userArgs))
			if _, err := tx.Exec(c.Request.Context(), userQuery, userArgs...); pkg.HandleErr(c, err) {
				return
			}
		}

		// Update analyst if needed
		if len(analystSet) > 0 {
			analystArgs = append(userArgs, analystArgs...)
			// analystArgs currently has userArgs followed by analystArgs; ensure id param at end
			analystArgs = append(analystArgs, id)
			analystQuery := fmt.Sprintf(`UPDATE analyst SET %s WHERE id = $%d`, strings.Join(analystSet, ", "), len(analystArgs))
			if _, err := tx.Exec(c.Request.Context(), analystQuery, analystArgs...); pkg.HandleErr(c, err) {
				return
			}
		}

		analyst, errScan := pkg.ScanAnalyst(tx.QueryRow(c.Request.Context(), `
			SELECT u.id, u.name, u.email, u.phone, u.time_created,
				   u.country_code, u.country_name, u.country_state, u.city, u.timezone,
				   a.hourly_rate, a.total_reviews, a.mean_rating
			FROM "user" u
			JOIN analyst a ON a.id = u.id
			WHERE u.id = $1`, id))

		if pkg.HandleErr(c, errScan) {
			return
		}

		if err = tx.Commit(c.Request.Context()); pkg.HandleErr(c, errScan) {
			return
		}

		c.JSON(http.StatusOK, pkg.AnalystResponse{Analyst: analyst})
	}
}

// CreateAnalyst godoc
// @Summary Criar papel de analista
// @Description Atribui o papel de analista a um usuário existente
// @Tags Analistas
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário"
// @Param analyst body pkg.Analyst true "Objeto analista (envie apenas `hourly_rate`)"
// @Success 201 {object} pkg.AnalystResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 409 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /users/{id}/analyst [post]
func CreateAnalyst(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		// Checando se o analista já existe
		var analystExists bool
		err = conn.QueryRow(c.Request.Context(),
			`SELECT EXISTS(SELECT 1 FROM analyst WHERE id = $1)`, id).Scan(&analystExists)

		if analystExists {
			c.JSON(http.StatusConflict, pkg.Internal(c.FullPath(), "Analyst already exists"))
			return
		}

		var request struct {
			Hourly_rate float64 `json:"hourly_rate"`
		}

		if err := c.BindJSON(&request); pkg.HandleErr(c, err) {
			return
		}

		analyst, err := services.AssignAnalystRole(c.Request.Context(), conn, id, request.Hourly_rate)
		if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusCreated, pkg.AnalystResponse{Analyst: *analyst})
	}
}

// UpdateAnalyst godoc
// @Summary Atualizar analista
// @Description Atualiza dados do analista pelo ID do usuário
// @Tags Analistas
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário"
// @Param analyst body pkg.Analyst true "Objeto analista"
// @Success 200 {object} pkg.AnalystResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /users/{id}/analyst [put]
func UpdateAnalyst(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		var analyst pkg.Analyst
		if err := c.BindJSON(&analyst); pkg.HandleErr(c, err) {
			return
		}

		// Validando parâmetros obrigatórios
		if analyst.Name == "" {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Received empty name"))
			return
		}
		if analyst.Email == "" {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Received empty email"))
			return
		}
		if len(analyst.Country_code) != 2 {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Country code needs to have ONLY 2 characters"))
			return
		}

		// Usar transação para garantir atomicidade entre updates nas tabelas
		tx, err := conn.Begin(c.Request.Context())

		if pkg.HandleErr(c, err) {
			return
		}

		defer func() { _ = tx.Rollback(c.Request.Context()) }()

		_, err = tx.Exec(c.Request.Context(),
			`UPDATE "user" SET name = $1, email = $2, phone = $3, country_code = $4,
			 country_name = $5, country_state = $6, city = $7, timezone = $8
			 WHERE id = $9`,
			analyst.Name, analyst.Email, analyst.Phone, analyst.Country_code, analyst.Country_name, analyst.Country_state,
			analyst.City, analyst.Timezone, id)

		if pkg.HandleErr(c, err) {
			return
		}

		_, err = tx.Exec(c.Request.Context(),
			`UPDATE analyst SET hourly_rate = $1, total_reviews = $2, mean_rating = $3
			 WHERE id = $4`,
			analyst.Hourly_rate, analyst.Total_reviews, analyst.Mean_rating, id)

		if pkg.HandleErr(c, err) {
			return
		}

		analyst, err = pkg.ScanAnalyst(tx.QueryRow(c.Request.Context(),
			`SELECT u.id, u.name, u.email, u.phone, u.time_created,
				   u.country_code, u.country_name, u.country_state, u.city, u.timezone,
				   a.hourly_rate, a.total_reviews, a.mean_rating
			FROM "user" u
			JOIN analyst a ON a.id = u.id
			WHERE u.id = $1`, id))

		if pkg.HandleErr(c, err) {
			return
		}

		if err = tx.Commit(c.Request.Context()); pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusOK, pkg.AnalystResponse{Analyst: analyst})
	}
}

// DeleteAnalyst godoc
// @Summary Excluir analista
// @Description Remove o papel de analista de um usuário
// @Tags Analistas
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário"
// @Success 204 "Deleção com sucesso"
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /users/{id}/analyst [delete]
func DeleteAnalyst(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		analystID, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		result, err := conn.Exec(c.Request.Context(), `DELETE FROM analyst WHERE id = $1`, analystID)

		if pkg.HandleErr(c, err) {
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Analyst not found"))
			return
		}

		c.Status(http.StatusNoContent)
	}
}
