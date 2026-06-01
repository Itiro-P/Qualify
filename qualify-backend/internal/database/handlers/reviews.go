package handlers

import (
	"main/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/n-r-w/squirrel"
)

const reviewSelect = "id, analyst_id, client_id, service_id, rating, comment, time_created"

func applyReviewFilters(builder squirrel.SelectBuilder, filters pkg.ReviewFilter) squirrel.SelectBuilder {
	if filters.AnalystId != nil {
		builder = builder.Where(squirrel.Eq{"analyst_id": *filters.AnalystId})
	}
	if filters.ClientId != nil {
		builder = builder.Where(squirrel.Eq{"client_id": *filters.ClientId})
	}
	if filters.ServiceId != nil {
		builder = builder.Where(squirrel.Eq{"service_id": *filters.ServiceId})
	}
	if filters.Rating != nil {
		builder = builder.Where(squirrel.Eq{"rating": *filters.Rating})
	}
	if filters.MinRating != nil {
		builder = builder.Where(squirrel.GtOrEq{"rating": *filters.MinRating})
	}
	if filters.MaxRating != nil {
		builder = builder.Where(squirrel.LtOrEq{"rating": *filters.MaxRating})
	}
	if filters.Comment != "" {
		builder = builder.Where(squirrel.ILike{"comment": pkg.PutPercent(filters.Comment)})
	}
	if order := filters.ValidateSort(pkg.ReviewSortFields); order != "" {
		builder = builder.OrderBy(order)
	} else {
		builder = builder.OrderBy("time_created DESC")
	}
	return builder.Limit(uint64(filters.PageSize)).Offset(uint64(filters.Offset()))
}

func reviewBuilder(filters pkg.ReviewFilter) squirrel.SelectBuilder {
	return applyReviewFilters(
		squirrel.Select(reviewSelect).From("review").PlaceholderFormat(squirrel.Dollar),
		filters,
	)
}

func scanReviews(c *gin.Context, conn *pgxpool.Pool, builder squirrel.SelectBuilder, filters pkg.ReviewFilter) {
	query, args, err := builder.ToSql()
	if pkg.HandleErr(c, err) {
		return
	}

	rows, err := conn.Query(c.Request.Context(), query, args...)
	if pkg.HandleErr(c, err) {
		return
	}
	defer rows.Close()

	reviews, err := pkg.ScanRows(c, rows, pkg.ScanReview)
	if pkg.HandleErr(c, err) {
		return
	}

	c.JSON(http.StatusOK, pkg.ReviewsResponse{
		Reviews:   reviews,
		Count:     len(reviews),
		Page:      filters.Page,
		Page_size: filters.PageSize,
	})
}

// GetReviews godoc
// @Summary Listar avaliações
// @Description Retorna avaliações com filtros e paginação
// @Tags Avaliações
// @Accept json
// @Produce json
// @Param analyst_id query int false "ID do analista"
// @Param client_id query int false "ID do cliente"
// @Param service_id query int false "ID do serviço"
// @Param rating query int false "Avaliação exata"
// @Param min_rating query int false "Avaliação mínima"
// @Param max_rating query int false "Avaliação máxima"
// @Param page query int false "Página"
// @Param page_size query int false "Tamanho da página"
// @Success 200 {object} pkg.ReviewsResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /reviews [get]
func GetReviews(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filters pkg.ReviewFilter
		err := c.ShouldBindQuery(&filters)
		if pkg.HandleErr(c, err) {
			return
		}
		filters.Normalize()

		scanReviews(c, conn, reviewBuilder(filters), filters)
	}
}

// GetReview godoc
// @Summary Obter avaliação
// @Description Retorna uma avaliação pelo ID
// @Tags Avaliações
// @Accept json
// @Produce json
// @Param id path int true "ID da avaliação"
// @Success 200 {object} pkg.ReviewResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /reviews/{id} [get]

func GetReview(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		query, args, err := squirrel.Select(reviewSelect).
			From("review").Where(squirrel.Eq{"id": id}).
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		review, err := pkg.ScanReview(conn.QueryRow(c.Request.Context(), query, args...))
		if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusOK, pkg.ReviewResponse{Review: review})
	}
}

// CreateReview godoc
// @Summary Criar avaliação
// @Description Cria uma nova avaliação para um serviço
// @Tags Avaliações
// @Accept json
// @Produce json
// @Param review body pkg.Review true "Objeto avaliação"
// @Success 201 {object} pkg.ReviewResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 409 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /reviews [post]
func CreateReview(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var review pkg.Review
		if err := c.BindJSON(&review); pkg.HandleErr(c, err) {
			return
		}

		if review.Rating < 1 || review.Rating > 5 {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Rating should be between 1 and 5"))
			return
		}

		query, args, err := squirrel.Insert("review").
			Columns("analyst_id", "client_id", "service_id", "rating", "comment").
			Values(review.Analyst_id, review.Client_id, review.Service_id, review.Rating, review.Comment).
			Suffix("RETURNING " + reviewSelect).
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		review, err = pkg.ScanReview(conn.QueryRow(c.Request.Context(), query, args...))
		if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusCreated, pkg.ReviewResponse{Review: review})
	}
}

// UpdateReview godoc
// @Summary Atualizar avaliação
// @Description Atualiza uma avaliação pelo ID
// @Tags Avaliações
// @Accept json
// @Produce json
// @Param id path int true "ID da avaliação"
// @Param review body pkg.Review true "Objeto avaliação"
// @Success 200 {object} pkg.ReviewResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /reviews/{id} [put]
func UpdateReview(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		var reviewExists bool
		err = conn.QueryRow(c.Request.Context(),
			`SELECT EXISTS(SELECT 1 FROM review WHERE id = $1)`, id).Scan(&reviewExists)
		if pkg.HandleErr(c, err) {
			return
		} else if !reviewExists {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Review does not exists"))
			return
		}

		var review pkg.Review
		if err := c.BindJSON(&review); pkg.HandleErr(c, err) {
			return
		}

		if review.Rating < 1 || review.Rating > 5 {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Rating should be between 1 and 5"))
			return
		}

		query, args, err := squirrel.Update("review").
			SetMap(map[string]any{
				"analyst_id": review.Analyst_id,
				"client_id":  review.Client_id,
				"service_id": review.Service_id,
				"rating":     review.Rating,
				"comment":    review.Comment}).
			Where(squirrel.Eq{"id": id}).
			Suffix("RETURNING " + reviewSelect).
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		review, err = pkg.ScanReview(conn.QueryRow(c.Request.Context(), query, args...))
		if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusOK, pkg.ReviewResponse{Review: review})
	}
}

// UpdateReviewPartial godoc
// @Summary Atualizar parcialmente avaliação
// @Description Atualiza um ou mais campos da avaliação pelo ID
// @Tags Avaliações
// @Accept json
// @Produce json
// @Param id path int true "ID da avaliação"
// @Param review body pkg.ReviewUpdateRequest true "Objeto avaliação"
// @Success 200 {object} pkg.ReviewResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /reviews/{id} [patch]
func UpdateReviewPartial(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		var reviewExists bool
		err = conn.QueryRow(c.Request.Context(),
			`SELECT EXISTS(SELECT 1 FROM review WHERE id = $1)`, id).Scan(&reviewExists)
		if pkg.HandleErr(c, err) {
			return
		} else if !reviewExists {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Review does not exists"))
			return
		}

		var req pkg.ReviewUpdateRequest
		if err := c.BindJSON(&req); pkg.HandleErr(c, err) {
			return
		}

		builder := squirrel.Update("review").PlaceholderFormat(squirrel.Dollar)
		hasFields := false

		if req.Comment != nil {
			builder = builder.Set("comment", *req.Comment)
			hasFields = true
		}
		if req.Rating != nil {
			if *req.Rating < 1 || *req.Rating > 5 {
				c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Rating should be between 1 and 5"))
				return
			}
			builder = builder.Set("rating", *req.Rating)
			hasFields = true
		}

		if !hasFields {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "No valid fields were given"))
			return
		}

		query, args, err := builder.
			Where(squirrel.Eq{"id": id}).
			Suffix("RETURNING " + reviewSelect).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		review, err := pkg.ScanReview(conn.QueryRow(c.Request.Context(), query, args...))
		if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusOK, pkg.ReviewResponse{Review: review})
	}
}

// DeleteReview godoc
// @Summary Excluir avaliação
// @Description Remove uma avaliação pelo ID
// @Tags Avaliações
// @Accept json
// @Produce json
// @Param id path int true "ID da avaliação"
// @Success 204 "Deleção com sucesso"
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /reviews/{id} [delete]
func DeleteReview(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		query, args, err := squirrel.Delete("review").
			Where(squirrel.Eq{"id": id}).
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		result, err := conn.Exec(c.Request.Context(), query, args...)
		if pkg.HandleErr(c, err) {
			return
		} else if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Review not found"))
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// GetClientReviews godoc
// @Summary Listar avaliações de um cliente
// @Description Retorna avaliações relacionadas a um cliente com filtros e paginação
// @Tags Avaliações
// @Accept json
// @Produce json
// @Param id path int true "ID do cliente"
// @Param analyst_id query int false "ID do analista"
// @Param service_id query int false "ID do serviço"
// @Param rating query int false "Avaliação exata"
// @Param min_rating query int false "Avaliação mínima"
// @Param max_rating query int false "Avaliação máxima"
// @Param page query int false "Página"
// @Param page_size query int false "Tamanho da página"
// @Success 200 {object} pkg.ReviewsResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /clients/{id}/reviews [get]
func GetClientReviews(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		// Checando se o cliente já existe
		var clientExists bool
		err = conn.QueryRow(c.Request.Context(),
			`SELECT EXISTS(SELECT 1 FROM client WHERE id = $1)`, id).Scan(&clientExists)
		if pkg.HandleErr(c, err) {
			return
		} else if !clientExists {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Client does not exists"))
			return
		}

		var filters pkg.ReviewFilter
		err = c.ShouldBindQuery(&filters)
		if pkg.HandleErr(c, err) {
			return
		}
		filters.Normalize()

		builder := applyReviewFilters(
			squirrel.Select(reviewSelect).
				From("review r").
				Join("client cl ON r.client_id = cl.id").
				Where(squirrel.Eq{"cl.id": id}).
				PlaceholderFormat(squirrel.Dollar), filters)

		scanReviews(c, conn, builder, filters)
	}
}

// GetAnalystReviews godoc
// @Summary Listar avaliações de um analista
// @Description Retorna avaliações relacionadas a um analista com filtros e paginação
// @Tags Avaliações
// @Accept json
// @Produce json
// @Param id path int true "ID do analista"
// @Param client_id query int false "ID do cliente"
// @Param service_id query int false "ID do serviço"
// @Param rating query int false "Avaliação exata"
// @Param min_rating query int false "Avaliação mínima"
// @Param max_rating query int false "Avaliação máxima"
// @Param page query int false "Página"
// @Param page_size query int false "Tamanho da página"
// @Success 200 {object} pkg.ReviewsResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /analysts/{id}/reviews [get]
func GetAnalystReviews(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		// Checando se o analista já existe
		var analystExists bool
		err = conn.QueryRow(c.Request.Context(),
			`SELECT EXISTS(SELECT 1 FROM analyst WHERE id = $1)`, id).Scan(&analystExists)
		if pkg.HandleErr(c, err) {
			return
		} else if !analystExists {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Analyst does not exists"))
			return
		}

		var filters pkg.ReviewFilter
		err = c.ShouldBindQuery(&filters)
		if pkg.HandleErr(c, err) {
			return
		}
		filters.Normalize()

		builder := applyReviewFilters(
			squirrel.Select(reviewSelect).
				From("review r").
				Join("analyst a ON r.analyst_id = a.id").
				Where(squirrel.Eq{"a.id": id}).
				PlaceholderFormat(squirrel.Dollar), filters)

		scanReviews(c, conn, builder, filters)
	}
}
