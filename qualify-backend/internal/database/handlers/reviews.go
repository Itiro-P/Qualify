package handlers

import (
	"errors"
	"fmt"
	"main/pkg"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

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
		query := `SELECT id, analyst_id, client_id, service_id, rating, comment, time_created
                  FROM review WHERE 1=1`

		args := []interface{}{}
		argCounter := 1

		if analystID := c.Query("analyst_id"); analystID != "" {
			if analystIDVal, err := strconv.Atoi(analystID); err == nil {
				query += fmt.Sprintf(" AND analyst_id = $%d", argCounter)
				args = append(args, analystIDVal)
				argCounter++
			}
		}

		if clientID := c.Query("client_id"); clientID != "" {
			if clientIDVal, err := strconv.Atoi(clientID); err == nil {
				query += fmt.Sprintf(" AND client_id = $%d", argCounter)
				args = append(args, clientIDVal)
				argCounter++
			}
		}

		if serviceID := c.Query("service_id"); serviceID != "" {
			if serviceIDVal, err := strconv.Atoi(serviceID); err == nil {
				query += fmt.Sprintf(" AND service_id = $%d", argCounter)
				args = append(args, serviceIDVal)
				argCounter++
			}
		}

		if rating := c.Query("rating"); rating != "" {
			if ratingVal, err := strconv.Atoi(rating); err == nil && ratingVal >= 1 && ratingVal <= 5 {
				query += fmt.Sprintf(" AND rating = $%d", argCounter)
				args = append(args, ratingVal)
				argCounter++
			}
		}

		if minRating := c.Query("min_rating"); minRating != "" {
			if minRatingVal, err := strconv.Atoi(minRating); err == nil {
				query += fmt.Sprintf(" AND rating >= $%d", argCounter)
				args = append(args, minRatingVal)
				argCounter++
			}
		}

		if maxRating := c.Query("max_rating"); maxRating != "" {
			if maxRatingVal, err := strconv.Atoi(maxRating); err == nil {
				query += fmt.Sprintf(" AND rating <= $%d", argCounter)
				args = append(args, maxRatingVal)
				argCounter++
			}
		}

		if comment := c.Query("comment"); comment != "" {
			query += fmt.Sprintf(" AND comment ILIKE $%d", argCounter)
			args = append(args, "%"+comment+"%")
			argCounter++
		}

		if fromDate := c.Query("from_date"); fromDate != "" {
			query += fmt.Sprintf(" AND time_created >= $%d", argCounter)
			args = append(args, fromDate)
			argCounter++
		}

		if toDate := c.Query("to_date"); toDate != "" {
			query += fmt.Sprintf(" AND time_created <= $%d", argCounter)
			args = append(args, toDate)
			argCounter++
		}

		if sortBy := c.Query("sort_by"); sortBy != "" {
			allowedSortFields := map[string]bool{
				"id": true, "analyst_id": true, "client_id": true, "service_id": true,
				"rating": true, "time_created": true,
			}
			if allowedSortFields[sortBy] {
				order := c.DefaultQuery("order", "DESC")
				if order == "ASC" || order == "DESC" {
					query += fmt.Sprintf(" ORDER BY %s %s", sortBy, order)
				}
			}
		} else {
			query += " ORDER BY time_created DESC"
		}

		page := 1
		if p := c.Query("page"); p != "" {
			if pVal, err := strconv.Atoi(p); err == nil && pVal > 0 {
				page = pVal
			}
		}

		pageSize := 20
		if ps := c.Query("page_size"); ps != "" {
			if psVal, err := strconv.Atoi(ps); err == nil && psVal > 0 && psVal <= 100 {
				pageSize = psVal
			}
		}

		offset := (page - 1) * pageSize
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCounter, argCounter+1)
		args = append(args, pageSize, offset)

		rows, err := conn.Query(c.Request.Context(), query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}
		defer rows.Close()

		var reviews []pkg.Review

		for rows.Next() {
			review, err := pkg.ScanReview(rows)
			if err != nil {
				c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
				return
			}
			reviews = append(reviews, review)
		}

		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}

		c.JSON(http.StatusOK, pkg.ReviewsResponse{
			Reviews:   reviews,
			Count:     len(reviews),
			Page:      page,
			Page_size: pageSize,
		})
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

		row := conn.QueryRow(c.Request.Context(), "SELECT id, analyst_id, client_id, service_id, rating, comment, time_created FROM review WHERE id = $1", id)
		review, err := pkg.ScanReview(row)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), err.Error()))
				return
			}
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
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
		if err := c.BindJSON(&review); err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		if review.Rating < 1 || review.Rating > 5 {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Rating should be between 1 and 5"))
			return
		}

		review, err := pkg.ScanReview(conn.QueryRow(c.Request.Context(),
			`INSERT INTO review (analyst_id, client_id, service_id, rating, comment)
             VALUES ($1, $2, $3, $4, $5)
             RETURNING id, analyst_id, client_id, service_id, rating, comment, time_created`,
			review.Analyst_id, review.Client_id, review.Service_id, review.Rating, review.Comment))
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				c.JSON(http.StatusConflict, pkg.Conflict(c.FullPath(), "Review already exists"))
				return
			}
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
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

		var review pkg.Review
		if err := c.BindJSON(&review); err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		if review.Rating < 1 || review.Rating > 5 {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Rating should be between 1 and 5"))
			return
		}

		review, err = pkg.ScanReview(conn.QueryRow(c.Request.Context(),
			`UPDATE review SET analyst_id = $1, client_id = $2, service_id = $3, rating = $4, comment = $5
             WHERE id = $6
             RETURNING id, analyst_id, client_id, service_id, rating, comment, time_created`,
			review.Analyst_id, review.Client_id, review.Service_id, review.Rating, review.Comment, id))

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), err.Error()))
				return
			}
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
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

		var review pkg.ReviewUpdateRequest
		if err := c.BindJSON(&review); err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		set := []string{}
		args := []interface{}{}
		i := 1

		if review.Comment != nil {
			set = append(set, fmt.Sprintf("comment = $%d", i))
			args = append(args, *review.Comment)
			i++
		}
		if review.Rating != nil {
			if *review.Rating < 1 || *review.Rating > 5 {
				c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Rating should be between 1 and 5"))
				return
			}
			set = append(set, fmt.Sprintf("rating = $%d", i))
			args = append(args, *review.Rating)
			i++
		}

		if len(set) == 0 {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Invalid arguments"))
			return
		}

		args = append(args, id)

		query := fmt.Sprintf("UPDATE review SET %s WHERE id = $%d RETURNING id, analyst_id, client_id, service_id, rating, comment, time_created",
			strings.Join(set, ", "), i)

		updatedReview, err := pkg.ScanReview(conn.QueryRow(c.Request.Context(), query, args...))
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), err.Error()))
				return
			}
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}

		c.JSON(http.StatusOK, pkg.ReviewResponse{Review: updatedReview})
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

		result, err := conn.Exec(c.Request.Context(), `DELETE FROM review WHERE id = $1`, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}

		if result.RowsAffected() == 0 {
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
// @Failure 500 {object} pkg.ErrorResponse
// @Router /clients/{id}/reviews [get]
func GetClientReviews(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		query := `SELECT r.id, r.analyst_id, r.client_id, r.service_id, r.rating, r.comment, r.time_created
                  FROM review r
                  JOIN client c ON r.client_id = c.id
                  WHERE c.user_id = $1`
		args := []interface{}{id}
		argCounter := 2

		if analystID := c.Query("analyst_id"); analystID != "" {
			if analystIDVal, err := strconv.Atoi(analystID); err == nil {
				query += fmt.Sprintf(" AND r.analyst_id = $%d", argCounter)
				args = append(args, analystIDVal)
				argCounter++
			}
		}
		if serviceID := c.Query("service_id"); serviceID != "" {
			if serviceIDVal, err := strconv.Atoi(serviceID); err == nil {
				query += fmt.Sprintf(" AND r.service_id = $%d", argCounter)
				args = append(args, serviceIDVal)
				argCounter++
			}
		}
		if rating := c.Query("rating"); rating != "" {
			if ratingVal, err := strconv.Atoi(rating); err == nil && ratingVal >= 1 && ratingVal <= 5 {
				query += fmt.Sprintf(" AND r.rating = $%d", argCounter)
				args = append(args, ratingVal)
				argCounter++
			}
		}
		if minRating := c.Query("min_rating"); minRating != "" {
			if minRatingVal, err := strconv.Atoi(minRating); err == nil {
				query += fmt.Sprintf(" AND r.rating >= $%d", argCounter)
				args = append(args, minRatingVal)
				argCounter++
			}
		}
		if maxRating := c.Query("max_rating"); maxRating != "" {
			if maxRatingVal, err := strconv.Atoi(maxRating); err == nil {
				query += fmt.Sprintf(" AND r.rating <= $%d", argCounter)
				args = append(args, maxRatingVal)
				argCounter++
			}
		}
		if comment := c.Query("comment"); comment != "" {
			query += fmt.Sprintf(" AND r.comment ILIKE $%d", argCounter)
			args = append(args, "%"+comment+"%")
			argCounter++
		}
		if fromDate := c.Query("from_date"); fromDate != "" {
			query += fmt.Sprintf(" AND r.time_created >= $%d", argCounter)
			args = append(args, fromDate)
			argCounter++
		}
		if toDate := c.Query("to_date"); toDate != "" {
			query += fmt.Sprintf(" AND r.time_created <= $%d", argCounter)
			args = append(args, toDate)
			argCounter++
		}

		allowedSortFields := map[string]bool{
			"id": true, "analyst_id": true, "client_id": true, "service_id": true,
			"rating": true, "time_created": true,
		}
		if sortBy := c.Query("sort_by"); sortBy != "" {
			if allowedSortFields[sortBy] {
				order := c.DefaultQuery("order", "DESC")
				if order == "ASC" || order == "DESC" {
					query += fmt.Sprintf(" ORDER BY r.%s %s", sortBy, order)
				}
			}
		} else {
			query += " ORDER BY r.time_created DESC"
		}

		page := 1
		if p := c.Query("page"); p != "" {
			if pVal, err := strconv.Atoi(p); err == nil && pVal > 0 {
				page = pVal
			}
		}
		pageSize := 20
		if ps := c.Query("page_size"); ps != "" {
			if psVal, err := strconv.Atoi(ps); err == nil && psVal > 0 && psVal <= 100 {
				pageSize = psVal
			}
		}

		offset := (page - 1) * pageSize
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCounter, argCounter+1)
		args = append(args, pageSize, offset)

		rows, err := conn.Query(c.Request.Context(), query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}
		defer rows.Close()

		var reviews []pkg.Review
		for rows.Next() {
			review, err := pkg.ScanReview(rows)
			if err != nil {
				c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
				return
			}
			reviews = append(reviews, review)
		}
		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}

		c.JSON(http.StatusOK, pkg.ReviewsResponse{Reviews: reviews, Count: len(reviews), Page: page, Page_size: pageSize})
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
// @Failure 500 {object} pkg.ErrorResponse
// @Router /analysts/{id}/reviews [get]
func GetAnalystReviews(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		query := `SELECT r.id, r.analyst_id, r.client_id, r.service_id, r.rating, r.comment, r.time_created
                  FROM review r
                  JOIN analyst a ON r.analyst_id = a.id
                  WHERE a.user_id = $1`
		args := []interface{}{id}
		argCounter := 2

		if clientID := c.Query("client_id"); clientID != "" {
			if clientIDVal, err := strconv.Atoi(clientID); err == nil {
				query += fmt.Sprintf(" AND r.client_id = $%d", argCounter)
				args = append(args, clientIDVal)
				argCounter++
			}
		}
		if serviceID := c.Query("service_id"); serviceID != "" {
			if serviceIDVal, err := strconv.Atoi(serviceID); err == nil {
				query += fmt.Sprintf(" AND r.service_id = $%d", argCounter)
				args = append(args, serviceIDVal)
				argCounter++
			}
		}
		if rating := c.Query("rating"); rating != "" {
			if ratingVal, err := strconv.Atoi(rating); err == nil && ratingVal >= 1 && ratingVal <= 5 {
				query += fmt.Sprintf(" AND r.rating = $%d", argCounter)
				args = append(args, ratingVal)
				argCounter++
			}
		}
		if minRating := c.Query("min_rating"); minRating != "" {
			if minRatingVal, err := strconv.Atoi(minRating); err == nil {
				query += fmt.Sprintf(" AND r.rating >= $%d", argCounter)
				args = append(args, minRatingVal)
				argCounter++
			}
		}
		if maxRating := c.Query("max_rating"); maxRating != "" {
			if maxRatingVal, err := strconv.Atoi(maxRating); err == nil {
				query += fmt.Sprintf(" AND r.rating <= $%d", argCounter)
				args = append(args, maxRatingVal)
				argCounter++
			}
		}
		if comment := c.Query("comment"); comment != "" {
			query += fmt.Sprintf(" AND r.comment ILIKE $%d", argCounter)
			args = append(args, "%"+comment+"%")
			argCounter++
		}
		if fromDate := c.Query("from_date"); fromDate != "" {
			query += fmt.Sprintf(" AND r.time_created >= $%d", argCounter)
			args = append(args, fromDate)
			argCounter++
		}
		if toDate := c.Query("to_date"); toDate != "" {
			query += fmt.Sprintf(" AND r.time_created <= $%d", argCounter)
			args = append(args, toDate)
			argCounter++
		}

		allowedSortFields := map[string]bool{
			"id": true, "analyst_id": true, "client_id": true, "service_id": true,
			"rating": true, "time_created": true,
		}
		if sortBy := c.Query("sort_by"); sortBy != "" {
			if allowedSortFields[sortBy] {
				order := c.DefaultQuery("order", "DESC")
				if order == "ASC" || order == "DESC" {
					query += fmt.Sprintf(" ORDER BY r.%s %s", sortBy, order)
				}
			}
		} else {
			query += " ORDER BY r.time_created DESC"
		}

		page := 1
		if p := c.Query("page"); p != "" {
			if pVal, err := strconv.Atoi(p); err == nil && pVal > 0 {
				page = pVal
			}
		}
		pageSize := 20
		if ps := c.Query("page_size"); ps != "" {
			if psVal, err := strconv.Atoi(ps); err == nil && psVal > 0 && psVal <= 100 {
				pageSize = psVal
			}
		}

		offset := (page - 1) * pageSize
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCounter, argCounter+1)
		args = append(args, pageSize, offset)

		rows, err := conn.Query(c.Request.Context(), query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}
		defer rows.Close()

		var reviews []pkg.Review
		for rows.Next() {
			review, err := pkg.ScanReview(rows)
			if err != nil {
				c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
				return
			}
			reviews = append(reviews, review)
		}
		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}

		c.JSON(http.StatusOK, pkg.ReviewsResponse{Reviews: reviews, Count: len(reviews), Page: page, Page_size: pageSize})
	}
}
