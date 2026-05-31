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

		if filters.Name != "" {
			builder = builder.Where(squirrel.ILike{"u.name": pkg.PutPercent(filters.Name)})
		}

		if filters.Email != "" {
			builder = builder.Where(squirrel.ILike{"u.email": pkg.PutPercent(filters.Email)})
		}

		if filters.Country != "" {
			builder = builder.Where(squirrel.ILike{"u.country_name": pkg.PutPercent(filters.Country)})
		}

		if filters.CountryCode != "" {
			builder = builder.Where(squirrel.ILike{"u.country_code": pkg.PutPercent(filters.CountryCode)})
		}

		if filters.CountryState != "" {
			builder = builder.Where(squirrel.ILike{"u.country_state": pkg.PutPercent(filters.CountryState)})
		}

		if filters.City != "" {
			builder = builder.Where(squirrel.ILike{"u.city": pkg.PutPercent(filters.City)})
		}

		if filters.Timezone != "" {
			builder = builder.Where(squirrel.ILike{"u.timezone": pkg.PutPercent(filters.Timezone)})
		}

		if filters.MinProposedBudget != nil {
			builder = builder.Where(squirrel.GtOrEq{"c.proposed_budget": filters.MinProposedBudget})
		}

		if filters.MaxProposedBudget != nil {
			builder = builder.Where(squirrel.LtOrEq{"c.proposed_budget": filters.MaxProposedBudget})
		}

		orderClause := filters.SortOptions.ValidateSort(pkg.ClientSortFields)

		if orderClause != "" {
			builder = builder.OrderBy(orderClause)
		} else {
			builder = builder.OrderBy("time_created DESC")
		}

		builder = builder.Limit(uint64(filters.PageSize)).Offset(uint64(filters.Offset()))

		query, args, err := builder.ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		rows, err := conn.Query(c.Request.Context(), query, args...)
		if pkg.HandleErr(c, err) {
			return
		}
		defer rows.Close()

		var clients []pkg.Client
		for rows.Next() {
			client, err := pkg.ScanClient(rows)
			if pkg.HandleErr(c, err) {
				return
			}
			clients = append(clients, client)
		}

		if err = rows.Err(); pkg.HandleErr(c, err) {
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
			Set("name", client.Name).
			Set("email", client.Email).
			Set("phone", client.Phone).
			Set("country_code", client.Country_code).
			Set("country_name", client.Country_name).
			Set("country_state", client.Country_state).
			Set("city", client.City).
			Set("timezone", client.Timezone).
			Where(squirrel.Eq{"id": id}).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if pkg.HandleErr(c, err) {
			return
		}
		if _, err = tx.Exec(c.Request.Context(), userQuery, userArgs...); pkg.HandleErr(c, err) {
			return
		}

		clientQuery, clientArgs, err := squirrel.Update("client").
			Set("proposed_budget", client.Proposed_budget).
			Where(squirrel.Eq{"id": id}).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if pkg.HandleErr(c, err) {
			return
		}
		if _, err = tx.Exec(c.Request.Context(), clientQuery, clientArgs...); pkg.HandleErr(c, err) {
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
		}

		if err = tx.Commit(c.Request.Context()); pkg.HandleErr(c, err) {
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
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		result, err := conn.Exec(c.Request.Context(), query, args...)
		if pkg.HandleErr(c, err) {
			return
		}

		if result.RowsAffected() == 0 {
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

		userBuilder := squirrel.Update(`"user"`).PlaceholderFormat(squirrel.Dollar)
		clientBuilder := squirrel.Update("client").PlaceholderFormat(squirrel.Dollar)
		userHasFields := false
		clientHasFields := false

		if req.Name != nil {
			userBuilder = userBuilder.Set("name", *req.Name)
			userHasFields = true
		}
		if req.Email != nil {
			userBuilder = userBuilder.Set("email", *req.Email)
			userHasFields = true
		}
		if req.Phone != nil {
			userBuilder = userBuilder.Set("phone", *req.Phone)
			userHasFields = true
		}
		if req.Country_code != nil {
			userBuilder = userBuilder.Set("country_code", *req.Country_code)
			userHasFields = true
		}
		if req.Country_name != nil {
			userBuilder = userBuilder.Set("country_name", *req.Country_name)
			userHasFields = true
		}
		if req.Country_state != nil {
			userBuilder = userBuilder.Set("country_state", *req.Country_state)
			userHasFields = true
		}
		if req.City != nil {
			userBuilder = userBuilder.Set("city", *req.City)
			userHasFields = true
		}
		if req.Timezone != nil {
			userBuilder = userBuilder.Set("timezone", *req.Timezone)
			userHasFields = true
		}

		if req.Proposed_budget != nil {
			clientBuilder = clientBuilder.Set("proposed_budget", *req.Proposed_budget)
			clientHasFields = true
		}

		if !userHasFields && !clientHasFields {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "No valid fields were given"))
			return
		}

		tx, err := conn.Begin(c.Request.Context())
		if pkg.HandleErr(c, err) {
			return
		}
		defer tx.Rollback(c.Request.Context())

		if userHasFields {
			query, args, err := userBuilder.Where(squirrel.Eq{"id": id}).ToSql()
			if pkg.HandleErr(c, err) {
				return
			}
			if _, err = tx.Exec(c.Request.Context(), query, args...); pkg.HandleErr(c, err) {
				return
			}
		}

		if clientHasFields {
			query, args, err := clientBuilder.Where(squirrel.Eq{"id": id}).ToSql()
			if pkg.HandleErr(c, err) {
				return
			}
			if _, err = tx.Exec(c.Request.Context(), query, args...); pkg.HandleErr(c, err) {
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
		}

		if err = tx.Commit(c.Request.Context()); pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusOK, pkg.ClientResponse{Client: client})
	}
}
