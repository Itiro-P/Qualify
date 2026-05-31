package handlers

import (
	"main/internal/database/services"
	"main/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/n-r-w/squirrel"
)

const analystSelect = `u.id, u.name, u.email, u.phone, u.time_created,
	u.country_code, u.country_name, u.country_state, u.city, u.timezone,
	a.hourly_rate, a.total_reviews, a.mean_rating`

const analystJoin = `"user" u JOIN analyst a ON a.id = u.id`

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
		builder := squirrel.Select(analystSelect).From(analystJoin).PlaceholderFormat(squirrel.Dollar)
		var filters pkg.AnalystFilter
		err := c.ShouldBindQuery(&filters)
		if pkg.HandleErr(c, err) {
			return
		}

		filters.Normalize()

		if filters.Name != "" {
			builder = builder.Where(squirrel.ILike{"u.name": pkg.PutPercent(filters.Name)})
		}

		if filters.Email != "" {
			builder = builder.Where(squirrel.ILike{"u.email": pkg.PutPercent(filters.Email)})
		}

		if filters.Country != "" {
			builder = builder.Where(squirrel.ILike{"u.country_name": pkg.PutPercent(filters.Country)})
		}

		if filters.CountryCode != "" {
			builder = builder.Where(squirrel.ILike{"u.country_code": pkg.PutPercent(filters.CountryCode)})
		}

		if filters.CountryState != "" {
			builder = builder.Where(squirrel.ILike{"u.country_state": pkg.PutPercent(filters.CountryState)})
		}

		if filters.City != "" {
			builder = builder.Where(squirrel.ILike{"u.city": pkg.PutPercent(filters.City)})
		}

		if filters.Timezone != "" {
			builder = builder.Where(squirrel.ILike{"u.timezone": pkg.PutPercent(filters.Timezone)})
		}

		if filters.MinHourlyRate != nil {
			builder = builder.Where(squirrel.GtOrEq{"a.hourly_rate": *filters.MinHourlyRate})
		}

		if filters.MaxHourlyRate != nil {
			builder = builder.Where(squirrel.LtOrEq{"a.hourly_rate": *filters.MaxHourlyRate})
		}

		if filters.MinTotalReviews != nil {
			builder = builder.Where(squirrel.GtOrEq{"a.total_reviews": *filters.MinTotalReviews})
		}

		if filters.MaxTotalReviews != nil {
			builder = builder.Where(squirrel.LtOrEq{"a.total_reviews": *filters.MaxTotalReviews})
		}

		if filters.MinRating != nil {
			builder = builder.Where(squirrel.GtOrEq{"a.mean_rating": *filters.MinRating})
		}

		if filters.MaxRating != nil {
			builder = builder.Where(squirrel.LtOrEq{"a.mean_rating": *filters.MaxRating})
		}

		orderClause := filters.SortOptions.ValidateSort(pkg.AnalystSortFields)

		if orderClause != "" {
			builder = builder.OrderBy(orderClause)
		} else {
			builder = builder.OrderBy("time_created DESC")
		}

		builder = builder.Limit(uint64(filters.PageSize)).Offset(uint64(filters.Offset()))

		query, args, err := builder.ToSql()
		if pkg.HandleErr(c, err) {
			return
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

		query, args, err := squirrel.Select(analystSelect).
			From(analystJoin).Where(squirrel.Eq{"u.id": id}).
			PlaceholderFormat(squirrel.Dollar).ToSql()

		if pkg.HandleErr(c, err) {
			return
		}

		analyst, err := pkg.ScanAnalyst(conn.QueryRow(c.Request.Context(), query, args...))

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

		userBuilder := squirrel.Update(`"user"`).PlaceholderFormat(squirrel.Dollar)
		analystBuilder := squirrel.Update("analyst").PlaceholderFormat(squirrel.Dollar)
		userHasFields := false
		analystHasFields := false

		if req.Name != nil {
			userBuilder = userBuilder.Set("name", *req.Name)
			userHasFields = true
		}
		if req.Email != nil {
			userBuilder = userBuilder.Set("email", *req.Email)
			userHasFields = true
		}
		if req.Phone != nil {
			userBuilder = userBuilder.Set("phone", *req.Phone)
			userHasFields = true
		}
		if req.Country_code != nil {
			userBuilder = userBuilder.Set("country_code", *req.Country_code)
			userHasFields = true
		}
		if req.Country_name != nil {
			userBuilder = userBuilder.Set("country_name", *req.Country_name)
			userHasFields = true
		}
		if req.Country_state != nil {
			userBuilder = userBuilder.Set("country_state", *req.Country_state)
			userHasFields = true
		}
		if req.City != nil {
			userBuilder = userBuilder.Set("city", *req.City)
			userHasFields = true
		}
		if req.Timezone != nil {
			userBuilder = userBuilder.Set("timezone", *req.Timezone)
			userHasFields = true
		}

		if req.Hourly_rate != nil {
			analystBuilder = analystBuilder.Set("hourly_rate", *req.Hourly_rate)
			analystHasFields = true
		}
		if req.Total_reviews != nil {
			analystBuilder = analystBuilder.Set("total_reviews", *req.Total_reviews)
			analystHasFields = true
		}
		if req.Mean_rating != nil {
			analystBuilder = analystBuilder.Set("mean_rating", *req.Mean_rating)
			analystHasFields = true
		}

		if !userHasFields && !analystHasFields {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "No valid fields were given"))
			return
		}

		tx, err := conn.Begin(c.Request.Context())
		if pkg.HandleErr(c, err) {
			return
		}
		defer tx.Rollback(c.Request.Context())

		if userHasFields {
			query, args, err := userBuilder.Where(squirrel.Eq{"id": id}).ToSql()
			if pkg.HandleErr(c, err) {
				return
			}
			if _, err = tx.Exec(c.Request.Context(), query, args...); pkg.HandleErr(c, err) {
				return
			}
		}

		if analystHasFields {
			query, args, err := analystBuilder.Where(squirrel.Eq{"id": id}).ToSql()
			if pkg.HandleErr(c, err) {
				return
			}
			if _, err = tx.Exec(c.Request.Context(), query, args...); pkg.HandleErr(c, err) {
				return
			}
		}

		selectQuery, selectArgs, err := squirrel.Select(analystSelect).
			From(analystJoin).Where(squirrel.Eq{"u.id": id}).
			PlaceholderFormat(squirrel.Dollar).ToSql()

		if pkg.HandleErr(c, err) {
			return
		}

		analyst, err := pkg.ScanAnalyst(tx.QueryRow(c.Request.Context(), selectQuery, selectArgs...))
		if pkg.HandleErr(c, err) {
			return
		}

		if err = tx.Commit(c.Request.Context()); pkg.HandleErr(c, err) {
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

		var analystExists bool
		err = conn.QueryRow(c.Request.Context(),
			`SELECT EXISTS(SELECT 1 FROM analyst WHERE id = $1)`, id).Scan(&analystExists)

		if pkg.HandleErr(c, err) {
			return
		}

		if analystExists {
			c.JSON(http.StatusConflict, pkg.Conflict(c.FullPath(), "Analyst already exists"))
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

		if analyst.Name == "" {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Received empty name"))
			return
		}
		if analyst.Email == "" {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Received empty email"))
			return
		}
		if len(analyst.Country_code) != 2 {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Country code must be exactly 2 characters"))
			return
		}

		tx, err := conn.Begin(c.Request.Context())
		if pkg.HandleErr(c, err) {
			return
		}
		defer tx.Rollback(c.Request.Context())

		userQuery, userArgs, err := squirrel.Update(`"user"`).
			Set("name", analyst.Name).
			Set("email", analyst.Email).
			Set("phone", analyst.Phone).
			Set("country_code", analyst.Country_code).
			Set("country_name", analyst.Country_name).
			Set("country_state", analyst.Country_state).
			Set("city", analyst.City).
			Set("timezone", analyst.Timezone).
			Where(squirrel.Eq{"id": id}).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if pkg.HandleErr(c, err) {
			return
		}
		if _, err = tx.Exec(c.Request.Context(), userQuery, userArgs...); pkg.HandleErr(c, err) {
			return
		}

		analystQuery, analystArgs, err := squirrel.Update("analyst").
			Set("hourly_rate", analyst.Hourly_rate).
			Set("total_reviews", analyst.Total_reviews).
			Set("mean_rating", analyst.Mean_rating).
			Where(squirrel.Eq{"id": id}).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if pkg.HandleErr(c, err) {
			return
		}
		if _, err = tx.Exec(c.Request.Context(), analystQuery, analystArgs...); pkg.HandleErr(c, err) {
			return
		}

		selectQuery, selectArgs, err := squirrel.Select(analystSelect).
			From(analystJoin).Where(squirrel.Eq{"u.id": id}).
			PlaceholderFormat(squirrel.Dollar).ToSql()

		if pkg.HandleErr(c, err) {
			return
		}

		analyst, err = pkg.ScanAnalyst(tx.QueryRow(c.Request.Context(), selectQuery, selectArgs...))
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
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		query, args, err := squirrel.Delete("analyst").
			Where(squirrel.Eq{"id": id}).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		result, err := conn.Exec(c.Request.Context(), query, args...)
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
