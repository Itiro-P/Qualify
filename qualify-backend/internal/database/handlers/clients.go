package handlers

import (
	"main/internal/database/services"
	"main/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/n-r-w/squirrel"
)

const clientSelect = `u.id, u.name, u.email, u.phone, u.time_created, 
	u.country_code, u.country_name, u.country_state, u.city, u.timezone, 
	c.proposed_budget`

const clientJoin = `"user" u JOIN client c ON c.id = u.id`

// GetClients godoc
// @Summary Listar clientes
// @Description Retorna uma lista de clientes com filtros opcionais
// @Tags Clientes
// @Accept json
// @Produce json
// @Param name query string false "Nome parcial para busca"
// @Param email query string false "Email parcial para busca"
// @Param country query string false "País"
// @Param country_code query string false "Código do país"
// @Param country_state query string false "Estado"
// @Param city query string false "Cidade"
// @Param timezone query string false "Fuso horário"
// @Param min_proposed_budget query number false "Orçamento mínimo"
// @Param max_proposed_budget query number false "Orçamento máximo"
// @Success 200 {object} pkg.ClientsResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /clients [get]
func GetClients(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		builder := squirrel.Select(clientSelect).From(clientJoin).PlaceholderFormat(squirrel.Dollar)
		var filters pkg.ClientFilter
		err := c.ShouldBindQuery(&filters)
		if pkg.HandleErr(c, err) {
			return
		}
		filters.Normalize()

		query, args, err := pkg.BuildFilterClient(pkg.BuildFilterUser(builder, filters.UserFilter), filters).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		rows, err := conn.Query(c.Request.Context(), query, args...)
		if pkg.HandleErr(c, err) {
			return
		}
		defer rows.Close()

		clients, err := pkg.ScanRows(c, rows, pkg.ScanClient)
		if pkg.HandleErr(c, err) {
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
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /users/{id}/client [get]
func GetClient(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		query, args, err := squirrel.Select(clientSelect).
			From(clientJoin).Where(squirrel.Eq{"u.id": id}).
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		client, err := pkg.ScanClient(conn.QueryRow(c.Request.Context(), query, args...))
		if pkg.HandleErr(c, err) {
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
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 409 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /users/{id}/client [post]
func CreateClient(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		// Checando se o cliente já existe
		var clientExists bool
		err = conn.QueryRow(c.Request.Context(),
			`SELECT EXISTS(SELECT 1 FROM client WHERE id = $1)`, id).Scan(&clientExists)
		if clientExists {
			c.JSON(http.StatusConflict, pkg.Internal(c.FullPath(), "Client already exists"))
			return
		}

		var request struct {
			Proposed_budget float64 `json:"proposed_budget"`
		}
		if err := c.BindJSON(&request); pkg.HandleErr(c, err) {
			return
		}

		client, err := services.AssignClientRole(c.Request.Context(), conn, id, request.Proposed_budget)
		if pkg.HandleErr(c, err) {
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
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /users/{id}/client [put]
func UpdateClient(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		var client pkg.Client
		if err := c.BindJSON(&client); pkg.HandleErr(c, err) {
			return
		}

		if client.Name == "" {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Received empty name"))
			return
		}
		if client.Email == "" {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Received empty email"))
			return
		}
		if len(client.Country_code) != 2 {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Country code must be exactly 2 characters"))
			return
		}

		tx, err := conn.Begin(c.Request.Context())
		if pkg.HandleErr(c, err) {
			return
		}
		defer tx.Rollback(c.Request.Context())

		userQuery, userArgs, err := squirrel.Update(`"user"`).
			SetMap(map[string]any{
				"name":          client.Name,
				"email":         client.Email,
				"phone":         client.Phone,
				"country_code":  client.Country_code,
				"country_name":  client.Country_name,
				"country_state": client.Country_state,
				"city":          client.City,
				"timezone":      client.Timezone,
			}).Where(squirrel.Eq{"id": id}).
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		} else if _, err = tx.Exec(c.Request.Context(), userQuery, userArgs...); pkg.HandleErr(c, err) {
			return
		}

		clientQuery, clientArgs, err := squirrel.Update("client").
			Set("proposed_budget", client.Proposed_budget).
			Where(squirrel.Eq{"id": id}).
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		} else if _, err = tx.Exec(c.Request.Context(), clientQuery, clientArgs...); pkg.HandleErr(c, err) {
			return
		}

		selectQuery, selectArgs, err := squirrel.Select(clientSelect).
			From(clientJoin).Where(squirrel.Eq{"u.id": id}).
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		client, err = pkg.ScanClient(tx.QueryRow(c.Request.Context(), selectQuery, selectArgs...))
		if pkg.HandleErr(c, err) {
			return
		} else if err = tx.Commit(c.Request.Context()); pkg.HandleErr(c, err) {
			return
		}

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
// @Success 204 "Deleção com sucesso"
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /users/{id}/client [delete]
func DeleteClient(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		query, args, err := squirrel.Delete("client").
			Where(squirrel.Eq{"id": id}).
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		result, err := conn.Exec(c.Request.Context(), query, args...)
		if pkg.HandleErr(c, err) {
			return
		} else if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Client not found"))
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// @Summary Atualizar parcialmente um cliente
// @Description Atualiza um ou mais campos do usuário/cliente
// @Tags Clientes
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário"
// @Param client body pkg.ClientUpdateRequest true "Campos opcionais para atualização"
// @Success 200 {object} pkg.ClientResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /users/{id}/client [patch]
func UpdateClientPartial(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		var req pkg.ClientUpdateRequest
		if err := c.BindJSON(&req); pkg.HandleErr(c, err) {
			return
		}

		userFields, clientFields := pkg.BuildUpdateClient(req)
		if len(userFields) <= 0 && len(clientFields) <= 0 {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "No valid fields were given"))
			return
		}

		tx, err := conn.Begin(c.Request.Context())
		if pkg.HandleErr(c, err) {
			return
		}
		defer tx.Rollback(c.Request.Context())

		if len(userFields) > 0 {
			query, args, err := squirrel.Update(`"user"`).
				SetMap(userFields).Where(squirrel.Eq{"id": id}).
				PlaceholderFormat(squirrel.Dollar).ToSql()
			if pkg.HandleErr(c, err) {
				return
			} else if _, err = tx.Exec(c.Request.Context(), query, args...); pkg.HandleErr(c, err) {
				return
			}
		}

		if len(clientFields) > 0 {
			query, args, err := squirrel.Update("client").
				SetMap(clientFields).Where(squirrel.Eq{"id": id}).
				PlaceholderFormat(squirrel.Dollar).ToSql()
			if pkg.HandleErr(c, err) {
				return
			} else if _, err = tx.Exec(c.Request.Context(), query, args...); pkg.HandleErr(c, err) {
				return
			}
		}

		selectQuery, selectArgs, err := squirrel.Select(clientSelect).
			From(clientJoin).Where(squirrel.Eq{"u.id": id}).
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		client, err := pkg.ScanClient(tx.QueryRow(c.Request.Context(), selectQuery, selectArgs...))
		if pkg.HandleErr(c, err) {
			return
		} else if err = tx.Commit(c.Request.Context()); pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusOK, pkg.ClientResponse{Client: client})
	}
}
