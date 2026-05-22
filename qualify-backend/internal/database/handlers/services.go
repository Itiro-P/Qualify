package handlers

import (
	"fmt"
	"main/pkg"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

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
		query := `SELECT id, title, content, proposal_letter_id, hourly_rate, status, time_created
		          FROM service WHERE 1=1`
		args := []interface{}{}
		argCounter := 1

		if status := c.Query("status"); status != "" {
			query += fmt.Sprintf(" AND status = $%d", argCounter)
			args = append(args, status)
			argCounter++
		}
		if proposalID := c.Query("proposal_letter_id"); proposalID != "" {
			if proposalIDVal, err := strconv.Atoi(proposalID); err == nil {
				query += fmt.Sprintf(" AND proposal_letter_id = $%d", argCounter)
				args = append(args, proposalIDVal)
				argCounter++
			}
		}

		if title := c.Query("title"); title != "" {
			query += fmt.Sprintf(" AND title ILIKE $%d", argCounter)
			args = append(args, "%"+title+"%")
			argCounter++
		}

		if content := c.Query("content"); content != "" {
			query += fmt.Sprintf(" AND content ILIKE $%d", argCounter)
			args = append(args, "%"+content+"%")
			argCounter++
		}

		if minRate := c.Query("min_hourly_rate"); minRate != "" {
			if minRateVal, err := strconv.ParseFloat(minRate, 64); err == nil {
				query += fmt.Sprintf(" AND hourly_rate >= $%d", argCounter)
				args = append(args, minRateVal)
				argCounter++
			}
		}

		if maxRate := c.Query("max_hourly_rate"); maxRate != "" {
			if maxRateVal, err := strconv.ParseFloat(maxRate, 64); err == nil {
				query += fmt.Sprintf(" AND hourly_rate <= $%d", argCounter)
				args = append(args, maxRateVal)
				argCounter++
			}
		}

		allowedSortFields := map[string]bool{
			"title": true, "hourly_rate": true, "status": true, "time_created": true,
		}
		if sortBy := c.Query("sort_by"); sortBy != "" {
			if allowedSortFields[sortBy] {
				order := c.DefaultQuery("order", "ASC")
				if order == "ASC" || order == "DESC" {
					query += fmt.Sprintf(" ORDER BY %s %s", sortBy, order)
				}
			}
		} else {
			query += " ORDER BY time_created DESC"
		}

		rows, err := conn.Query(c.Request.Context(), query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}
		defer rows.Close()

		var services []pkg.Service
		for rows.Next() {
			var s pkg.Service
			if err := rows.Scan(&s.Id, &s.Title, &s.Content, &s.Proposal_letter_id,
				&s.Hourly_rate, &s.Status, &s.Time_created); err != nil {
				c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
				return
			}
			services = append(services, s)
		}
		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}
		c.JSON(http.StatusOK, pkg.ServicesResponse{Services: services, Count: len(services)})
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
		id := c.Param("id")
		serviceID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		var s pkg.Service
		err = conn.QueryRow(c.Request.Context(),
			`SELECT id, title, content, proposal_letter_id, hourly_rate, status, time_created
			 FROM service WHERE id = $1`, serviceID,
		).Scan(&s.Id, &s.Title, &s.Content, &s.Proposal_letter_id,
			&s.Hourly_rate, &s.Status, &s.Time_created)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), err.Error()))
				return
			}
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}

		c.JSON(http.StatusOK, pkg.ServiceResponse{Service: s})
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
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /services [post]
func CreateService(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var service pkg.Service
		if err := c.BindJSON(&service); err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		err := conn.QueryRow(c.Request.Context(),
			`INSERT INTO service (proposal_letter_id, title, content, hourly_rate, status)
			 VALUES ($1, $2, $3, $4, $5)
			 RETURNING id, time_created`,
			service.Proposal_letter_id, service.Title, service.Content, service.Hourly_rate, service.Status).
			Scan(&service.Id, &service.Time_created)

		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
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
		id := c.Param("id")
		serviceID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		var service pkg.Service
		if err := c.BindJSON(&service); err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		// Validando parâmetros obrigatórios
		if service.Title == "" {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Received empty title"))
			return
		}
		if service.Hourly_rate < 0 {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Received invalid hourly rate. Must be positive"))
			return
		}

		err = conn.QueryRow(c.Request.Context(),
			`UPDATE service SET proposal_letter_id = $1, title = $2, content = $3, hourly_rate = $4, status = $5
			 WHERE id = $6
			 RETURNING id, proposal_letter_id, title, content, hourly_rate, status, time_created`,
			service.Proposal_letter_id, service.Title, service.Content, service.Hourly_rate, service.Status, serviceID).
			Scan(&service.Id, &service.Proposal_letter_id, &service.Title, &service.Content, &service.Hourly_rate, &service.Status, &service.Time_created)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), err.Error()))
				return
			}
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
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
		id := c.Param("id")
		serviceID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		var service pkg.ServiceUpdateRequest
		if err := c.BindJSON(&service); err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		set := []string{}
		args := []interface{}{}
		i := 1

		if service.Content != nil {
			set = append(set, fmt.Sprintf("content = $%d", i))
			args = append(args, *service.Content)
			i++
		}
		if service.Hourly_rate != nil {
			if *service.Hourly_rate < 0 {
				c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Received invalid hourly rate. Must be positive"))
				return
			}
			set = append(set, fmt.Sprintf("hourly_rate = $%d", i))
			args = append(args, *service.Hourly_rate)
			i++
		}
		if service.Status != nil {
			set = append(set, fmt.Sprintf("status = $%d", i))
			args = append(args, *service.Status)
			i++
		}
		if service.Title != nil {
			set = append(set, fmt.Sprintf("title = $%d", i))
			args = append(args, *service.Title)
			i++
		}

		if len(set) == 0 {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Invalid arguments"))
			return
		}

		args = append(args, serviceID)

		query := fmt.Sprintf("UPDATE service SET %s WHERE id = $%d RETURNING id, proposal_letter_id, title, content, hourly_rate, status, time_created",
			strings.Join(set, ", "), i)

		var updatedService pkg.Service
		err = conn.QueryRow(c.Request.Context(), query, args...).Scan(&updatedService.Id, &updatedService.Proposal_letter_id,
			&updatedService.Title, &updatedService.Content, &updatedService.Hourly_rate,
			&updatedService.Status, &updatedService.Time_created)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), err.Error()))
				return
			}
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}

		c.JSON(http.StatusOK, pkg.ServiceResponse{Service: updatedService})
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
		id := c.Param("id")
		serviceID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		result, err := conn.Exec(c.Request.Context(),
			`DELETE FROM service WHERE id = $1`, serviceID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Service not found"))
			return
		}

		c.Status(http.StatusNoContent)
	}
}
