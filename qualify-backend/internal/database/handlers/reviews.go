package handlers

import (
	"fmt"
	"main/pkg"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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
// @Failure 500 {object} map[string]string
// @Router /reviews [get]
func GetReviews(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Build query with filters
		query := `SELECT id, analyst_id, client_id, service_id, rating, comment, time_created 
                  FROM review WHERE 1=1`

		args := []interface{}{}
		argCounter := 1

		// Filter by analyst_id
		if analystID := c.Query("analyst_id"); analystID != "" {
			if analystIDVal, err := strconv.Atoi(analystID); err == nil {
				query += fmt.Sprintf(" AND analyst_id = $%d", argCounter)
				args = append(args, analystIDVal)
				argCounter++
			}
		}

		// Filter by client_id
		if clientID := c.Query("client_id"); clientID != "" {
			if clientIDVal, err := strconv.Atoi(clientID); err == nil {
				query += fmt.Sprintf(" AND client_id = $%d", argCounter)
				args = append(args, clientIDVal)
				argCounter++
			}
		}

		// Filter by service_id
		if serviceID := c.Query("service_id"); serviceID != "" {
			if serviceIDVal, err := strconv.Atoi(serviceID); err == nil {
				query += fmt.Sprintf(" AND service_id = $%d", argCounter)
				args = append(args, serviceIDVal)
				argCounter++
			}
		}

		// Rating filter (exact match)
		if rating := c.Query("rating"); rating != "" {
			if ratingVal, err := strconv.Atoi(rating); err == nil && ratingVal >= 1 && ratingVal <= 5 {
				query += fmt.Sprintf(" AND rating = $%d", argCounter)
				args = append(args, ratingVal)
				argCounter++
			}
		}

		// Rating range filter
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

		// Comment filter (partial match)
		if comment := c.Query("comment"); comment != "" {
			query += fmt.Sprintf(" AND comment ILIKE $%d", argCounter)
			args = append(args, "%"+comment+"%")
			argCounter++
		}

		// Date range filters
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

		// Optional: Add sorting
		if sortBy := c.Query("sort_by"); sortBy != "" {
			// Validando sortBy para evitar SQL injection - apenas campos permitidos
			allowedSortFields := map[string]bool{
				"id": true, "analyst_id": true, "client_id": true, "service_id": true,
				"rating": true, "time_created": true,
			}
			if allowedSortFields[sortBy] {
				order := c.DefaultQuery("order", "DESC") // Default DESC for reviews (newest first)
				if order == "ASC" || order == "DESC" {
					query += fmt.Sprintf(" ORDER BY %s %s", sortBy, order)
				}
			}
		} else {
			// Default sorting by newest first
			query += " ORDER BY time_created DESC"
		}

		// Pagination
		page := 1
		if p := c.Query("page"); p != "" {
			if pVal, err := strconv.Atoi(p); err == nil && pVal > 0 {
				page = pVal
			}
		}

		pageSize := 20 // Default page size for reviews
		if ps := c.Query("page_size"); ps != "" {
			if psVal, err := strconv.Atoi(ps); err == nil && psVal > 0 && psVal <= 100 {
				pageSize = psVal
			}
		}

		offset := (page - 1) * pageSize
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCounter, argCounter+1)
		args = append(args, pageSize, offset)

		// Execute query
		rows, err := conn.Query(c.Request.Context(), query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching reviews: " + err.Error()})
			return
		}
		defer rows.Close()

		var reviews []pkg.Review

		// Iterate through results
		for rows.Next() {
			var review pkg.Review
			err := rows.Scan(
				&review.Id,
				&review.Analyst_id,
				&review.Client_id,
				&review.Service_id,
				&review.Rating,
				&review.Comment,
				&review.Time_created,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning review data: " + err.Error()})
				return
			}
			reviews = append(reviews, review)
		}

		// Check for errors from iterating over rows
		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating reviews: " + err.Error()})
			return
		}

		// Return results
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
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reviews/{id} [get]
func GetReview(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		reviewID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review id"})
			return
		}
		var review pkg.Review
		err = conn.QueryRow(c.Request.Context(), "SELECT id, analyst_id, client_id, service_id, rating, comment, time_created FROM review WHERE id = $1", reviewID).Scan(&review.Id, &review.Analyst_id, &review.Client_id, &review.Service_id, &review.Rating, &review.Comment, &review.Time_created)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security     BearerAuth
// @Router /reviews [post]
func CreateReview(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var review pkg.Review
		if err := c.BindJSON(&review); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Validando que a avaliação esteja entre 1-5
		if review.Rating < 1 || review.Rating > 5 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Rating must be between 1 and 5"})
			return
		}

		err := conn.QueryRow(c.Request.Context(),
			`INSERT INTO review (analyst_id, client_id, service_id, rating, comment)
			 VALUES ($1, $2, $3, $4, $5)
			 RETURNING id, time_created`,
			review.Analyst_id, review.Client_id, review.Service_id, review.Rating, review.Comment).
			Scan(&review.Id, &review.Time_created)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security     BearerAuth
// @Router /reviews/{id} [put]
func UpdateReview(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		reviewID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review id"})
			return
		}
		var review pkg.Review
		if err := c.BindJSON(&review); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Validando que a avaliação esteja entre 1-5
		if review.Rating < 1 || review.Rating > 5 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Rating must be between 1 and 5"})
			return
		}

		err = conn.QueryRow(c.Request.Context(),
			`UPDATE review SET analyst_id = $1, client_id = $2, service_id = $3, rating = $4, comment = $5
			 WHERE id = $6
			 RETURNING id, analyst_id, client_id, service_id, rating, comment, time_created`,
			review.Analyst_id, review.Client_id, review.Service_id, review.Rating, review.Comment, reviewID).
			Scan(&review.Id, &review.Analyst_id, &review.Client_id, &review.Service_id, &review.Rating, &review.Comment, &review.Time_created)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security     BearerAuth
// @Router /reviews/{id} [patch]
func UpdateReviewPartial(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		reviewID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review id"})
			return
		}
		var review pkg.ReviewUpdateRequest
		if err := c.BindJSON(&review); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
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
				c.JSON(http.StatusBadRequest, gin.H{"error": "Rating must be between 1 and 5"})
				return
			}
			set = append(set, fmt.Sprintf("rating = $%d", i))
			args = append(args, *review.Rating)
			i++
		}

		if len(set) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
			return
		}

		args = append(args, reviewID)

		query := fmt.Sprintf("UPDATE review SET %s WHERE id = $%d RETURNING id, analyst_id, client_id, service_id, rating, comment, time_created",
			strings.Join(set, ", "), i)

		var updatedReview pkg.Review
		err = conn.QueryRow(c.Request.Context(), query, args...).Scan(&updatedReview.Id, &updatedReview.Analyst_id, &updatedReview.Client_id, &updatedReview.Service_id, &updatedReview.Rating, &updatedReview.Comment, &updatedReview.Time_created)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security     BearerAuth
// @Router /reviews/{id} [delete]
func DeleteReview(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		reviewID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review id"})
			return
		}

		result, err := conn.Exec(c.Request.Context(), `DELETE FROM review WHERE id = $1`, reviewID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "review deleted successfully"})
	}
}
