package handlers

import (
	"fmt"
	"main/pkg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func GetReviews(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Build query with filters
		query := `SELECT id, service_id, rating, comment, time_created 
                  FROM review WHERE 1=1`

		args := []interface{}{}
		argCounter := 1

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
			// Validate sortBy to prevent SQL injection
			allowedSortFields := map[string]bool{
				"id": true, "service_id": true,
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
		c.JSON(http.StatusOK, gin.H{
			"reviews":   reviews,
			"count":     len(reviews),
			"page":      page,
			"page_size": pageSize,
		})
	}
}

func GetReview(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		reviewID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review id"})
			return
		}
		var review pkg.Review
		err = conn.QueryRow(c.Request.Context(), "SELECT id, service_id, rating, comment, time_created FROM review WHERE id = $1", reviewID).Scan(&review.Id, &review.Service_id, &review.Rating, &review.Comment, &review.Time_created)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, review)
	}
}

func CreateReview(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		var review pkg.Review
		if err := c.BindJSON(&review); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Validate rating is between 1-5
		if review.Rating < 1 || review.Rating > 5 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Rating must be between 1 and 5"})
			return
		}

		err := conn.QueryRow(c.Request.Context(),
			`INSERT INTO review (service_id, rating, comment)
			 VALUES ($1, $2, $3)
			 RETURNING id, time_created`,
			review.Service_id, review.Rating, review.Comment).
			Scan(&review.Id, &review.Time_created)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, review)
	}
}

func UpdateReview(conn *pgx.Conn) gin.HandlerFunc {
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

		// Validate rating is between 1-5
		if review.Rating < 1 || review.Rating > 5 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Rating must be between 1 and 5"})
			return
		}

		err = conn.QueryRow(c.Request.Context(),
			`UPDATE review SET service_id = $1, rating = $2, comment = $3
			 WHERE id = $4
			 RETURNING id, service_id, rating, comment, time_created`,
			review.Service_id, review.Rating, review.Comment, reviewID).
			Scan(&review.Id, &review.Service_id, &review.Rating, &review.Comment, &review.Time_created)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, review)
	}
}

func DeleteReview(conn *pgx.Conn) gin.HandlerFunc {
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
