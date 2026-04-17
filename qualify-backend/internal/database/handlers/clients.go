package handlers

import (
	"fmt"
	"main/internal/database/services"
	"main/pkg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func GetClients(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `
			SELECT u.id, u.name, u.email, u.phone, u.time_created,
			       u.country_code, u.country_name, u.country_state, u.city, u.timezone,
			       c.proposed_budget
			FROM "user" u
			JOIN client c ON c.id = u.id
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

		if minBudget := c.Query("min_proposed_budget"); minBudget != "" {
			if minBudgetVal, err := strconv.ParseFloat(minBudget, 64); err == nil {
				query += fmt.Sprintf(" AND c.proposed_budget >= $%d", argCounter)
				args = append(args, minBudgetVal)
				argCounter++
			}
		}

		if maxBudget := c.Query("max_proposed_budget"); maxBudget != "" {
			if maxBudgetVal, err := strconv.ParseFloat(maxBudget, 64); err == nil {
				query += fmt.Sprintf(" AND c.proposed_budget <= $%d", argCounter)
				args = append(args, maxBudgetVal)
				argCounter++
			}
		}

		allowedSortFields := map[string]bool{
			"name": true, "country_name": true, "city": true,
			"proposed_budget": true, "time_created": true,
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching clients: " + err.Error()})
			return
		}
		defer rows.Close()

		var clients []pkg.Client
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

		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating clients: " + err.Error()})
			return
		}

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
		err = conn.QueryRow(c.Request.Context(), `
			SELECT u.id, u.name, u.email, u.phone, u.time_created, 
			       u.country_code, u.country_name, u.country_state, u.city, u.timezone,
			       c.proposed_budget
			FROM "user" u
			JOIN client c ON c.id = u.id
			WHERE u.id = $1`, clientID).Scan(
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
		userIDParam := c.Param("id")
		userID, err := strconv.Atoi(userIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		var request struct {
			Proposed_budget float64 `json:"proposed_budget"`
		}
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		client, err := services.AssignClientRole(c.Request.Context(), conn, userID, request.Proposed_budget)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign client role: " + err.Error()})
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

		// Update user table
		_, err = conn.Exec(c.Request.Context(),
			`UPDATE "user" SET name = $1, email = $2, phone = $3, country_code = $4, 
			 country_name = $5, country_state = $6, city = $7, timezone = $8
			 WHERE id = $9`,
			client.Name, client.Email, client.Phone, client.Country_code, client.Country_name, client.Country_state,
			client.City, client.Timezone, clientID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user: " + err.Error()})
			return
		}

		// Update client table
		_, err = conn.Exec(c.Request.Context(),
			`UPDATE client SET proposed_budget = $1
			 WHERE id = $2`,
			client.Proposed_budget, clientID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update client: " + err.Error()})
			return
		}

		// Fetch the updated client
		err = conn.QueryRow(c.Request.Context(), `
			SELECT u.id, u.name, u.email, u.phone, u.time_created, 
			       u.country_code, u.country_name, u.country_state, u.city, u.timezone,
			       c.proposed_budget
			FROM "user" u
			JOIN client c ON c.id = u.id
			WHERE u.id = $1`, clientID).Scan(
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated client: " + err.Error()})
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
