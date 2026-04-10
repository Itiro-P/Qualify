package handlers

import (
	"context"
	"fmt"
	"main/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func GetAnalysts(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Build query with filters
		query := `SELECT id, name, email, phone, time_created, country_code, 
                         country_name, country_state, city, timezone, 
                         hourly_rate, total_reviews, mean_rating 
                  FROM analyst WHERE 1=1`

		args := []interface{}{}
		argCounter := 1

		// Name filter (partial match)
		if name := c.Query("name"); name != "" {
			query += fmt.Sprintf(" AND name ILIKE $%d", argCounter)
			args = append(args, "%"+name+"%")
			argCounter++
		}

		// Country filter
		if country := c.Query("country"); country != "" {
			query += fmt.Sprintf(" AND country_name ILIKE $%d", argCounter)
			args = append(args, "%"+country+"%")
			argCounter++
		}

		// City filter
		if city := c.Query("city"); city != "" {
			query += fmt.Sprintf(" AND city ILIKE $%d", argCounter)
			args = append(args, "%"+city+"%")
			argCounter++
		}

		// Hourly rate range filter
		if minRate := c.Query("min_hourly_rate"); minRate != "" {
			query += fmt.Sprintf(" AND hourly_rate >= $%d", argCounter)
			args = append(args, minRate)
			argCounter++
		}

		if maxRate := c.Query("max_hourly_rate"); maxRate != "" {
			query += fmt.Sprintf(" AND hourly_rate <= $%d", argCounter)
			args = append(args, maxRate)
			argCounter++
		}

		// Total reviews filter (minimum)
		if minReviews := c.Query("min_total_reviews"); minReviews != "" {
			query += fmt.Sprintf(" AND total_reviews >= $%d", argCounter)
			args = append(args, minReviews)
			argCounter++
		}

		// Mean rating filter (minimum)
		if minRating := c.Query("min_mean_rating"); minRating != "" {
			query += fmt.Sprintf(" AND mean_rating >= $%d", argCounter)
			args = append(args, minRating)
			argCounter++
		}

		// Optional: Add sorting
		if sortBy := c.Query("sort_by"); sortBy != "" {
			order := c.DefaultQuery("order", "ASC")
			query += fmt.Sprintf(" ORDER BY %s %s", sortBy, order)
		}

		// Execute query
		rows, err := conn.Query(context.Background(), query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar prestadores: " + err.Error()})
			return
		}
		defer rows.Close()

		var analysts []pkg.Analyst

		// Iterate through results
		for rows.Next() {
			var analyst pkg.Analyst
			err := rows.Scan(
				&analyst.ID,
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
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao escanear data: " + err.Error()})
				return
			}
			analysts = append(analysts, analyst)
		}

		// Check for errors from iterating over rows
		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao iterar prestadores: " + err.Error()})
			return
		}

		// Return results
		c.JSON(http.StatusOK, gin.H{
			"analysts": analysts,
			"count":    len(analysts),
		})
	}
}

func GetAnalyst(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var analyst pkg.Analyst
		err := conn.QueryRow(context.Background(), "SELECT id, name FROM \"analyst\" WHERE id = $1", id).Scan(&analyst.ID, &analyst.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar prestador"})
			return
		}

		c.JSON(http.StatusOK, analyst)
	}
}
