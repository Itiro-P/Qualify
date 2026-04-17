package handlers

import (
	"fmt"
	"main/pkg"
	"net/http"
	"strconv"

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
		c.JSON(http.StatusOK, gin.H{"proposals": proposals, "count": len(proposals)})
	}
}

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

		rows, err := conn.Query(c.Request.Context(),
			`SELECT id, title, content, proposal_letter_id, hourly_rate, status, time_created
			 FROM service WHERE proposal_letter_id = $1 ORDER BY time_created`, proposalID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var services []pkg.Service
		for rows.Next() {
			var s pkg.Service
			if err := rows.Scan(&s.Id, &s.Title, &s.Content, &s.Proposal_letter_id,
				&s.Hourly_rate, &s.Status, &s.Time_created); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			services = append(services, s)
		}
		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"proposal": p, "services": services})
	}
}

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

		c.JSON(http.StatusCreated, proposal)
	}
}

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

		c.JSON(http.StatusOK, proposal)
	}
}

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
