package handlers

import (
	"context"
	"fmt"
	"main/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func GetProposalLetters(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `SELECT id, client_id, analyst_id, proposed_hourly_rate,
		                 title, content, time_created
		          FROM proposal_letter WHERE 1=1`
		args := []interface{}{}
		argCounter := 1

		if clientID := c.Query("client_id"); clientID != "" {
			query += fmt.Sprintf(" AND client_id = $%d", argCounter)
			args = append(args, clientID)
			argCounter++
		}
		if analystID := c.Query("analyst_id"); analystID != "" {
			query += fmt.Sprintf(" AND analyst_id = $%d", argCounter)
			args = append(args, analystID)
			argCounter++
		}

		query += " ORDER BY time_created DESC"

		rows, err := conn.Query(context.Background(), query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar propostas: " + err.Error()})
			return
		}
		defer rows.Close()

		var proposals []pkg.ProposalLetter
		for rows.Next() {
			var p pkg.ProposalLetter
			if err := rows.Scan(&p.Id, &p.Client_id, &p.Analyst_id, &p.Proposed_hourly_rate,
				&p.Title, &p.Content, &p.Time_created); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao escanear proposta: " + err.Error()})
				return
			}
			proposals = append(proposals, p)
		}
		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao iterar propostas: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"proposals": proposals, "count": len(proposals)})
	}
}

// GetProposalLetter retorna a proposta e os serviços vinculados a ela.
func GetProposalLetter(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var p pkg.ProposalLetter
		err := conn.QueryRow(context.Background(),
			`SELECT id, client_id, analyst_id, proposed_hourly_rate,
			        title, content, time_created
			 FROM proposal_letter WHERE id = $1`, id,
		).Scan(&p.Id, &p.Client_id, &p.Analyst_id, &p.Proposed_hourly_rate,
			&p.Title, &p.Content, &p.Time_created)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar proposta: " + err.Error()})
			return
		}

		rows, err := conn.Query(context.Background(),
			`SELECT id, title, content, proposal_letter_id, hourly_rate, status, time_created
			 FROM service WHERE proposal_letter_id = $1 ORDER BY time_created`, id,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar serviços: " + err.Error()})
			return
		}
		defer rows.Close()

		var services []pkg.Service
		for rows.Next() {
			var s pkg.Service
			if err := rows.Scan(&s.Id, &s.Title, &s.Content, &s.Proposal_letter_id,
				&s.Hourly_rate, &s.Status, &s.Time_created); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao escanear serviço: " + err.Error()})
				return
			}
			services = append(services, s)
		}
		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao iterar serviços: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"proposal": p, "services": services})
	}
}
