package handlers

import (
	"context"
	"main/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func GetCertifications(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := conn.Query(context.Background(),
			`SELECT id, name, year, description FROM certification ORDER BY year DESC`,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar certificações: " + err.Error()})
			return
		}
		defer rows.Close()

		var certs []pkg.Certification
		for rows.Next() {
			var cert pkg.Certification
			if err := rows.Scan(&cert.Id, &cert.Name, &cert.Year, &cert.Description); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao escanear certificação: " + err.Error()})
				return
			}
			certs = append(certs, cert)
		}
		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao iterar certificações: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"certifications": certs, "count": len(certs)})
	}
}

func GetCertification(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var cert pkg.Certification
		err := conn.QueryRow(context.Background(),
			`SELECT id, name, year, description FROM certification WHERE id = $1`, id,
		).Scan(&cert.Id, &cert.Name, &cert.Year, &cert.Description)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar certificação: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, cert)
	}
}
