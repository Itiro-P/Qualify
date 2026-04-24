package handlers

import (
	"main/pkg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// GetCertifications godoc
// @Summary Listar certificações
// @Description Retorna lista de certificações
// @Tags Certificações
// @Accept json
// @Produce json
// @Success 200 {object} pkg.CertificationsResponse
// @Failure 500 {object} map[string]string
// @Router /certifications [get]
func GetCertifications(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := conn.Query(c.Request.Context(),
			`SELECT id, name, year, description FROM certification ORDER BY year DESC`,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusOK, pkg.CertificationsResponse{
			Certifications: certs,
			Count:          len(certs),
		})
	}
}

// GetCertification godoc
// @Summary Obter certificação
// @Description Retorna certificação por ID
// @Tags Certificações
// @Accept json
// @Produce json
// @Param id path int true "ID da certificação"
// @Success 200 {object} pkg.CertificationResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /certifications/{id} [get]
func GetCertification(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		certID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid certification id"})
			return
		}

		var cert pkg.Certification
		err = conn.QueryRow(c.Request.Context(),
			`SELECT id, name, year, description FROM certification WHERE id = $1`, certID,
		).Scan(&cert.Id, &cert.Name, &cert.Year, &cert.Description)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "certification not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, pkg.CertificationResponse{
			Certification: cert,
		})
	}
}

// CreateCertification godoc
// @Summary Criar certificação
// @Description Cria uma nova certificação
// @Tags Certificações
// @Accept json
// @Produce json
// @Param certification body pkg.Certification true "Objeto certificação"
// @Success 201 {object} pkg.CertificationResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /certifications [post]
func CreateCertification(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		var cert pkg.Certification
		if err := c.BindJSON(&cert); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Validate required fields
		if cert.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "certification name is required"})
			return
		}
		if cert.Year < 1900 || cert.Year > 2030 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "certification year must be between 1900 and 2030"})
			return
		}

		err := conn.QueryRow(c.Request.Context(),
			`INSERT INTO certification (name, year, description) VALUES ($1, $2, $3) RETURNING id`,
			cert.Name, cert.Year, cert.Description,
		).Scan(&cert.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, pkg.CertificationResponse{
			Certification: cert,
		})
	}
}

// UpdateCertification godoc
// @Summary Atualizar certificação
// @Description Atualiza uma certificação pelo ID
// @Tags Certificações
// @Accept json
// @Produce json
// @Param id path int true "ID da certificação"
// @Param certification body pkg.Certification true "Objeto certificação"
// @Success 200 {object} pkg.CertificationResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /certifications/{id} [put]
func UpdateCertification(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		certID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid certification id"})
			return
		}

		var cert pkg.Certification
		if err := c.BindJSON(&cert); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Validate required fields
		if cert.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "certification name is required"})
			return
		}
		if cert.Year < 1900 || cert.Year > 2030 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "certification year must be between 1900 and 2030"})
			return
		}

		err = conn.QueryRow(c.Request.Context(),
			`UPDATE certification SET name = $1, year = $2, description = $3 WHERE id = $4 RETURNING id, name, year, description`,
			cert.Name, cert.Year, cert.Description, certID,
		).Scan(&cert.Id, &cert.Name, &cert.Year, &cert.Description)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "certification not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, pkg.CertificationResponse{
			Certification: cert,
		})
	}
}

// DeleteCertification godoc
// @Summary Excluir certificação
// @Description Exclui certificação pelo ID
// @Tags Certificações
// @Accept json
// @Produce json
// @Param id path int true "ID da certificação"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /certifications/{id} [delete]
func DeleteCertification(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		certID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid certification id"})
			return
		}

		result, err := conn.Exec(c.Request.Context(),
			`DELETE FROM certification WHERE id = $1`, certID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "certification not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "certification deleted successfully"})
	}
}

// GetAnalystCertifications godoc
// @Summary Listar certificações do analista
// @Description Retorna as certificações associadas a um analista
// @Tags Certificações
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário (analista)"
// @Success 200 {object} pkg.CertificationsResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/analyst/certifications [get]
func GetAnalystCertifications(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		rows, err := conn.Query(c.Request.Context(),
			`SELECT c.id, c.name, c.year, c.description 
			 FROM certification c 
			 JOIN analyst_certification ac ON c.id = ac.certification_id 
			 WHERE ac.analyst_id = $1 
			 ORDER BY c.year DESC`,
			userID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusOK, pkg.CertificationsResponse{
			Certifications: certs,
			Count:          len(certs),
		})
	}
}

// CreateAnalystCertification godoc
// @Summary Associar certificação ao analista
// @Description Associa uma certificação a um analista
// @Tags Certificações
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário (analista)"
// @Param certification body pkg.AnalystCertification true "Objeto associação"
// @Success 201 {object} pkg.AnalystCertificationResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/analyst/certifications [post]
func CreateAnalystCertification(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		analystID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		var ac pkg.AnalystCertification
		if err := c.BindJSON(&ac); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		ac.Analyst_id = analystID

		// Validate that analyst exists
		var analystExists bool
		err = conn.QueryRow(c.Request.Context(),
			`SELECT EXISTS(SELECT 1 FROM analyst WHERE id = $1)`, ac.Analyst_id,
		).Scan(&analystExists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !analystExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "analyst not found"})
			return
		}

		// Validate that certification exists
		var certExists bool
		err = conn.QueryRow(c.Request.Context(),
			`SELECT EXISTS(SELECT 1 FROM certification WHERE id = $1)`, ac.Certification_id,
		).Scan(&certExists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !certExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "certification not found"})
			return
		}

		_, err = conn.Exec(c.Request.Context(),
			`INSERT INTO analyst_certification (analyst_id, certification_id) VALUES ($1, $2)`,
			ac.Analyst_id, ac.Certification_id,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, pkg.AnalystCertificationResponse{
			Analyst_certification: ac,
		})
	}
}

// DeleteAnalystCertification godoc
// @Summary Remover certificação do analista
// @Description Remove associação de certificação de um analista
// @Tags Certificações
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário (analista)"
// @Param certification_id query int true "ID da certificação"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/analyst/certifications [delete]
func DeleteAnalystCertification(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		userIDVal, err := strconv.Atoi(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		certID := c.Query("certification_id")
		if certID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "certification_id query parameter required"})
			return
		}
		certIDVal, err := strconv.Atoi(certID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid certification_id"})
			return
		}

		result, err := conn.Exec(c.Request.Context(),
			`DELETE FROM analyst_certification WHERE analyst_id = $1 AND certification_id = $2`,
			userIDVal, certIDVal,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "analyst certification not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "analyst certification deleted successfully"})
	}
}
