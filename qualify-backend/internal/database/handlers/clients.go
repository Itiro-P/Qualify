package handlers

import (
	"fmt"
	"main/internal/database/services"
	"main/pkg"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// GetClients godoc
// @Summary Listar clientes
// @Description Retorna uma lista de clientes com filtros opcionais
// @Tags Clientes
// @Accept json
// @Produce json
// @Param name query string false "Nome parcial para busca"
// @Param country query string false "País"
// @Param city query string false "Cidade"
// @Param min_proposed_budget query number false "Orçamento mínimo"
// @Param max_proposed_budget query number false "Orçamento máximo"
// @Success 200 {object} pkg.ClientsResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /clients [get]
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

		c.JSON(http.StatusOK, pkg.ClientsResponse{
			Clients: clients,
			Count:   len(clients),
		})
	}
}

// GetClient godoc
// @Summary Obter cliente
// @Description Retorna os detalhes de um cliente pelo ID do usuário
// @Tags Clientes
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário"
// @Success 200 {object} pkg.ClientResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/client [get]
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

		c.JSON(http.StatusOK, pkg.ClientResponse{Client: client})
	}
}

// CreateClient godoc
// @Summary Criar papel de cliente
// @Description Atribui o papel de cliente a um usuário existente
// @Tags Clientes
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário"
// @Param proposed_budget body object true "{\"proposed_budget\": number}"
// @Success 201 {object} pkg.ClientResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/client [post]
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

		c.JSON(http.StatusCreated, pkg.ClientResponse{Client: *client})
	}
}

// UpdateClient godoc
// @Summary Atualizar cliente
// @Description Atualiza dados do cliente pelo ID do usuário
// @Tags Clientes
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário"
// @Param client body pkg.Client true "Objeto cliente"
// @Success 200 {object} pkg.ClientResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/client [put]
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

		// Validando parâmetros obrigatórios
		if client.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user name is required"})
			return
		}
		if client.Email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user email is required"})
			return
		}
		if len(client.Country_code) != 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "country_code must be exactly 2 characters"})
			return
		}

		// Make update atomic across user + client tables
		tx, err := conn.Begin(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction: " + err.Error()})
			return
		}
		// ensure rollback on error
		committed := false
		defer func() {
			if !committed {
				_ = tx.Rollback(c.Request.Context())
			}
		}()

		_, err = tx.Exec(c.Request.Context(),
			`UPDATE "user" SET name = $1, email = $2, phone = $3, country_code = $4, 
			 country_name = $5, country_state = $6, city = $7, timezone = $8
			 WHERE id = $9`,
			client.Name, client.Email, client.Phone, client.Country_code, client.Country_name, client.Country_state,
			client.City, client.Timezone, clientID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user: " + err.Error()})
			return
		}

		_, err = tx.Exec(c.Request.Context(),
			`UPDATE client SET proposed_budget = $1
			 WHERE id = $2`,
			client.Proposed_budget, clientID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update client: " + err.Error()})
			return
		}

		// Fetch the updated client within transaction
		err = tx.QueryRow(c.Request.Context(), `
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

		if err = tx.Commit(c.Request.Context()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction: " + err.Error()})
			return
		}
		committed = true

		c.JSON(http.StatusOK, pkg.ClientResponse{Client: client})
	}
}

// DeleteClient godoc
// @Summary Excluir cliente
// @Description Remove o papel de cliente de um usuário
// @Tags Clientes
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/client [delete]
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

// UpdateClientPartial godoc
// @Summary Atualizar parcialmente um cliente
// @Description Atualiza um ou mais campos do usuário/cliente
// @Tags Clientes
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário"
// @Param client body pkg.ClientUpdateRequest true "Campos opcionais para atualização"
// @Success 200 {object} pkg.ClientResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/client [patch]
func UpdateClientPartial(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		clientID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client id"})
			return
		}

		var req pkg.ClientUpdateRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		userSet := []string{}
		userArgs := []interface{}{}
		i := 1
		if req.Name != nil {
			userSet = append(userSet, fmt.Sprintf("name = $%d", i)); userArgs = append(userArgs, *req.Name); i++
		}
		if req.Email != nil {
			userSet = append(userSet, fmt.Sprintf("email = $%d", i)); userArgs = append(userArgs, *req.Email); i++
		}
		if req.Phone != nil {
			userSet = append(userSet, fmt.Sprintf("phone = $%d", i)); userArgs = append(userArgs, *req.Phone); i++
		}
		if req.Country_code != nil {
			userSet = append(userSet, fmt.Sprintf("country_code = $%d", i)); userArgs = append(userArgs, *req.Country_code); i++
		}
		if req.Country_name != nil {
			userSet = append(userSet, fmt.Sprintf("country_name = $%d", i)); userArgs = append(userArgs, *req.Country_name); i++
		}
		if req.Country_state != nil {
			userSet = append(userSet, fmt.Sprintf("country_state = $%d", i)); userArgs = append(userArgs, *req.Country_state); i++
		}
		if req.City != nil {
			userSet = append(userSet, fmt.Sprintf("city = $%d", i)); userArgs = append(userArgs, *req.City); i++
		}
		if req.Timezone != nil {
			userSet = append(userSet, fmt.Sprintf("timezone = $%d", i)); userArgs = append(userArgs, *req.Timezone); i++
		}

		clientSet := []string{}
		clientArgs := []interface{}{}
		if req.Proposed_budget != nil {
			clientSet = append(clientSet, fmt.Sprintf("proposed_budget = $%d", i)); clientArgs = append(clientArgs, *req.Proposed_budget); i++
		}

		if len(userSet) == 0 && len(clientSet) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
			return
		}

		tx, err := conn.Begin(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to begin transaction: " + err.Error()})
			return
		}
		defer func() { if err != nil { _ = tx.Rollback(c.Request.Context()) } }()

		if len(userSet) > 0 {
			userArgs = append(userArgs, clientID)
			userQuery := fmt.Sprintf(`UPDATE "user" SET %s WHERE id = $%d`, strings.Join(userSet, ", "), len(userArgs))
			if _, err := tx.Exec(c.Request.Context(), userQuery, userArgs...); err != nil {
				_ = tx.Rollback(c.Request.Context())
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		if len(clientSet) > 0 {
			clientArgs = append(userArgs, clientArgs...)
			clientArgs = append(clientArgs, clientID)
			clientQuery := fmt.Sprintf(`UPDATE client SET %s WHERE id = $%d`, strings.Join(clientSet, ", "), len(clientArgs))
			if _, err := tx.Exec(c.Request.Context(), clientQuery, clientArgs...); err != nil {
				_ = tx.Rollback(c.Request.Context())
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		var client pkg.Client
		err = tx.QueryRow(c.Request.Context(), `
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
			_ = tx.Rollback(c.Request.Context())
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err = tx.Commit(c.Request.Context()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transaction: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, pkg.ClientResponse{Client: client})
	}
}
