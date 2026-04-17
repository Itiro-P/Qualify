package handlers

import (
	"fmt"
	"main/pkg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func GetClients(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Build query with filters
		query := `SELECT id, name, email, phone, time_created, country_code, 
                         country_name, country_state, city, timezone, 
                         proposed_budget 
                  FROM client WHERE 1=1`

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

		// Proposed budget range filter
		if minBudget := c.Query("min_proposed_budget"); minBudget != "" {
			min, err := strconv.ParseFloat(minBudget, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid min_proposed_budget"})
				return
			}
			query += fmt.Sprintf(" AND proposed_budget >= $%d", argCounter)
			args = append(args, min)
			argCounter++
		}

		if maxBudget := c.Query("max_proposed_budget"); maxBudget != "" {
			max, err := strconv.ParseFloat(maxBudget, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid max_proposed_budget"})
				return
			}
			query += fmt.Sprintf(" AND proposed_budget <= $%d", argCounter)
			args = append(args, max)
			argCounter++
		}

		// Optional: Add sorting
		if sortBy := c.Query("sort_by"); sortBy != "" {
			// Validate sortBy to prevent SQL injection
			allowedSortFields := map[string]bool{
				"name": true, "country_name": true, "city": true,
				"proposed_budget": true, "time_created": true,
			}
			if allowedSortFields[sortBy] {
				order := c.DefaultQuery("order", "ASC")
				if order == "ASC" || order == "DESC" {
					query += fmt.Sprintf(" ORDER BY %s %s", sortBy, order)
				}
			}
		}

		// Execute query
		rows, err := conn.Query(c.Request.Context(), query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching clients: " + err.Error()})
			return
		}
		defer rows.Close()

		var clients []pkg.Client

		// Iterate through results
		for rows.Next() {
			var client pkg.Client
			err := rows.Scan(
				&client.Id,
				&client.Name,
				&client.Email,
				&client.Phone,
				&client.Time_created,
				&client.Country_code,
				&client.Country_name,
				&client.Country_state,
				&client.City,
				&client.Timezone,
				&client.Proposed_budget,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning client data: " + err.Error()})
				return
			}
			clients = append(clients, client)
		}

		// Check for errors from iterating over rows
		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating clients: " + err.Error()})
			return
		}

		// Return results
		c.JSON(http.StatusOK, gin.H{
			"clients": clients,
			"count":   len(clients),
		})
	}
}

func GetClient(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		clientID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client id"})
			return
		}
		var client pkg.Client
		err = conn.QueryRow(c.Request.Context(), `SELECT id, name, email, phone, time_created, country_code, 
                         country_name, country_state, city, timezone, proposed_budget 
                  FROM client WHERE id = $1`, clientID).Scan(
			&client.Id,
			&client.Name,
			&client.Email,
			&client.Phone,
			&client.Time_created,
			&client.Country_code,
			&client.Country_name,
			&client.Country_state,
			&client.City,
			&client.Timezone,
			&client.Proposed_budget,
		)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, client)
	}
}

func CreateClient(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		var client pkg.Client
		if err := c.BindJSON(&client); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		err := conn.QueryRow(c.Request.Context(),
			`INSERT INTO client (name, email, phone, country_code, country_name, country_state, city, timezone, proposed_budget)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			 RETURNING id, time_created`,
			client.Name, client.Email, client.Phone, client.Country_code, client.Country_name, client.Country_state,
			client.City, client.Timezone, client.Proposed_budget).
			Scan(&client.Id, &client.Time_created)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, client)
	}
}

func UpdateClient(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		clientID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client id"})
			return
		}
		var client pkg.Client
		if err := c.BindJSON(&client); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		err = conn.QueryRow(c.Request.Context(),
			`UPDATE client SET name = $1, email = $2, phone = $3, country_code = $4, 
			 country_name = $5, country_state = $6, city = $7, timezone = $8, proposed_budget = $9
			 WHERE id = $10
			 RETURNING id, name, email, phone, time_created, country_code, country_name, country_state, city, timezone, proposed_budget`,
			client.Name, client.Email, client.Phone, client.Country_code, client.Country_name, client.Country_state,
			client.City, client.Timezone, client.Proposed_budget, clientID).
			Scan(&client.Id, &client.Name, &client.Email, &client.Phone, &client.Time_created, &client.Country_code,
				&client.Country_name, &client.Country_state, &client.City, &client.Timezone, &client.Proposed_budget)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, client)
	}
}

func DeleteClient(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		clientID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client id"})
			return
		}

		result, err := conn.Exec(c.Request.Context(), `DELETE FROM client WHERE id = $1`, clientID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "client deleted successfully"})
	}
}
