package handlers

import (
	"main/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/n-r-w/squirrel"
)

const certificationSelect = `id, name, year, description, institution`

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
// @Param sort_by query string false "Campo para ordenar: name,year,institution" enums(name,year,institution)
// @Param order query string false "Direção: ASC ou DESC" enums(ASC,DESC)
// @Param page query int false "Página"
// @Param page_size query int false "Tamanho da página"
// @Success 200 {object} pkg.CertificationsResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /certifications [get]
func GetCertifications(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filters pkg.CertificationFilter
		err := c.ShouldBindQuery(&filters)
		if pkg.HandleErr(c, err) {
			return
		}
		filters.Normalize()
		builder := squirrel.Select(certificationSelect).From("certification").PlaceholderFormat(squirrel.Dollar)

		if filters.Name != "" {
			builder = builder.Where(squirrel.ILike{"name": pkg.PutPercent(filters.Name)})
		}
		if filters.Year != nil {
			builder = builder.Where(squirrel.Eq{"year": *filters.Year})
		}
		if filters.FromYear != nil {
			builder = builder.Where(squirrel.GtOrEq{"year": *filters.FromYear})
		}
		if filters.ToYear != nil {
			builder = builder.Where(squirrel.LtOrEq{"year": *filters.ToYear})
		}
		if filters.Institution != "" {
			builder = builder.Where(squirrel.ILike{"institution": pkg.PutPercent(filters.Institution)}) // ← faltava
		}

		orderClause := filters.SortOptions.ValidateSort(pkg.CertificationSortFields)

		if orderClause != "" {
			builder = builder.OrderBy(orderClause)
		} else {
			builder = builder.OrderBy("year DESC")
		}

		builder = builder.Limit(uint64(filters.PageSize)).Offset(uint64(filters.Offset()))

		query, args, err := builder.ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		rows, err := conn.Query(c.Request.Context(), query, args...)
		if pkg.HandleErr(c, err) {
			return
		}

		defer rows.Close()

		certs, err := pkg.ScanRows(c, rows, pkg.ScanCertification)
		if pkg.HandleErr(c, err) {
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

		query, args, err := squirrel.Select(certificationSelect).
			From("certification").Where(squirrel.Eq{"id": id}).
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		certification, err := pkg.ScanCertification(conn.QueryRow(c.Request.Context(), query, args...))
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
		} else if exists {
			c.JSON(http.StatusConflict, pkg.Conflict(c.FullPath(), "Certification already exists"))
			return
		}

		insertQuery, insertArgs, err := squirrel.Insert("certification").
			Columns("name", "year", "description", "institution").
			Values(cert.Name, cert.Year, cert.Description, cert.Institution).
			Suffix("RETURNING id").
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if pkg.HandleErr(c, err) {
			return
		} else if pkg.HandleErr(c, conn.QueryRow(c.Request.Context(), insertQuery, insertArgs...).Scan(&cert.Id)) {
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
		if err := c.BindJSON(&cert); pkg.HandleErr(c, err) {
			return
		}

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

		query, args, err := squirrel.Update("certification").
			SetMap(map[string]any{
				"name":        cert.Name,
				"year":        cert.Year,
				"description": cert.Description,
				"institution": cert.Institution,
			}).Where(squirrel.Eq{"id": id}).
			Suffix("RETURNING " + certificationSelect).
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		cert, err = pkg.ScanCertification(conn.QueryRow(c.Request.Context(), query, args...))
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

		var req pkg.CertificationUpdateRequest
		if err := c.BindJSON(&req); pkg.HandleErr(c, err) {
			return
		}

		builder := squirrel.Update("certification").PlaceholderFormat(squirrel.Dollar)
		hasFields := false

		if req.Name != nil {
			builder = builder.Set("name", *req.Name)
			hasFields = true
		}
		if req.Year != nil {
			if *req.Year < 1900 || *req.Year > 2030 {
				c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Received invalid year. Must be between 1900 and 2030"))
				return
			}
			builder = builder.Set("year", *req.Year)
			hasFields = true
		}
		if req.Description != nil {
			builder = builder.Set("description", *req.Description)
			hasFields = true
		}
		if req.Institution != nil {
			builder = builder.Set("institution", *req.Institution)
			hasFields = true
		}
		if !hasFields {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "No valid fields were given"))
			return
		}

		query, args, err := builder.Where(squirrel.Eq{"id": id}).
			Suffix("RETURNING " + certificationSelect).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		cert, err := pkg.ScanCertification(conn.QueryRow(c.Request.Context(), query, args...))
		if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusOK, pkg.CertificationResponse{Certification: cert})
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

		query, args, err := squirrel.Delete("certification").
			Where(squirrel.Eq{"id": id}).
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		result, err := conn.Exec(c.Request.Context(), query, args...)
		if pkg.HandleErr(c, err) {
			return
		} else if result.RowsAffected() == 0 {
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
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /users/{id}/analyst/certifications [get]
func GetAnalystCertifications(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		query, args, err := squirrel.Select("c.id, c.name, c.year, c.description, c.institution").
			From("certification c").
			Join("analyst_certification ac ON c.id = ac.certification_id").
			Where(squirrel.Eq{"ac.analyst_id": id}).OrderBy("c.year DESC").
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		rows, err := conn.Query(c.Request.Context(), query, args...)
		if pkg.HandleErr(c, err) {
			return
		}
		defer rows.Close()

		certs, err := pkg.ScanRows(c, rows, pkg.ScanCertification)
		if pkg.HandleErr(c, err) {
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

		query, args, err := squirrel.Select(certificationSelect).
			From("certification").Where(squirrel.Eq{"id": certID}).
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		cert, err := pkg.ScanCertification(conn.QueryRow(c.Request.Context(), query, args...))
		if pkg.HandleErr(c, err) {
			return
		}

		insertQuery, insertArgs, err := squirrel.Insert("analyst_certification").
			Columns("analyst_id", "certification_id").Values(analystID, certID).
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		} else if _, err = conn.Exec(c.Request.Context(), insertQuery, insertArgs...); pkg.HandleErr(c, err) {
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
		if err := c.BindJSON(&cert); pkg.HandleErr(c, err) {
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
		} else if !analystExists {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Analyst not found"))
			return
		}

		err = conn.QueryRow(c.Request.Context(),
			`SELECT id FROM certification WHERE name = $1 AND institution = $2`,
			cert.Name, cert.Institution).Scan(&cert.Id)

		if err == pgx.ErrNoRows {
			insertQuery, insertArgs, err := squirrel.Insert("certification").
				Columns("name", "year", "description", "institution").
				Values(cert.Name, cert.Year, cert.Description, cert.Institution).
				Suffix("RETURNING " + certificationSelect).
				PlaceholderFormat(squirrel.Dollar).
				ToSql()
			if pkg.HandleErr(c, err) {
				return
			}
			cert, err = pkg.ScanCertification(conn.QueryRow(c.Request.Context(), insertQuery, insertArgs...))
			if pkg.HandleErr(c, err) {
				return
			}
		} else if pkg.HandleErr(c, err) {
			return
		}

		assocQuery, assocArgs, err := squirrel.Insert("analyst_certification").
			Columns("analyst_id", "certification_id").
			Values(id, cert.Id).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if pkg.HandleErr(c, err) {
			return
		} else if _, err = conn.Exec(c.Request.Context(), assocQuery, assocArgs...); pkg.HandleErr(c, err) {
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

		query, args, err := squirrel.Delete("analyst_certification").
			Where(squirrel.Eq{"analyst_id": userID, "certification_id": certID}).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		result, err := conn.Exec(c.Request.Context(), query, args...)
		if pkg.HandleErr(c, err) {
			return
		} else if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Analyst certification not found"))
			return
		}

		c.Status(http.StatusNoContent)
	}
}
