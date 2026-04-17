package handlers

import (
	"fmt"
	"main/pkg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func GetAnalysts(conn *pgx.Conn) gin.HandlerFunc {
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

		if country := c.Query("country"); country != "" {
			query += fmt.Sprintf(" AND u.country_name ILIKE $%d", argCounter)
			args = append(args, "%"+country+"%")
			argCounter++
		}

		if city := c.Query("city"); city != "" {
			query += fmt.Sprintf(" AND u.city ILIKE $%d", argCounter)
			args = append(args, "%"+city+"%")
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
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar analistas: " + err.Error()})
			return
		}
		defer rows.Close()

		var analysts []pkg.Analyst

		for rows.Next() {
			var analyst pkg.Analyst
			err := rows.Scan(
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
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao escanear data: " + err.Error()})
				return
			}
			analysts = append(analysts, analyst)
		}

		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao iterar analistas: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"analysts": analysts,
			"count":    len(analysts),
		})
	}
}

func GetAnalyst(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		analystID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid analyst id"})
			return
		}

		var analyst pkg.Analyst
		err = conn.QueryRow(c.Request.Context(), `
			SELECT u.id, u.name, u.email, u.phone, u.time_created, 
			       u.country_code, u.country_name, u.country_state, u.city, u.timezone,
			       a.hourly_rate, a.total_reviews, a.mean_rating
			FROM "user" u
			JOIN analyst a ON a.id = u.id
			WHERE u.id = $1`, analystID).Scan(
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
		)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "analyst not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, analyst)
	}
}

func CreateAnalyst(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		var analyst pkg.Analyst
		if err := c.BindJSON(&analyst); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		err := conn.QueryRow(c.Request.Context(),
			`INSERT INTO analyst (name, email, phone, country_code, country_name, country_state, city, timezone, hourly_rate, total_reviews, mean_rating)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			 RETURNING id, time_created`,
			analyst.Name, analyst.Email, analyst.Phone, analyst.Country_code, analyst.Country_name, analyst.Country_state,
			analyst.City, analyst.Timezone, analyst.Hourly_rate, analyst.Total_reviews, analyst.Mean_rating).
			Scan(&analyst.Id, &analyst.Time_created)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, analyst)
	}
}

func UpdateAnalyst(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		analystID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid analyst id"})
			return
		}
		var analyst pkg.Analyst
		if err := c.BindJSON(&analyst); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		err = conn.QueryRow(c.Request.Context(),
			`UPDATE analyst SET name = $1, email = $2, phone = $3, country_code = $4, 
			 country_name = $5, country_state = $6, city = $7, timezone = $8, hourly_rate = $9, total_reviews = $10, mean_rating = $11
			 WHERE id = $12
			 RETURNING id, name, email, phone, time_created, country_code, country_name, country_state, city, timezone, hourly_rate, total_reviews, mean_rating`,
			analyst.Name, analyst.Email, analyst.Phone, analyst.Country_code, analyst.Country_name, analyst.Country_state,
			analyst.City, analyst.Timezone, analyst.Hourly_rate, analyst.Total_reviews, analyst.Mean_rating, analystID).
			Scan(&analyst.Id, &analyst.Name, &analyst.Email, &analyst.Phone, &analyst.Time_created, &analyst.Country_code,
				&analyst.Country_name, &analyst.Country_state, &analyst.City, &analyst.Timezone,
				&analyst.Hourly_rate, &analyst.Total_reviews, &analyst.Mean_rating)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "analyst not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, analyst)
	}
}

func DeleteAnalyst(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		analystID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid analyst id"})
			return
		}

		result, err := conn.Exec(c.Request.Context(), `DELETE FROM analyst WHERE id = $1`, analystID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "analyst not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "analyst deleted successfully"})
	}
}
