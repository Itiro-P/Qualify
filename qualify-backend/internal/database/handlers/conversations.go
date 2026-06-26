package handlers

import (
	"main/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/n-r-w/squirrel"
)

const conversationSelect = `conv.id, conv.service_id, conv.proposal_id, conv.client_id, conv.analyst_id, conv.created_at`

// GetAnalystConversations godoc
// @Summary Obter conversas do analista
// @Description Retorna as conversas em que um usuário está relacionado como analista
// @Tags Conversas
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário (analista)"
// @Success 200 {object} pkg.ConversationResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /users/{id}/analyst/conversations [get]
func GetAnalystConversations(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		query, args, err := squirrel.Select(conversationSelect).
			From("conversation conv").Where(squirrel.Eq{"conv.analyst_id": id}).
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		rows, err := conn.Query(c.Request.Context(), query, args...)
		if pkg.HandleErr(c, err) {
			return
		}
		defer rows.Close()

		conversations, err := pkg.ScanRows(c, rows, pkg.ScanConversation)
		if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusOK, pkg.ConversationResponse{Conversations: conversations, Count: len(conversations)})
	}
}

// GetAnalystConversation godoc
// @Summary Obter conversa do analista pelo ID
// @Description Retorna uma conversa específica em que o usuário está relacionado como analista
// @Tags Conversas
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário (analista)"
// @Param conv_id path int true "ID da conversa"
// @Success 200 {object} pkg.ConversationResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /users/{id}/analyst/conversations/{conv_id} [get]
func GetAnalystConversation(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		convID, err := pkg.ParsePathParam(c, "conv_id")
		if err != nil {
			return
		}

		query, args, err := squirrel.Select(conversationSelect).
			From("conversation conv").
			Where(squirrel.Eq{"conv.id": convID, "conv.analyst_id": id}).
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		conversation, err := pkg.ScanConversation(conn.QueryRow(c.Request.Context(), query, args...))
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Conversation not found"))
			return
		} else if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusOK, pkg.ConversationResponse{Conversations: []pkg.Conversation{conversation}, Count: 1})
	}
}

// GetClientConversations godoc
// @Summary Obter conversas do cliente
// @Description Retorna as conversas em que um usuário está relacionado como cliente
// @Tags Conversas
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário (cliente)"
// @Success 200 {object} pkg.ConversationResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /users/{id}/client/conversations [get]
func GetClientConversations(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		query, args, err := squirrel.Select(conversationSelect).
			From("conversation conv").Where(squirrel.Eq{"conv.client_id": id}).
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		rows, err := conn.Query(c.Request.Context(), query, args...)
		if pkg.HandleErr(c, err) {
			return
		}
		defer rows.Close()

		conversations, err := pkg.ScanRows(c, rows, pkg.ScanConversation)
		if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusOK, pkg.ConversationResponse{Conversations: conversations, Count: len(conversations)})
	}
}

// GetClientConversation godoc
// @Summary Obter conversa do cliente pelo ID
// @Description Retorna uma conversa específica em que o usuário está relacionado como cliente
// @Tags Conversas
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário (cliente)"
// @Param conv_id path int true "ID da conversa"
// @Success 200 {object} pkg.ConversationResponse
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Router /users/{id}/client/conversations/{conv_id} [get]
func GetClientConversation(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := pkg.ParseIdParam(c)
		if err != nil {
			return
		}

		convID, err := pkg.ParsePathParam(c, "conv_id")
		if err != nil {
			return
		}

		query, args, err := squirrel.Select(conversationSelect).
			From("conversation conv").
			Where(squirrel.Eq{"conv.id": convID, "conv.client_id": id}).
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		conversation, err := pkg.ScanConversation(conn.QueryRow(c.Request.Context(), query, args...))
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Conversation not found"))
			return
		} else if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusOK, pkg.ConversationResponse{Conversations: []pkg.Conversation{conversation}, Count: 1})
	}
}

// CreateConversation godoc
// @Summary Criar conversa
// @Description Cria uma nova conversa entre um cliente e um analista
// @Tags Conversas
// @Accept json
// @Produce json
// @Param conversation body pkg.ConversationCreateRequest true "Objeto conversa"
// @Success 201 {object} pkg.Conversation
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 409 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security BearerAuth
// @Router /conversations [post]
func CreateConversation(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req pkg.ConversationCreateRequest
		if err := c.BindJSON(&req); pkg.HandleErr(c, err) {
			return
		}

		ctx := c.Request.Context()

		var analystExists bool
		err := conn.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM analyst WHERE id = $1)`, req.Analyst_id).Scan(&analystExists)
		if pkg.HandleErr(c, err) {
			return
		} else if !analystExists {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Analyst not found"))
			return
		}

		var clientExists bool
		err = conn.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM client WHERE id = $1)`, req.Client_id).Scan(&clientExists)
		if pkg.HandleErr(c, err) {
			return
		} else if !clientExists {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Client not found"))
			return
		}

		if req.Service_id == nil || *req.Service_id == 0 {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "service_id is required"))
			return
		}
		if req.Proposal_id == nil || *req.Proposal_id == 0 {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "proposal_id is required"))
			return
		}

		var serviceExists bool
		err = conn.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM "service" WHERE id = $1)`, *req.Service_id).Scan(&serviceExists)
		if pkg.HandleErr(c, err) {
			return
		} else if !serviceExists {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Service not found"))
			return
		}

		var proposalExists bool
		err = conn.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM proposal_letter WHERE id = $1)`, *req.Proposal_id).Scan(&proposalExists)
		if pkg.HandleErr(c, err) {
			return
		} else if !proposalExists {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Proposal not found"))
			return
		}

		ins := squirrel.Insert("conversation").Columns("analyst_id", "client_id").Values(req.Analyst_id, req.Client_id)

		if req.Service_id != nil && *req.Service_id != 0 {
			ins = squirrel.Insert("conversation").
				Columns("analyst_id", "client_id", "service_id").
				Values(req.Analyst_id, req.Client_id, req.Service_id)
		}
		if req.Proposal_id != nil && *req.Proposal_id != 0 {
			ins = squirrel.Insert("conversation").
				Columns("analyst_id", "client_id", "proposal_id").
				Values(req.Analyst_id, req.Client_id, req.Proposal_id)
		}
		if req.Service_id != nil && *req.Service_id != 0 && req.Proposal_id != nil && *req.Proposal_id != 0 {
			ins = squirrel.Insert("conversation").
				Columns("analyst_id", "client_id", "service_id", "proposal_id").
				Values(req.Analyst_id, req.Client_id, req.Service_id, req.Proposal_id)
		}

		ins = ins.Suffix("RETURNING id, service_id, proposal_id, client_id, analyst_id, created_at")

		query, args, err := ins.PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		conversation, err := pkg.ScanConversation(conn.QueryRow(ctx, query, args...))
		if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusCreated, conversation)
	}
}

// UpdateConversation godoc
// @Summary Atualizar conversa
// @Description Atualiza os campos de uma conversa pelo ID
// @Tags Conversas
// @Accept json
// @Produce json
// @Param conv_id path int true "ID da conversa"
// @Param conversation body pkg.ConversationUpdateRequest true "Objeto conversa"
// @Success 200 {object} pkg.Conversation
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security BearerAuth
// @Router /conversations/{conv_id} [put]
func UpdateConversation(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		convID, err := pkg.ParsePathParam(c, "conv_id")
		if err != nil {
			return
		}

		var req pkg.ConversationUpdateRequest
		if err := c.BindJSON(&req); pkg.HandleErr(c, err) {
			return
		}

		upd := squirrel.Update("conversation")

		if req.Analyst_id != nil {
			upd = upd.Set("analyst_id", *req.Analyst_id)
		}
		if req.Client_id != nil {
			upd = upd.Set("client_id", *req.Client_id)
		}
		if req.Service_id != nil {
			if *req.Service_id == 0 {
				c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "service_id is required"))
				return
			}

			var serviceExists bool
			err = conn.QueryRow(c.Request.Context(), `SELECT EXISTS(SELECT 1 FROM "service" WHERE id = $1)`, *req.Service_id).Scan(&serviceExists)
			if pkg.HandleErr(c, err) {
				return
			} else if !serviceExists {
				c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Service not found"))
				return
			}
			upd = upd.Set("service_id", *req.Service_id)
		}
		if req.Proposal_id != nil {
			if *req.Proposal_id == 0 {
				c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "proposal_id is required"))
				return
			}

			var proposalExists bool
			err = conn.QueryRow(c.Request.Context(), `SELECT EXISTS(SELECT 1 FROM proposal_letter WHERE id = $1)`, *req.Proposal_id).Scan(&proposalExists)
			if pkg.HandleErr(c, err) {
				return
			} else if !proposalExists {
				c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Proposal not found"))
				return
			}
			upd = upd.Set("proposal_id", *req.Proposal_id)
		}

		query, args, err := upd.
			Where(squirrel.Eq{"id": convID}).
			Suffix("RETURNING id, service_id, proposal_id, client_id, analyst_id, created_at").
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		conversation, err := pkg.ScanConversation(conn.QueryRow(c.Request.Context(), query, args...))
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Conversation not found"))
			return
		} else if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusOK, conversation)
	}
}

// UpdateConversationPartial godoc
// @Summary Atualizar conversa parcialmente
// @Description Atualiza apenas os campos enviados de uma conversa pelo ID
// @Tags Conversas
// @Accept json
// @Produce json
// @Param conv_id path int true "ID da conversa"
// @Param conversation body pkg.ConversationUpdateRequest true "Campos a atualizar"
// @Success 200 {object} pkg.Conversation
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security BearerAuth
// @Router /conversations/{conv_id} [patch]
func UpdateConversationPartial(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		convID, err := pkg.ParsePathParam(c, "conv_id")
		if err != nil {
			return
		}

		var req pkg.ConversationUpdateRequest
		if err := c.BindJSON(&req); pkg.HandleErr(c, err) {
			return
		}

		upd := squirrel.Update("conversation")
		anyField := false

		if req.Analyst_id != nil {
			upd = upd.Set("analyst_id", *req.Analyst_id)
			anyField = true
		}
		if req.Client_id != nil {
			upd = upd.Set("client_id", *req.Client_id)
			anyField = true
		}
		if req.Service_id != nil {
			if *req.Service_id == 0 {
				c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "service_id is required"))
				return
			}

			var serviceExists bool
			err = conn.QueryRow(c.Request.Context(), `SELECT EXISTS(SELECT 1 FROM "service" WHERE id = $1)`, *req.Service_id).Scan(&serviceExists)
			if pkg.HandleErr(c, err) {
				return
			} else if !serviceExists {
				c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Service not found"))
				return
			}
			upd = upd.Set("service_id", *req.Service_id)
			anyField = true
		}
		if req.Proposal_id != nil {
			if *req.Proposal_id == 0 {
				c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "proposal_id is required"))
				return
			}

			var proposalExists bool
			err = conn.QueryRow(c.Request.Context(), `SELECT EXISTS(SELECT 1 FROM proposal_letter WHERE id = $1)`, *req.Proposal_id).Scan(&proposalExists)
			if pkg.HandleErr(c, err) {
				return
			} else if !proposalExists {
				c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Proposal not found"))
				return
			}
			upd = upd.Set("proposal_id", *req.Proposal_id)
			anyField = true
		}

		if !anyField {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "No fields provided for update"))
			return
		}

		query, args, err := upd.
			Where(squirrel.Eq{"id": convID}).
			Suffix("RETURNING id, service_id, proposal_id, client_id, analyst_id, created_at").
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		conversation, err := pkg.ScanConversation(conn.QueryRow(c.Request.Context(), query, args...))
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Conversation not found"))
			return
		} else if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusOK, conversation)
	}
}

// DeleteConversation godoc
// @Summary Excluir conversa
// @Description Remove uma conversa pelo ID
// @Tags Conversas
// @Accept json
// @Produce json
// @Param conv_id path int true "ID da conversa"
// @Success 204 "Deleção com sucesso"
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Failure 500 {object} pkg.ErrorResponse
// @Security BearerAuth
// @Router /conversations/{conv_id} [delete]
func DeleteConversation(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		convID, err := pkg.ParsePathParam(c, "conv_id")
		if err != nil {
			return
		}

		query, args, err := squirrel.Delete("conversation").
			Where(squirrel.Eq{"id": convID}).
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		result, err := conn.Exec(c.Request.Context(), query, args...)
		if pkg.HandleErr(c, err) {
			return
		} else if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Conversation not found"))
			return
		}

		c.Status(http.StatusNoContent)
	}
}
