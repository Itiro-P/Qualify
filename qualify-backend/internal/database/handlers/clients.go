package handlers

import (
	"fmt"
	"main/pkg"
	"net/http"

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
			query += fmt.Sprintf(" AND proposed_budget >= $%d", argCounter)
			args = append(args, minBudget)
			argCounter++
		}

		if maxBudget := c.Query("max_proposed_budget"); maxBudget != "" {
			query += fmt.Sprintf(" AND proposed_budget <= $%d", argCounter)
			args = append(args, maxBudget)
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
		var client pkg.Client
		err := conn.QueryRow(c.Request.Context(), "SELECT id, name FROM \"client\" WHERE id = $1", id).Scan(&client.Id, &client.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar cliente"})
			return
		}

		c.JSON(http.StatusOK, client)
	}
}
