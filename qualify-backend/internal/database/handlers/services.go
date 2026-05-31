package handlers

import (
	"main/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/n-r-w/squirrel"
)

const serviceSelect = "id, title, content, proposal_letter_id, hourly_rate, status, time_created"

func serviceBuilder(filters pkg.ServiceFilter) squirrel.SelectBuilder {
	builder := squirrel.Select(serviceSelect).
		From("service").
		PlaceholderFormat(squirrel.Dollar)

	if filters.Status != "" {
		builder = builder.Where(squirrel.Eq{"status": filters.Status})
	}
	if filters.ProposalId != nil {
		builder = builder.Where(squirrel.Eq{"proposal_letter_id": *filters.ProposalId})
	}
	if filters.Title != "" {
		builder = builder.Where(squirrel.ILike{"title": pkg.PutPercent(filters.Title)})
	}
	if filters.Content != "" {
		builder = builder.Where(squirrel.ILike{"content": pkg.PutPercent(filters.Content)})
	}
	if filters.MinHourlyRate != nil {
		builder = builder.Where(squirrel.GtOrEq{"hourly_rate": *filters.MinHourlyRate})
	}
	if filters.MaxHourlyRate != nil {
		builder = builder.Where(squirrel.LtOrEq{"hourly_rate": *filters.MaxHourlyRate})
	}

	if order := filters.ValidateSort(pkg.ServiceSortFields); order != "" {
		builder = builder.OrderBy(order)
	} else {
		builder = builder.OrderBy("time_created DESC")
	}

	builder = builder.Limit(uint64(filters.PageSize)).Offset(uint64(filters.Offset()))
	return builder
}

func scanServices(c *gin.Context, conn *pgxpool.Pool, builder squirrel.SelectBuilder) {
	query, args, err := builder.ToSql()
	if pkg.HandleErr(c, err) {
		return
	}

	rows, err := conn.Query(c.Request.Context(), query, args...)
	if pkg.HandleErr(c, err) {
		return
	}
	defer rows.Close()

	var services []pkg.Service
	for rows.Next() {
		s, err := pkg.ScanService(rows)
		if pkg.HandleErr(c, err) {
			return
		}
		services = append(services, s)
	}
	if err = rows.Err(); pkg.HandleErr(c, err) {
		return
	}

	c.JSON(http.StatusOK, pkg.ServicesResponse{Services: services, Count: len(services)})
}

// GetServices godoc
// @Summary Listar serviços
// @Description Retorna lista de serviços com filtros
// @Tags Serviços
// @Accept json
// @Produce json
// @Param status query string false "Status do serviço"
// @Param proposal_letter_id query int false "ID da proposta"
// @Param title query string false "Título parcial"
// @Param content query string false "Conteúdo parcial"
// @Param min_hourly_rate query number false "Valor mínimo por hora"
// @Param max_hourly_rate query number false "Valor máximo por hora"
// @Param sort_by query string false "Campo para ordenar: title,hourly_rate,status,time_created"
// @Param order query string false "Direção: ASC ou DESC"
// @Success 200 {object} pkg.ServicesResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /services [get]
func GetServices(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filters pkg.ServiceFilter
		err := c.ShouldBindQuery(&filters)
		if pkg.HandleErr(c, err) {
			return
		}
		filters.Normalize()

		scanServices(c, conn, serviceBuilder(filters))
	}
}

// GetService godoc
// @Summary Obter serviço
// @Description Retorna um serviço pelo ID
// @Tags Serviços
// @Accept json
// @Produce json
// @Param id path int true "ID do serviço"
// @Success 200 {object} pkg.ServiceResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /services/{id} [get]
func GetService(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		query, args, err := squirrel.Select(serviceSelect).
			From("service").Where(squirrel.Eq{"id": id}).
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		service, err := pkg.ScanService(conn.QueryRow(c.Request.Context(), query, args...))
		if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusOK, pkg.ServiceResponse{Service: service})
	}
}

// CreateService godoc
// @Summary Criar serviço
// @Description Cria um novo serviço associado a uma proposta
// @Tags Serviços
// @Accept json
// @Produce json
// @Param service body pkg.Service true "Objeto serviço"
// @Success 201 {object} pkg.ServiceResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 409 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /services [post]

func CreateService(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var service pkg.Service
		if err := c.BindJSON(&service); pkg.HandleErr(c, err) {
			return
		}

		query, args, err := squirrel.Insert("service").
			Columns("proposal_letter_id", "title", "content", "hourly_rate", "status").
			Values(service.Proposal_letter_id, service.Title, service.Content, service.Hourly_rate, service.Status).
			Suffix("RETURNING " + serviceSelect).
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		service, err = pkg.ScanService(conn.QueryRow(c.Request.Context(), query, args...))
		if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusCreated, pkg.ServiceResponse{Service: service})
	}
}

// UpdateService godoc
// @Summary Atualizar serviço
// @Description Atualiza um serviço pelo ID
// @Tags Serviços
// @Accept json
// @Produce json
// @Param id path int true "ID do serviço"
// @Param service body pkg.Service true "Objeto serviço"
// @Success 200 {object} pkg.ServiceResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /services/{id} [put]
func UpdateService(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		var service pkg.Service
		if err := c.BindJSON(&service); pkg.HandleErr(c, err) {
			return
		}

		if service.Title == "" {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Received empty title"))
			return
		}
		if service.Hourly_rate < 0 {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Received invalid hourly rate. Must be positive"))
			return
		}

		query, args, err := squirrel.Update("service").
			Set("proposal_letter_id", service.Proposal_letter_id).
			Set("title", service.Title).
			Set("content", service.Content).
			Set("hourly_rate", service.Hourly_rate).
			Set("status", service.Status).
			Where(squirrel.Eq{"id": id}).
			Suffix("RETURNING " + serviceSelect).
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		service, err = pkg.ScanService(conn.QueryRow(c.Request.Context(), query, args...))
		if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusOK, pkg.ServiceResponse{Service: service})
	}
}

// UpdateServicePartial godoc
// @Summary Atualizar parcialmente um serviço
// @Description Atualiza um ou mais campos do serviço pelo ID
// @Tags Serviços
// @Accept json
// @Produce json
// @Param id path int true "ID do serviço"
// @Param service body pkg.ServiceUpdateRequest true "Objeto serviço"
// @Success 200 {object} pkg.ServiceResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /services/{id} [patch]
func UpdateServicePartial(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		var req pkg.ServiceUpdateRequest
		if err := c.BindJSON(&req); pkg.HandleErr(c, err) {
			return
		}

		builder := squirrel.Update("service").PlaceholderFormat(squirrel.Dollar)
		hasFields := false

		if req.Title != nil {
			builder = builder.Set("title", *req.Title)
			hasFields = true
		}
		if req.Content != nil {
			builder = builder.Set("content", *req.Content)
			hasFields = true
		}
		if req.Hourly_rate != nil {
			if *req.Hourly_rate < 0 {
				c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Received invalid hourly rate. Must be positive"))
				return
			}
			builder = builder.Set("hourly_rate", *req.Hourly_rate)
			hasFields = true
		}
		if req.Status != nil {
			builder = builder.Set("status", *req.Status)
			hasFields = true
		}

		if !hasFields {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "No valid fields were given"))
			return
		}

		query, args, err := builder.
			Where(squirrel.Eq{"id": id}).
			Suffix("RETURNING " + serviceSelect).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		service, err := pkg.ScanService(conn.QueryRow(c.Request.Context(), query, args...))
		if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusOK, pkg.ServiceResponse{Service: service})
	}
}

// DeleteService godoc
// @Summary Excluir serviço
// @Description Remove um serviço pelo ID
// @Tags Serviços
// @Accept json
// @Produce json
// @Param id path int true "ID do serviço"
// @Success 204 "Deleção com sucesso"
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /services/{id} [delete]
func DeleteService(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		query, args, err := squirrel.Delete("service").
			Where(squirrel.Eq{"id": id}).
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		result, err := conn.Exec(c.Request.Context(), query, args...)
		if pkg.HandleErr(c, err) {
			return
		}
		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Service not found"))
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// GetClientServices de um cliente godoc
// @Summary Listar serviços
// @Description Retorna lista de serviços relacionados a um cliente com filtros
// @Tags Serviços
// @Accept json
// @Produce json
// @Param id path int true "ID do cliente"
// @Param status query string false "Status do serviço"
// @Param proposal_letter_id query int false "ID da proposta"
// @Param title query string false "Título parcial"
// @Param content query string false "Conteúdo parcial"
// @Param min_hourly_rate query number false "Valor mínimo por hora"
// @Param max_hourly_rate query number false "Valor máximo por hora"
// @Param sort_by query string false "Campo para ordenar: title,hourly_rate,status,time_created"
// @Param order query string false "Direção: ASC ou DESC"
// @Success 200 {object} pkg.ServicesResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /clients/{id}/services [get]
func GetClientServices(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		var filters pkg.ServiceFilter
		err = c.ShouldBindQuery(&filters)
		if pkg.HandleErr(c, err) {
			return
		}
		filters.Normalize()

		builder := serviceBuilder(filters).
			From("service s").
			Join("proposal_letter p ON s.proposal_letter_id = p.id").
			Join("client cl ON p.client_id = cl.id").
			Where(squirrel.Eq{"cl.id": id})

		scanServices(c, conn, builder)
	}
}

// GetAnalystServices godoc
// @Summary Listar serviços de um analista
// @Description Retorna lista de serviços relacionados a um analista com filtros
// @Tags Serviços
// @Accept json
// @Produce json
// @Param id path int true "ID do analista"
// @Param status query string false "Status do serviço"
// @Param proposal_letter_id query int false "ID da proposta"
// @Param title query string false "Título parcial"
// @Param content query string false "Conteúdo parcial"
// @Param min_hourly_rate query number false "Valor mínimo por hora"
// @Param max_hourly_rate query number false "Valor máximo por hora"
// @Param sort_by query string false "Campo para ordenar: title,hourly_rate,status,time_created"
// @Param order query string false "Direção: ASC ou DESC"
// @Success 200 {object} pkg.ServicesResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /analysts/{id}/services [get]
func GetAnalystServices(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		var filters pkg.ServiceFilter
		err = c.ShouldBindQuery(&filters)
		if pkg.HandleErr(c, err) {
			return
		}
		filters.Normalize()

		builder := serviceBuilder(filters).
			From("service s").
			Join("proposal_letter p ON s.proposal_letter_id = p.id").
			Join("analyst a ON p.analyst_id = a.id").
			Where(squirrel.Eq{"a.id": id})

		scanServices(c, conn, builder)
	}
}
