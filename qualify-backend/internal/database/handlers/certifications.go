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

// GetCertifications godoc
// @Summary Listar certificações
// @Description Retorna lista de certificações
// @Tags Certificações
// @Accept json
// @Produce json
// @Param name query string false "Nome parcial"
// @Success 200 {object} pkg.CertificationsResponse
// @Failure 500 pkg.ErrorResponse
// @Router /certifications [get]
func GetCertifications(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `SELECT id, name, year, description, institution FROM certification WHERE 1=1`
		args := []interface{}{}
		argCounter := 1

		if name := c.Query("name"); name != "" {
			query += fmt.Sprintf(" AND name ILIKE $%d", argCounter)
			args = append(args, "%"+name+"%")
			argCounter++
		}

		query += " ORDER BY year DESC"

		rows, err := conn.Query(c.Request.Context(), query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}
		defer rows.Close()

		var certs []pkg.Certification
		for rows.Next() {
			var cert pkg.Certification
			if err := rows.Scan(&cert.Id, &cert.Name, &cert.Year, &cert.Description, &cert.Institution); err != nil {
				c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
				return
			}
			certs = append(certs, cert)
		}

		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
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
// @Failure 400 pkg.ErrorResponse
// @Failure 404 pkg.ErrorResponse
// @Failure 500 pkg.ErrorResponse
// @Router /certifications/{id} [get]
func GetCertification(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		certID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		var cert pkg.Certification
		err = conn.QueryRow(c.Request.Context(),
			`SELECT id, name, year, description, institution FROM certification WHERE id = $1`, certID,
		).Scan(&cert.Id, &cert.Name, &cert.Year, &cert.Description, &cert.Institution)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), err.Error()))
				return
			}
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
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
// @Failure 400 pkg.ErrorResponse
// @Failure 500 pkg.ErrorResponse
// @Security     BearerAuth
// @Router /certifications [post]
func CreateCertification(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var cert pkg.Certification
		if err := c.BindJSON(&cert); err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		// Validando parâmetros obrigatórios
		if cert.Name == "" {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Received empty name"))
			return
		}

		if cert.Description == "" {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Received empty description"))
			return
		}

		if cert.Institution == "" {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Received empty institution"))
			return
		}

		if cert.Year < 1900 || cert.Year > 2030 {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Received invalid year. Must be between 1900 and 2030"))
			return
		}

		err := conn.QueryRow(c.Request.Context(),
			`INSERT INTO certification (name, year, description, institution) VALUES ($1, $2, $3, $4) RETURNING id`,
			cert.Name, cert.Year, cert.Description, cert.Institution,
		).Scan(&cert.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
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
// @Failure 400 pkg.ErrorResponse
// @Failure 404 pkg.ErrorResponse
// @Failure 500 pkg.ErrorResponse
// @Security     BearerAuth
// @Router /certifications/{id} [put]
func UpdateCertification(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		certID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		var cert pkg.Certification
		if err := c.BindJSON(&cert); err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		// Validando parâmetros obrigatórios
		if cert.Name == "" {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Received empty name"))
			return
		}

		if cert.Description == "" {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Received empty description"))
			return
		}

		if cert.Institution == "" {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Received empty institution"))
			return
		}

		if cert.Year < 1900 || cert.Year > 2030 {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Received invalid year. Must be between 1900 and 2030"))
			return
		}

		err = conn.QueryRow(c.Request.Context(),
			`UPDATE certification SET name = $1, year = $2, description = $3, institution = $4 WHERE id = $5 RETURNING id, name, year, description, institution`,
			cert.Name, cert.Year, cert.Description, cert.Institution, certID,
		).Scan(&cert.Id, &cert.Name, &cert.Year, &cert.Description, &cert.Institution)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), err.Error()))
				return
			}
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}
		c.JSON(http.StatusOK, pkg.CertificationResponse{
			Certification: cert,
		})
	}
}

// UpdateCertificationPartial godoc
// @Summary Atualizar certificação parcialmente
// @Description Atualiza um ou mais campos da certificação pelo ID
// @Tags Certificações
// @Accept json
// @Produce json
// @Param id path int true "ID da certificação"
// @Param certification body pkg.CertificationUpdateRequest true "Objeto certificação"
// @Success 200 {object} pkg.CertificationResponse
// @Failure 400 pkg.ErrorResponse
// @Failure 404 pkg.ErrorResponse
// @Failure 500 pkg.ErrorResponse
// @Security     BearerAuth
// @Router /certifications/{id} [patch]
func UpdateCertificationPartial(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		certID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		var cert pkg.CertificationUpdateRequest
		if err := c.BindJSON(&cert); err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		set := []string{}
		args := []interface{}{}
		i := 1

		if cert.Name != nil {
			set = append(set, fmt.Sprintf("name = $%d", i))
			args = append(args, *cert.Name)
			i++
		}
		if cert.Year != nil {
			if *cert.Year < 1900 || *cert.Year > 2030 {
				c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Received invalid year. Must be between 1900 and 2030"))
				return
			}
			set = append(set, fmt.Sprintf("year = $%d", i))
			args = append(args, *cert.Year)
			i++
		}
		if cert.Description != nil {
			set = append(set, fmt.Sprintf("description = $%d", i))
			args = append(args, *cert.Description)
			i++
		}
		if cert.Institution != nil {
			set = append(set, fmt.Sprintf("institution = $%d", i))
			args = append(args, *cert.Institution)
			i++
		}

		if len(set) == 0 {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Invalid arguments"))
			return
		}

		args = append(args, certID)

		query := fmt.Sprintf("UPDATE certification SET %s WHERE id = $%d RETURNING id, name, year, description, institution",
			strings.Join(set, ", "), i)

		var updatedCert pkg.Certification
		err = conn.QueryRow(c.Request.Context(), query, args...).Scan(&updatedCert.Id, &updatedCert.Name, &updatedCert.Year, &updatedCert.Description, &updatedCert.Institution)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), err.Error()))
				return
			}
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}

		c.JSON(http.StatusOK, pkg.CertificationResponse{Certification: updatedCert})
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
// @Failure 400 pkg.ErrorResponse
// @Failure 404 pkg.ErrorResponse
// @Failure 500 pkg.ErrorResponse
// @Security     BearerAuth
// @Router /certifications/{id} [delete]
func DeleteCertification(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		certID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		result, err := conn.Exec(c.Request.Context(),
			`DELETE FROM certification WHERE id = $1`, certID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Certification not found"))
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
// @Failure 400 pkg.ErrorResponse
// @Failure 500 pkg.ErrorResponse
// @Router /users/{id}/analyst/certifications [get]
func GetAnalystCertifications(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		rows, err := conn.Query(c.Request.Context(),
			`SELECT c.id, c.name, c.year, c.description, c.institution 
			 FROM certification c 
			 JOIN analyst_certification ac ON c.id = ac.certification_id 
			 WHERE ac.analyst_id = $1 
			 ORDER BY c.year DESC`,
			userID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}
		defer rows.Close()

		var certs []pkg.Certification
		for rows.Next() {
			var cert pkg.Certification
			if err := rows.Scan(&cert.Id, &cert.Name, &cert.Year, &cert.Description, &cert.Institution); err != nil {
				c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
				return
			}
			certs = append(certs, cert)
		}
		if err = rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
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
// @Failure 400 pkg.ErrorResponse
// @Failure 500 pkg.ErrorResponse
// @Security     BearerAuth
// @Router /users/{id}/analyst/certifications [post]
func CreateAnalystCertification(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		analystID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		var ac pkg.AnalystCertification
		if err := c.BindJSON(&ac); err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}
		ac.Analyst_id = analystID

		// Validando parâmetros obrigatórios
		var analystExists bool
		err = conn.QueryRow(c.Request.Context(),
			`SELECT EXISTS(SELECT 1 FROM analyst WHERE id = $1)`, ac.Analyst_id,
		).Scan(&analystExists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}
		if !analystExists {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Analyst does not exists"))
			return
		}

		// Validando parâmetros obrigatórios
		var certExists bool
		err = conn.QueryRow(c.Request.Context(),
			`SELECT EXISTS(SELECT 1 FROM certification WHERE id = $1)`, ac.Certification_id,
		).Scan(&certExists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}
		if !certExists {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Certification does not exists"))
			return
		}

		_, err = conn.Exec(c.Request.Context(),
			`INSERT INTO analyst_certification (analyst_id, certification_id) VALUES ($1, $2)`,
			ac.Analyst_id, ac.Certification_id,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
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
// @Failure 400 pkg.ErrorResponse
// @Failure 404 pkg.ErrorResponse
// @Failure 500 pkg.ErrorResponse
// @Security     BearerAuth
// @Router /users/{id}/analyst/certifications [delete]
func DeleteAnalystCertification(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		userIDVal, err := strconv.Atoi(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		certID := c.Query("certification_id")
		if certID == "" {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Invalid certification ID"))
			return
		}
		certIDVal, err := strconv.Atoi(certID)
		if err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		result, err := conn.Exec(c.Request.Context(),
			`DELETE FROM analyst_certification WHERE analyst_id = $1 AND certification_id = $2`,
			userIDVal, certIDVal,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, pkg.Internal(c.FullPath(), err.Error()))
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Analyst certification not found"))
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "analyst certification deleted successfully"})
	}
}
