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
// @Param max_total_reviews query int false "Quantidade máxima de avaliações totais"
// @Param min_rating query int false "Avaliação média mínima"
// @Param max_rating query int false "Avaliação média máxima"
// @Param skills   query   string   false   "Nomes de habilidades (ex: c++,java,TypeScript)"
// @Param sort_by query string false "Campo para ordenar: name,country_name,country_state,city,timezone,hourly_rate,total_reviews,mean_rating,time_created"
// @Param order query string false "Direção: ASC ou DESC"
// @Param page query int false "Página"
// @Param page_size query int false "Tamanho da página"
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
		query, args, err := pkg.BuildFilterAnalyst(pkg.BuildFilterUser(builder, filters.UserFilter), filters).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		rows, err := conn.Query(c.Request.Context(), query, args...)
		if pkg.HandleErr(c, err) {
			return
		}

		defer rows.Close()

		analysts, err := pkg.ScanRows(c, rows, pkg.ScanAnalyst)
		if pkg.HandleErr(c, err) {
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

		userFields, analystFields := pkg.BuildUpdateAnalyst(req)

		if len(userFields) == 0 && len(analystFields) == 0 {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "No valid fields were given"))
			return
		}

		tx, err := conn.Begin(c.Request.Context())
		if pkg.HandleErr(c, err) {
			return
		}
		defer func() { _ = tx.Rollback(c.Request.Context()) }()

		if len(userFields) > 0 {
			query, args, err := squirrel.Update(`"user"`).
				SetMap(userFields).Where(squirrel.Eq{"id": id}).
				PlaceholderFormat(squirrel.Dollar).ToSql()
			if pkg.HandleErr(c, err) {
				return
			} else if _, err = tx.Exec(c.Request.Context(), query, args...); pkg.HandleErr(c, err) {
				return
			}
		}

		if len(analystFields) > 0 {
			query, args, err := squirrel.Update("analyst").
				SetMap(analystFields).Where(squirrel.Eq{"id": id}).
				PlaceholderFormat(squirrel.Dollar).ToSql()
			if pkg.HandleErr(c, err) {
				return
			} else if _, err = tx.Exec(c.Request.Context(), query, args...); pkg.HandleErr(c, err) {
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
		} else if err = tx.Commit(c.Request.Context()); pkg.HandleErr(c, err) {
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
// @Param analyst body pkg.AnalystCreateRequest true "Objeto analista (envie apenas `hourly_rate`)"
// @Success 201 {object} pkg.AnalystResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
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
		} else if analystExists {
			c.JSON(http.StatusConflict, pkg.Internal(c.FullPath(), "Analyst already exists"))
			return
		}

		var request pkg.AnalystCreateRequest
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
			SetMap(map[string]any{
				"name":          analyst.Name,
				"email":         analyst.Email,
				"phone":         analyst.Phone,
				"country_code":  analyst.Country_code,
				"country_name":  analyst.Country_name,
				"country_state": analyst.Country_state,
				"city":          analyst.City,
				"timezone":      analyst.Timezone,
			}).Where(squirrel.Eq{"id": id}).PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		} else if _, err = tx.Exec(c.Request.Context(), userQuery, userArgs...); pkg.HandleErr(c, err) {
			return
		}

		analystQuery, analystArgs, err := squirrel.Update("analyst").
			SetMap(map[string]any{
				"hourly_rate":   analyst.Hourly_rate,
				"total_reviews": analyst.Total_reviews,
				"mean_rating":   analyst.Mean_rating,
			}).Where(squirrel.Eq{"id": id}).PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		} else if _, err = tx.Exec(c.Request.Context(), analystQuery, analystArgs...); pkg.HandleErr(c, err) {
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
		} else if err = tx.Commit(c.Request.Context()); pkg.HandleErr(c, err) {
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
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		result, err := conn.Exec(c.Request.Context(), query, args...)
		if pkg.HandleErr(c, err) {
			return
		} else if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Analyst not found"))
			return
		}

		c.Status(http.StatusNoContent)
	}
}
