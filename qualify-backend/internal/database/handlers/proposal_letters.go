package handlers

import (
	"main/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/n-r-w/squirrel"
)

const proposalSelect = "id, client_id, analyst_id, proposed_hourly_rate, title, content, time_created"

func applyProposalFilters(builder squirrel.SelectBuilder, filters pkg.ProposalFilter) squirrel.SelectBuilder {
	if filters.ClientId != nil {
		builder = builder.Where(squirrel.Eq{"client_id": *filters.ClientId})
	}
	if filters.AnalystId != nil {
		builder = builder.Where(squirrel.Eq{"analyst_id": *filters.AnalystId})
	}
	if filters.Title != "" {
		builder = builder.Where(squirrel.ILike{"title": pkg.PutPercent(filters.Title)})
	}
	if filters.Content != "" {
		builder = builder.Where(squirrel.ILike{"content": pkg.PutPercent(filters.Content)})
	}
	if filters.MinProposedHourlyRate != nil {
		builder = builder.Where(squirrel.GtOrEq{"proposed_hourly_rate": *filters.MinProposedHourlyRate})
	}
	if filters.MaxProposedHourlyRate != nil {
		builder = builder.Where(squirrel.LtOrEq{"proposed_hourly_rate": *filters.MaxProposedHourlyRate})
	}
	if order := filters.ValidateSort(pkg.ProposalSortFields); order != "" {
		builder = builder.OrderBy(order)
	} else {
		builder = builder.OrderBy("time_created DESC")
	}
	return builder.Limit(uint64(filters.PageSize)).Offset(uint64(filters.Offset()))
}

func proposalBuilder(filters pkg.ProposalFilter) squirrel.SelectBuilder {
	return applyProposalFilters(
		squirrel.Select(proposalSelect).From("proposal_letter").PlaceholderFormat(squirrel.Dollar),
		filters,
	)
}

func scanProposals(c *gin.Context, conn *pgxpool.Pool, builder squirrel.SelectBuilder) {
	query, args, err := builder.ToSql()
	if pkg.HandleErr(c, err) {
		return
	}

	rows, err := conn.Query(c.Request.Context(), query, args...)
	if pkg.HandleErr(c, err) {
		return
	}
	defer rows.Close()

	proposals, err := pkg.ScanRows(c, rows, pkg.ScanProposalLetter)
	if pkg.HandleErr(c, err) {
		return
	}

	c.JSON(http.StatusOK, pkg.ProposalLettersResponse{Proposal_letters: proposals, Count: len(proposals)})
}

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
// @Param page query int false "Página"
// @Param page_size query int false "Tamanho da página"
// @Success 200 {object} pkg.ProposalLettersResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /proposals [get]
func GetProposalLetters(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filters pkg.ProposalFilter
		err := c.ShouldBindQuery(&filters)
		if pkg.HandleErr(c, err) {
			return
		}
		filters.Normalize()

		scanProposals(c, conn, proposalBuilder(filters))
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
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		query, args, err := squirrel.Select(proposalSelect).
			From("proposal_letter").Where(squirrel.Eq{"id": id}).
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		proposal, err := pkg.ScanProposalLetter(conn.QueryRow(c.Request.Context(), query, args...))
		if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusOK, pkg.ProposalLetterResponse{Proposal_letter: proposal})
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
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /proposals [post]
func CreateProposalLetter(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var proposal pkg.ProposalLetter
		if err := c.BindJSON(&proposal); pkg.HandleErr(c, err) {
			return
		}

		// Checando se o cliente já existe
		var analystExists bool
		err := conn.QueryRow(c.Request.Context(),
			`SELECT EXISTS(SELECT 1 FROM analyst WHERE id = $1)`, proposal.Analyst_id).Scan(&analystExists)
		if pkg.HandleErr(c, err) {
			return
		} else if !analystExists {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Analyst does not exists"))
			return
		}

		// Checando se o cliente já existe
		var clientExists bool
		err = conn.QueryRow(c.Request.Context(),
			`SELECT EXISTS(SELECT 1 FROM client WHERE id = $1)`, proposal.Client_id).Scan(&clientExists)
		if pkg.HandleErr(c, err) {
			return
		} else if !clientExists {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Client does not exists"))
			return
		}

		query, args, err := squirrel.Insert("proposal_letter").
			Columns("title", "content", "client_id", "analyst_id", "proposed_hourly_rate").
			Values(proposal.Title, proposal.Content, proposal.Client_id, proposal.Analyst_id, proposal.Proposed_hourly_rate).
			Suffix("RETURNING " + proposalSelect).
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		proposal, err = pkg.ScanProposalLetter(conn.QueryRow(c.Request.Context(), query, args...))
		if pkg.HandleErr(c, err) {
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
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		var proposal pkg.ProposalLetter
		if err := c.BindJSON(&proposal); pkg.HandleErr(c, err) {
			return
		} else if proposal.Title == "" {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Received empty title"))
			return
		}

		query, args, err := squirrel.Update("proposal_letter").
			SetMap(map[string]any{
				"title":                proposal.Title,
				"content":              proposal.Content,
				"client_id":            proposal.Client_id,
				"analyst_id":           proposal.Analyst_id,
				"proposed_hourly_rate": proposal.Proposed_hourly_rate}).
			Where(squirrel.Eq{"id": id}).
			Suffix("RETURNING " + proposalSelect).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		proposal, err = pkg.ScanProposalLetter(conn.QueryRow(c.Request.Context(), query, args...))
		if pkg.HandleErr(c, err) {
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
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		var req pkg.ProposalLetterUpdateRequest
		if err := c.BindJSON(&req); pkg.HandleErr(c, err) {
			return
		}

		builder := squirrel.Update("proposal_letter").PlaceholderFormat(squirrel.Dollar)
		hasFields := false

		if req.Title != nil {
			builder = builder.Set("title", *req.Title)
			hasFields = true
		}
		if req.Content != nil {
			builder = builder.Set("content", *req.Content)
			hasFields = true
		}
		if req.Proposed_hourly_rate != nil {
			builder = builder.Set("proposed_hourly_rate", *req.Proposed_hourly_rate)
			hasFields = true
		}

		if !hasFields {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "No valid fields were given"))
			return
		}

		query, args, err := builder.
			Where(squirrel.Eq{"id": id}).
			Suffix("RETURNING " + proposalSelect).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		proposal, err := pkg.ScanProposalLetter(conn.QueryRow(c.Request.Context(), query, args...))
		if pkg.HandleErr(c, err) {
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
// @Success 204 "Deleção com sucesso"
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security     BearerAuth
// @Router /proposals/{id} [delete]
func DeleteProposalLetter(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		query, args, err := squirrel.Delete("proposal_letter").
			Where(squirrel.Eq{"id": id}).
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		result, err := conn.Exec(c.Request.Context(), query, args...)
		if pkg.HandleErr(c, err) {
			return
		} else if result.RowsAffected() == 0 {
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
// @Param page query int false "Página"
// @Param page_size query int false "Tamanho da página"
// @Success 200 {object} pkg.ProposalLettersResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /users/{id}/client/proposals [get]
func GetClientProposalLetters(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		// Checando se o cliente já existe
		var clientExists bool
		err = conn.QueryRow(c.Request.Context(),
			`SELECT EXISTS(SELECT 1 FROM client WHERE id = $1)`, id).Scan(&clientExists)
		if pkg.HandleErr(c, err) {
			return
		} else if !clientExists {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Client does not exists"))
			return
		}

		var filters pkg.ProposalFilter
		err = c.ShouldBindQuery(&filters)
		if pkg.HandleErr(c, err) {
			return
		}
		filters.Normalize()

		builder := squirrel.Select("p." + proposalSelect).
			From("proposal_letter p").
			Join("client cl ON p.client_id = cl.id").
			Where(squirrel.Eq{"cl.id": id}).
			PlaceholderFormat(squirrel.Dollar)
		builder = applyProposalFilters(builder, filters)

		scanProposals(c, conn, builder)
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
// @Param page query int false "Página"
// @Param page_size query int false "Tamanho da página"
// @Success 200 {object} pkg.ProposalLettersResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /users/{id}/analyst/proposals [get]
func GetAnalystProposalLetters(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		// Checando se o analista já existe
		var analystExists bool
		err = conn.QueryRow(c.Request.Context(),
			`SELECT EXISTS(SELECT 1 FROM analyst WHERE id = $1)`, id).Scan(&analystExists)
		if pkg.HandleErr(c, err) {
			return
		} else if !analystExists {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Analyst does not exists"))
			return
		}

		var filters pkg.ProposalFilter
		err = c.ShouldBindQuery(&filters)
		if pkg.HandleErr(c, err) {
			return
		}
		filters.Normalize()

		builder := squirrel.Select("p." + proposalSelect).
			From("proposal_letter p").
			Join("client cl ON p.client_id = cl.id").
			Where(squirrel.Eq{"cl.id": id}).
			PlaceholderFormat(squirrel.Dollar)
		builder = applyProposalFilters(builder, filters)

		scanProposals(c, conn, builder)
	}
}
