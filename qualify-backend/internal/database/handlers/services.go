package handlers

import (
	"context"
	"fmt"
	"main/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func GetServices(conn *pgx.Conn) gin.HandlerFunc {
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
			query += fmt.Sprintf(" AND proposal_letter_id = $%d", argCounter)
			args = append(args, proposalID)
			argCounter++
		}

		query += " ORDER BY time_created DESC"

		rows, err := conn.Query(context.Background(), query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar serviços: " + err.Error()})
			return
		}
		defer rows.Close()

		var services []pkg.Service
		for rows.Next() {
			var s pkg.Service
			if err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Proposal_letter_id,
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
		c.JSON(http.StatusOK, gin.H{"services": services, "count": len(services)})
	}
}

func GetService(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var s pkg.Service
		err := conn.QueryRow(context.Background(),
			`SELECT id, title, content, proposal_letter_id, hourly_rate, status, time_created
			 FROM service WHERE id = $1`, id,
		).Scan(&s.ID, &s.Title, &s.Content, &s.Proposal_letter_id,
			&s.Hourly_rate, &s.Status, &s.Time_created)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar serviço: " + err.Error()})
			return
		}

		rows, err := conn.Query(context.Background(),
			`SELECT id, service_id, rating, comment, time_created
			 FROM review WHERE service_id = $1 ORDER BY time_created DESC`, id,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar reviews: " + err.Error()})
			return
		}
		defer rows.Close()

		var reviews []pkg.Review
		for rows.Next() {
			var r pkg.Review
			if err := rows.Scan(&r.ID, &r.Service_id, &r.Rating, &r.Comment, &r.Time_created); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao escanear review: " + err.Error()})
				return
			}
			reviews = append(reviews, r)
		}
		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao iterar reviews: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"service": s, "reviews": reviews})
	}
}
