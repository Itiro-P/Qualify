package handlers

import (
	"fmt"
	"main/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func GetReviews(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Build query with filters
		query := `SELECT id, client_id, analyst_id, rating, comment, time_created 
                  FROM review WHERE 1=1`

		args := []interface{}{}
		argCounter := 1

		// Filter by client_id
		if clientID := c.Query("client_id"); clientID != "" {
			query += fmt.Sprintf(" AND client_id = $%d", argCounter)
			args = append(args, clientID)
			argCounter++
		}

		// Filter by analyst_id
		if analystID := c.Query("analyst_id"); analystID != "" {
			query += fmt.Sprintf(" AND analyst_id = $%d", argCounter)
			args = append(args, analystID)
			argCounter++
		}

		// Rating filter (exact match)
		if rating := c.Query("rating"); rating != "" {
			query += fmt.Sprintf(" AND rating = $%d", argCounter)
			args = append(args, rating)
			argCounter++
		}

		// Rating range filter
		if minRating := c.Query("min_rating"); minRating != "" {
			query += fmt.Sprintf(" AND rating >= $%d", argCounter)
			args = append(args, minRating)
			argCounter++
		}

		if maxRating := c.Query("max_rating"); maxRating != "" {
			query += fmt.Sprintf(" AND rating <= $%d", argCounter)
			args = append(args, maxRating)
			argCounter++
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
				"id": true, "client_id": true, "analyst_id": true,
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
			fmt.Sscanf(p, "%d", &page)
		}

		pageSize := 20 // Default page size for reviews
		if ps := c.Query("page_size"); ps != "" {
			fmt.Sscanf(ps, "%d", &pageSize)
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
				&review.Client_id,
				&review.Analyst_id,
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
		var review pkg.Review
		err := conn.QueryRow(c.Request.Context(), "SELECT id, name FROM \"review\" WHERE id = $1", id).Scan(&review.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar avaliação"})
			return
		}

		c.JSON(http.StatusOK, review)
	}
}
