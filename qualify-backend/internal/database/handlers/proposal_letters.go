package handlers

import (
	"errors"
	"fmt"
	"main/pkg"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GetProposalLetters godoc
// @Summary Listar propostas
// @Description Retorna lista de cartas de proposta (proposals)
// @Tags Propostas
// @Accept json
// @Produce json
// @Param client_id query int false "ID do cliente"
// @Param analyst_id query int false "ID do analista"
// @Param title query string false "Título parcial"
// @Param content query string false "Conteúdo parcial"
// @Param min_proposed_hourly_rate query number false "Valor mínimo por hora proposto"
// @Param max_proposed_hourly_rate query number false "Valor máximo por hora proposto"
// @Param sort_by query string false "Campo para ordenar: title,proposed_hourly_rate,time_created"
// @Param order query string false "Direção: ASC ou DESC"
// @Success 200 {object} pkg.ProposalLettersResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /proposals [get]
func GetProposalLetters(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `SELECT id, client_id, analyst_id, proposed_hourly_rate,
                         title, content, time_created
                  FROM proposal_letter WHERE 1=1`
		args := []interface{}{}
		argCounter := 1

		if clientID := c.Query("client_id"); clientID != "" {
			if clientIDVal, err := strconv.Atoi(clientID); err == nil {
				query += fmt.Sprintf(" AND client_id = $%d", argCounter)
				args = append(args, clientIDVal)
				argCounter++
			}
		}
		if analystID := c.Query("analyst_id"); analystID != "" {
			if analystIDVal, err := strconv.Atoi(analystID); err == nil {
				query += fmt.Sprintf(" AND analyst_id = $%d", argCounter)
				args = append(args, analystIDVal)
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

		if minRate := c.Query("min_proposed_hourly_rate"); minRate != "" {
			if minRateVal, err := strconv.ParseFloat(minRate, 64); err == nil {
				query += fmt.Sprintf(" AND proposed_hourly_rate >= $%d", argCounter)
				args = append(args, minRateVal)
				argCounter++
			}
		}

		if maxRate := c.Query("max_proposed_hourly_rate"); maxRate != "" {
			if maxRateVal, err := strconv.ParseFloat(maxRate, 64); err == nil {
				query += fmt.Sprintf(" AND proposed_hourly_rate <= $%d", argCounter)
				args = append(args, maxRateVal)
				argCounter++
			}
		}

		allowedSortFields := map[string]bool{
			"title": true, "proposed_hourly_rate": true, "time_created": true,
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

		var proposals []pkg.ProposalLetter
		for rows.Next() {
			p, err := pkg.ScanProposalLetter(rows)
			if err != nil {
				c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
				return
			}
			proposals = append(proposals, p)
		}

		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}

		c.JSON(http.StatusOK, pkg.ProposalLettersResponse{Proposal_letters: proposals, Count: len(proposals)})
	}
}

// GetProposalLetter godoc
// @Summary Obter proposta
// @Description Retorna uma proposta específica pelo ID
// @Tags Propostas
// @Accept json
// @Produce json
// @Param id path int true "ID da proposta"
// @Success 200 {object} pkg.ProposalLetterResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /proposals/{id} [get]
func GetProposalLetter(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		proposalID, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		row := conn.QueryRow(c.Request.Context(),
			`SELECT id, client_id, analyst_id, proposed_hourly_rate,
                    title, content, time_created
             FROM proposal_letter WHERE id = $1`, proposalID,
		)

		p, err := pkg.ScanProposalLetter(row)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), err.Error()))
				return
			}
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}

		c.JSON(http.StatusOK, pkg.ProposalLetterResponse{Proposal_letter: p})
	}
}

// CreateProposalLetter godoc
// @Summary Criar proposta
// @Description Cria uma nova carta de proposta
// @Tags Propostas
// @Accept json
// @Produce json
// @Param proposal body pkg.ProposalLetter true "Objeto proposta"
// @Success 201 {object} pkg.ProposalLetterResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 409 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /proposals [post]
func CreateProposalLetter(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var proposal pkg.ProposalLetter
		if err := c.BindJSON(&proposal); err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		err := conn.QueryRow(c.Request.Context(),
			`INSERT INTO proposal_letter (title, content, client_id, analyst_id, proposed_hourly_rate)
            VALUES ($1, $2, $3, $4, $5)
            RETURNING id, client_id, analyst_id, proposed_hourly_rate, title, content, time_created`,
			proposal.Title, proposal.Content, proposal.Client_id, proposal.Analyst_id, proposal.Proposed_hourly_rate).
			Scan(&proposal.Id, &proposal.Client_id, &proposal.Analyst_id, &proposal.Proposed_hourly_rate, &proposal.Title, &proposal.Content, &proposal.Time_created)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				switch pgErr.Code {
				case "23505": // unique_violation
					c.JSON(http.StatusConflict, pkg.Conflict(c.FullPath(), "Proposal letter already exists"))
					return
				case "23503": // foreign_key_violation (client_id ou analyst_id não existe)
					c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Client or analyst not found"))
					return
				}
			}
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}

		row := conn.QueryRow(c.Request.Context(),
			`INSERT INTO proposal_letter (title, content, client_id, analyst_id, proposed_hourly_rate)
             VALUES ($1, $2, $3, $4, $5)
             RETURNING id, client_id, analyst_id, proposed_hourly_rate, title, content, time_created`,
			proposal.Title, proposal.Content, proposal.Client_id, proposal.Analyst_id, proposal.Proposed_hourly_rate)

		proposal, err = pkg.ScanProposalLetter(row)
		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}

		c.JSON(http.StatusCreated, pkg.ProposalLetterResponse{Proposal_letter: proposal})
	}
}

// UpdateProposalLetter godoc
// @Summary Atualizar proposta
// @Description Atualiza uma proposta existente pelo ID
// @Tags Propostas
// @Accept json
// @Produce json
// @Param id path int true "ID da proposta"
// @Param proposal body pkg.ProposalLetter true "Objeto proposta"
// @Success 200 {object} pkg.ProposalLetterResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /proposals/{id} [put]
func UpdateProposalLetter(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		proposalID, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}
		var proposal pkg.ProposalLetter
		if err := c.BindJSON(&proposal); err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		// Validando parâmetros obrigatórios
		if proposal.Title == "" {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Received empty title"))
			return
		}

		row := conn.QueryRow(c.Request.Context(),
			`UPDATE proposal_letter SET title = $1, content = $2, client_id = $3, analyst_id = $4, proposed_hourly_rate = $5
             WHERE id = $6
             RETURNING id, client_id, analyst_id, proposed_hourly_rate, title, content, time_created`,
			proposal.Title, proposal.Content, proposal.Client_id, proposal.Analyst_id, proposal.Proposed_hourly_rate, proposalID)

		proposal, err = pkg.ScanProposalLetter(row)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), err.Error()))
				return
			}
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}

		c.JSON(http.StatusOK, pkg.ProposalLetterResponse{Proposal_letter: proposal})
	}
}

// UpdateProposalLetterPartial godoc
// @Summary Atualizar parcialmente proposta
// @Description Atualiza um ou mais campos da proposta existente pelo ID
// @Tags Propostas
// @Accept json
// @Produce json
// @Param id path int true "ID da proposta"
// @Param proposal body pkg.ProposalLetterUpdateRequest true "Objeto proposta"
// @Success 200 {object} pkg.ProposalLetterResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /proposals/{id} [patch]
func UpdateProposalLetterPartial(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		proposalID, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		var req pkg.ProposalLetterUpdateRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		set := []string{}
		args := []interface{}{}
		argID := 1

		if req.Title != nil {
			set = append(set, fmt.Sprintf("title = $%d", argID))
			args = append(args, *req.Title)
			argID++
		}
		if req.Content != nil {
			set = append(set, fmt.Sprintf("content = $%d", argID))
			args = append(args, *req.Content)
			argID++
		}
		if req.Proposed_hourly_rate != nil {
			set = append(set, fmt.Sprintf("proposed_hourly_rate = $%d", argID))
			args = append(args, *req.Proposed_hourly_rate)
			argID++
		}

		if len(set) == 0 {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Invalid arguments"))
			return
		}

		// Adiciona o ID como último argumento
		args = append(args, proposalID)
		query := fmt.Sprintf(
			"UPDATE proposal_letter SET %s WHERE id = $%d RETURNING id, client_id, analyst_id, proposed_hourly_rate, title, content, time_created",
			strings.Join(set, ", "), argID,
		)

		row := conn.QueryRow(c.Request.Context(), query, args...)
		p, err := pkg.ScanProposalLetter(row)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), err.Error()))
				return
			}
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}

		c.JSON(http.StatusOK, gin.H{"proposal_letter": p})
	}
}

// DeleteProposalLetter godoc
// @Summary Excluir proposta
// @Description Remove uma proposta pelo ID
// @Tags Propostas
// @Accept json
// @Produce json
// @Param id path int true "ID da proposta"
// @Success 204 "Deleção com sucesso"
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /proposals/{id} [delete]
func DeleteProposalLetter(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		proposalID, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		result, err := conn.Exec(c.Request.Context(),
			`DELETE FROM proposal_letter WHERE id = $1`, proposalID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Proposal letter not found"))
			return
		}

		c.Status(http.StatusNoContent)
	}
}

// GetClientProposalLetters godoc
// @Summary Listar propostas de cliente
// @Description Retorna lista de cartas de proposta (proposals) relacionadas a um cliente específico
// @Tags Propostas
// @Accept json
// @Produce json
// @Param id path int true "ID do cliente"
// @Param analyst_id query int false "ID do analista"
// @Param title query string false "Título parcial"
// @Param content query string false "Conteúdo parcial"
// @Param min_proposed_hourly_rate query number false "Valor mínimo por hora proposto"
// @Param max_proposed_hourly_rate query number false "Valor máximo por hora proposto"
// @Param sort_by query string false "Campo para ordenar: title,proposed_hourly_rate,time_created"
// @Param order query string false "Direção: ASC ou DESC"
// @Success 200 {object} pkg.ProposalLettersResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /clients/{id}/proposals [get]
func GetClientProposalLetters(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		query := `SELECT p.id, p.client_id, p.analyst_id, p.proposed_hourly_rate,
          p.title, p.content, p.time_created
          FROM proposal_letter p
          JOIN client c ON p.client_id = c.id
          WHERE c.user_id = $1`
		args := []interface{}{clientID}
		argCounter := 2

		if analystID := c.Query("analyst_id"); analystID != "" {
			if analystIDVal, err := strconv.Atoi(analystID); err == nil {
				query += fmt.Sprintf(" AND analyst_id = $%d", argCounter)
				args = append(args, analystIDVal)
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

		if minRate := c.Query("min_proposed_hourly_rate"); minRate != "" {
			if minRateVal, err := strconv.ParseFloat(minRate, 64); err == nil {
				query += fmt.Sprintf(" AND proposed_hourly_rate >= $%d", argCounter)
				args = append(args, minRateVal)
				argCounter++
			}
		}

		if maxRate := c.Query("max_proposed_hourly_rate"); maxRate != "" {
			if maxRateVal, err := strconv.ParseFloat(maxRate, 64); err == nil {
				query += fmt.Sprintf(" AND proposed_hourly_rate <= $%d", argCounter)
				args = append(args, maxRateVal)
				argCounter++
			}
		}

		allowedSortFields := map[string]bool{
			"title": true, "proposed_hourly_rate": true, "time_created": true,
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

		var proposals []pkg.ProposalLetter
		for rows.Next() {
			p, err := pkg.ScanProposalLetter(rows)
			if err != nil {
				c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
				return
			}
			proposals = append(proposals, p)
		}

		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}

		c.JSON(http.StatusOK, pkg.ProposalLettersResponse{Proposal_letters: proposals, Count: len(proposals)})
	}
}

// GetAnalystProposalLetters godoc
// @Summary Listar propostas de analista
// @Description Retorna lista de cartas de proposta (proposals) relacionadas a um analista específico
// @Tags Propostas
// @Accept json
// @Produce json
// @Param id path int true "ID do analista"
// @Param client_id query int false "ID do cliente"
// @Param title query string false "Título parcial"
// @Param content query string false "Conteúdo parcial"
// @Param min_proposed_hourly_rate query number false "Valor mínimo por hora proposto"
// @Param max_proposed_hourly_rate query number false "Valor máximo por hora proposto"
// @Param sort_by query string false "Campo para ordenar: title,proposed_hourly_rate,time_created"
// @Param order query string false "Direção: ASC ou DESC"
// @Success 200 {object} pkg.ProposalLettersResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /analysts/{id}/proposals [get]
func GetAnalystProposalLetters(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		analystID, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		query := `SELECT p.id, p.client_id, p.analyst_id, p.proposed_hourly_rate,
          p.title, p.content, p.time_created
          FROM proposal_letter p
          JOIN analyst a ON p.analyst_id = a.id
          WHERE a.user_id = $1`
		args := []interface{}{analystID}
		argCounter := 2

		if clientID := c.Query("client_id"); clientID != "" {
			if clientIDVal, err := strconv.Atoi(clientID); err == nil {
				query += fmt.Sprintf(" AND client_id = $%d", argCounter)
				args = append(args, clientIDVal)
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

		if minRate := c.Query("min_proposed_hourly_rate"); minRate != "" {
			if minRateVal, err := strconv.ParseFloat(minRate, 64); err == nil {
				query += fmt.Sprintf(" AND proposed_hourly_rate >= $%d", argCounter)
				args = append(args, minRateVal)
				argCounter++
			}
		}

		if maxRate := c.Query("max_proposed_hourly_rate"); maxRate != "" {
			if maxRateVal, err := strconv.ParseFloat(maxRate, 64); err == nil {
				query += fmt.Sprintf(" AND proposed_hourly_rate <= $%d", argCounter)
				args = append(args, maxRateVal)
				argCounter++
			}
		}

		allowedSortFields := map[string]bool{
			"title": true, "proposed_hourly_rate": true, "time_created": true,
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

		var proposals []pkg.ProposalLetter
		for rows.Next() {
			p, err := pkg.ScanProposalLetter(rows)
			if err != nil {
				c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
				return
			}
			proposals = append(proposals, p)
		}

		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}

		c.JSON(http.StatusOK, pkg.ProposalLettersResponse{Proposal_letters: proposals, Count: len(proposals)})
	}
}
