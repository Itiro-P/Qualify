package handlers

import (
	"fmt"
	"main/pkg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// GetProposalLetters godoc
// @Summary Listar propostas
// @Description Retorna lista de cartas de proposta (proposals)
// @Tags Propostas
// @Accept json
// @Produce json
// @Param client_id query int false "ID do cliente"
// @Param analyst_id query int false "ID do analista"
// @Success 200 {object} pkg.ProposalLettersResponse
// @Failure 500 {object} map[string]string
// @Router /proposals [get]
func GetProposalLetters(conn *pgx.Conn) gin.HandlerFunc {
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

		query += " ORDER BY time_created DESC"

		rows, err := conn.Query(c.Request.Context(), query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var proposals []pkg.ProposalLetter
		for rows.Next() {
			var p pkg.ProposalLetter
			if err := rows.Scan(&p.Id, &p.Client_id, &p.Analyst_id, &p.Proposed_hourly_rate,
				&p.Title, &p.Content, &p.Time_created); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			proposals = append(proposals, p)
		}
		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /proposals/{id} [get]
func GetProposalLetter(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		proposalID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid proposal id"})
			return
		}

		var p pkg.ProposalLetter
		err = conn.QueryRow(c.Request.Context(),
			`SELECT id, client_id, analyst_id, proposed_hourly_rate,
			        title, content, time_created
			 FROM proposal_letter WHERE id = $1`, proposalID,
		).Scan(&p.Id, &p.Client_id, &p.Analyst_id, &p.Proposed_hourly_rate,
			&p.Title, &p.Content, &p.Time_created)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "proposal letter not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /proposals [post]
func CreateProposalLetter(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		var proposal pkg.ProposalLetter
		if err := c.BindJSON(&proposal); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		err := conn.QueryRow(c.Request.Context(),
			`INSERT INTO proposal_letter (title, content, client_id, analyst_id, proposed_hourly_rate)
			 VALUES ($1, $2, $3, $4, $5)
			 RETURNING id, time_created`,
			proposal.Title, proposal.Content, proposal.Client_id, proposal.Analyst_id, proposal.Proposed_hourly_rate).
			Scan(&proposal.Id, &proposal.Time_created)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /proposals/{id} [put]
func UpdateProposalLetter(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		proposalID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid proposal id"})
			return
		}
		var proposal pkg.ProposalLetter
		if err := c.BindJSON(&proposal); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Validando parâmetros obrigatórios
		if proposal.Title == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "proposal title is required"})
			return
		}

		err = conn.QueryRow(c.Request.Context(),
			`UPDATE proposal_letter SET title = $1, content = $2, client_id = $3, analyst_id = $4, proposed_hourly_rate = $5
			 WHERE id = $6
			 RETURNING id, title, content, client_id, analyst_id, proposed_hourly_rate, time_created`,
			proposal.Title, proposal.Content, proposal.Client_id, proposal.Analyst_id, proposal.Proposed_hourly_rate, proposalID).
			Scan(&proposal.Id, &proposal.Title, &proposal.Content, &proposal.Client_id, &proposal.Analyst_id, &proposal.Proposed_hourly_rate, &proposal.Time_created)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "proposal letter not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, pkg.ProposalLetterResponse{Proposal_letter: proposal})
	}
}

// DeleteProposalLetter godoc
// @Summary Excluir proposta
// @Description Remove uma proposta pelo ID
// @Tags Propostas
// @Accept json
// @Produce json
// @Param id path int true "ID da proposta"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /proposals/{id} [delete]
func DeleteProposalLetter(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		proposalID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid proposal id"})
			return
		}

		result, err := conn.Exec(c.Request.Context(),
			`DELETE FROM proposal_letter WHERE id = $1`, proposalID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "proposal letter not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "proposal letter deleted successfully"})
	}
}
