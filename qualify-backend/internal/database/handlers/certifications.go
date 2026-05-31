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
// @Param institution query string false "Instituição parcial"
// @Param year query int false "Ano"
// @Param from_year query int false "Ano inicial"
// @Param to_year query int false "Ano final"
// @Param sort_by query string false "Campo para ordenar: name,year,institution"
// @Param order query string false "Direção: ASC ou DESC"
// @Success 200 {object} pkg.CertificationsResponse
// @Failure 500 {object} pkg.ErrorResponse
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

		if institution := c.Query("institution"); institution != "" {
			query += fmt.Sprintf(" AND institution ILIKE $%d", argCounter)
			args = append(args, "%"+institution+"%")
			argCounter++
		}

		if year := c.Query("year"); year != "" {
			if yearVal, err := strconv.Atoi(year); err == nil {
				query += fmt.Sprintf(" AND year = $%d", argCounter)
				args = append(args, yearVal)
				argCounter++
			}
		}

		if fromYear := c.Query("from_year"); fromYear != "" {
			if fromYearVal, err := strconv.Atoi(fromYear); err == nil {
				query += fmt.Sprintf(" AND year >= $%d", argCounter)
				args = append(args, fromYearVal)
				argCounter++
			}
		}

		if toYear := c.Query("to_year"); toYear != "" {
			if toYearVal, err := strconv.Atoi(toYear); err == nil {
				query += fmt.Sprintf(" AND year <= $%d", argCounter)
				args = append(args, toYearVal)
				argCounter++
			}
		}

		allowedSortFields := map[string]bool{
			"name": true, "year": true, "institution": true,
		}
		if sortBy := c.Query("sort_by"); sortBy != "" {
			if allowedSortFields[sortBy] {
				order := c.DefaultQuery("order", "ASC")
				if order == "ASC" || order == "DESC" {
					query += fmt.Sprintf(" ORDER BY %s %s", sortBy, order)
				}
			}
		} else {
			query += " ORDER BY year DESC"
		}

		rows, err := conn.Query(c.Request.Context(), query, args...)

		if pkg.HandleErr(c, err) {
			return
		}

		defer rows.Close()

		var certs []pkg.Certification
		for rows.Next() {
			cert, err := pkg.ScanCertification(rows)
			if pkg.HandleErr(c, err) {
				return
			}
			certs = append(certs, cert)
		}

		if err = rows.Err(); pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusOK, pkg.CertificationsResponse{Certifications: certs, Count: len(certs)})
	}
}

// GetCertification godoc
// @Summary Obter certificação
// @Description Retorna uma certificação pelo ID
// @Tags Certificações
// @Accept json
// @Produce json
// @Param id path int true "ID da certificação"
// @Success 200 {object} pkg.CertificationResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /certifications/{id} [get]
func GetCertification(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		certification, err := pkg.ScanCertification(conn.QueryRow(c.Request.Context(),
			`SELECT id, name, year, description, institution
			FROM certification WHERE id = $1`, id))

		if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusOK, pkg.CertificationResponse{Certification: certification})
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
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 409 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /certifications [post]
func CreateCertification(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var cert pkg.Certification
		if err := c.BindJSON(&cert); pkg.HandleErr(c, err) {
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

		var exists bool
		err := conn.QueryRow(c.Request.Context(),
			`SELECT EXISTS(SELECT 1 FROM certification WHERE name = $1 AND institution = $2)`,
			cert.Name, cert.Institution).Scan(&exists)

		if pkg.HandleErr(c, err) {
			return
		}

		if exists {
			c.JSON(http.StatusConflict, pkg.Conflict(c.FullPath(), "Certification already exists"))
			return
		}

		if pkg.HandleErr(c, conn.QueryRow(c.Request.Context(),
			`INSERT INTO certification (name, year, description, institution) VALUES ($1, $2, $3, $4) RETURNING id`,
			cert.Name, cert.Year, cert.Description, cert.Institution).Scan(&cert.Id)) {
			return
		}
		c.JSON(http.StatusCreated, pkg.CertificationResponse{Certification: cert})
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
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /certifications/{id} [put]
func UpdateCertification(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
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

		cert, err = pkg.ScanCertification(conn.QueryRow(c.Request.Context(),
			`UPDATE certification SET name = $1, year = $2, description = $3, institution = $4 WHERE id = $5 RETURNING id, name, year, description, institution`,
			cert.Name, cert.Year, cert.Description, cert.Institution, id))

		if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusOK, pkg.CertificationResponse{Certification: cert})
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
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /certifications/{id} [patch]
func UpdateCertificationPartial(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
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

		args = append(args, id)

		query := fmt.Sprintf("UPDATE certification SET %s WHERE id = $%d RETURNING id, name, year, description, institution",
			strings.Join(set, ", "), i)

		updatedCert, err := pkg.ScanCertification(conn.QueryRow(c.Request.Context(), query, args...))

		if pkg.HandleErr(c, err) {
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
// @Success 204 "Deleção com sucesso"
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /certifications/{id} [delete]
func DeleteCertification(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		result, err := conn.Exec(c.Request.Context(),
			`DELETE FROM certification WHERE id = $1`, id)

		if pkg.HandleErr(c, err) {
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Certification not found"))
			return
		}

		c.Status(http.StatusNoContent)
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
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /users/{id}/analyst/certifications [get]
func GetAnalystCertifications(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		rows, err := conn.Query(c.Request.Context(),
			`SELECT c.id, c.name, c.year, c.description, c.institution
			 FROM certification c
			 JOIN analyst_certification ac ON c.id = ac.certification_id
			 WHERE ac.analyst_id = $1
			 ORDER BY c.year DESC`, id)

		if pkg.HandleErr(c, err) {
			return
		}
		defer rows.Close()

		var certs []pkg.Certification
		for rows.Next() {
			cert, err := pkg.ScanCertification(rows)
			if pkg.HandleErr(c, err) {
				return
			}
			certs = append(certs, cert)
		}

		if err = rows.Err(); pkg.HandleErr(c, err) {
			return
		}
		c.JSON(http.StatusOK, pkg.CertificationsResponse{Certifications: certs, Count: len(certs)})
	}
}

// AssociateAnalystCertification godoc
// @Summary Associar certificação existente ao analista
// @Description Associa uma certificação já existente a um analista pelo ID
// @Tags Certificações
// @Accept json
// @Produce json
// @Param id path int true "ID do analista"
// @Param cert_id path int true "ID da certificação"
// @Success 201 {object} pkg.CertificationResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 409 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security BearerAuth
// @Router /users/{id}/analyst/certifications/{cert_id} [post]
func AssociateAnalystCertification(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		analystID, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		certID, err := pkg.ParsePathParam(c, "cert_id")
		if err != nil {
			return
		}

		cert, err := pkg.ScanCertification(conn.QueryRow(c.Request.Context(),
			`SELECT id, name, year, description, institution FROM certification WHERE id = $1`, certID))

		if pkg.HandleErr(c, err) {
			return
		}

		_, err = conn.Exec(c.Request.Context(),
			`INSERT INTO analyst_certification (analyst_id, certification_id) VALUES ($1, $2)`,
			analystID, certID)

		if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusCreated, pkg.CertificationResponse{Certification: cert})
	}
}

// CreateAnalystCertification godoc
// @Summary Criar certificação e associar ao analista
// @Description Cria uma certificação (se não existir) e a associa a um analista
// @Tags Certificações
// @Accept json
// @Produce json
// @Param id path int true "ID do analista"
// @Param certification body pkg.Certification true "Objeto certificação"
// @Success 201 {object} pkg.CertificationResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 409 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security BearerAuth
// @Router /users/{id}/analyst/certifications [post]
func CreateAnalystCertification(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		var cert pkg.Certification
		if err := c.BindJSON(&cert); err != nil {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), err.Error()))
			return
		}

		if cert.Year < 1900 || cert.Year > 2030 {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Received invalid year. Must be between 1900 and 2030"))
			return
		}

		var analystExists bool
		err = conn.QueryRow(c.Request.Context(),
			`SELECT EXISTS(SELECT 1 FROM analyst WHERE id = $1)`, id).Scan(&analystExists)

		if pkg.HandleErr(c, err) {
			return
		}

		if !analystExists {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Analyst not found"))
			return
		}

		err = conn.QueryRow(c.Request.Context(),
			`SELECT id FROM certification WHERE name = $1 AND institution = $2`,
			cert.Name, cert.Institution).Scan(&cert.Id)

		if err == pgx.ErrNoRows {
			cert, err = pkg.ScanCertification(conn.QueryRow(c.Request.Context(),
				`INSERT INTO certification (name, year, description, institution)
				 VALUES ($1, $2, $3, $4)
				 RETURNING id, name, year, description, institution`,
				cert.Name, cert.Year, cert.Description, cert.Institution))
			if pkg.HandleErr(c, err) {
				return
			}
		} else if pkg.HandleErr(c, err) {
			return
		}

		_, err = conn.Exec(c.Request.Context(),
			`INSERT INTO analyst_certification (analyst_id, certification_id) VALUES ($1, $2)`,
			id, cert.Id)

		if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusCreated, pkg.CertificationResponse{Certification: cert})
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
// @Success 204 "Deleção com sucesso"
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /users/{id}/analyst/certifications [delete]
func DeleteAnalystCertification(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		certID, err := pkg.ParsePathQuery(c, "certification_id")
		if err != nil {
			return
		}

		result, err := conn.Exec(c.Request.Context(),
			`DELETE FROM analyst_certification WHERE analyst_id = $1 AND certification_id = $2`,
			userID, certID)

		if pkg.HandleErr(c, err) {
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Analyst certification not found"))
			return
		}

		c.Status(http.StatusNoContent)
	}
}
