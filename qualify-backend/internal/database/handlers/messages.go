package handlers

import (
	"main/pkg"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/n-r-w/squirrel"
)

// GetConversationMessages godoc
// @Summary Obter mensagens de uma conversa
// @Description Retorna todas as mensagens de uma conversa específica
// @Tags Mensagens
// @Accept json
// @Produce json
// @Param conv_id path int true "ID da conversa"
// @Success 200 {object} pkg.MessageResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Security BearerAuth
// @Router /conversations/{conv_id}/messages [get]
func GetConversationMessages(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		convID, err := pkg.ParsePathParam(c, "conv_id")
		if err != nil {
			return
		}

		query, args, err := squirrel.Select("m.id", "m.conversation_id", "m.sender_id", "m.content", "m.created_at", "m.read_at").
			From("message m").
			Where(squirrel.Eq{"m.conversation_id": convID}).
			OrderBy("m.created_at ASC").
			PlaceholderFormat(squirrel.Dollar).ToSql()

		if pkg.HandleErr(c, err) {
			return
		}

		rows, err := conn.Query(c.Request.Context(), query, args...)
		if pkg.HandleErr(c, err) {
			return
		}
		defer rows.Close()

		messages, err := pkg.ScanRows(c, rows, pkg.ScanMessage)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, pkg.MessageResponse{Messages: messages, Count: len(messages)})
	}
}

// CreateMessage godoc
// @Summary Criar mensagem
// @Description Adiciona uma nova mensagem a uma conversa
// @Tags Mensagens
// @Accept json
// @Produce json
// @Param conv_id path int true "ID da conversa"
// @Param message body pkg.MessageCreateRequest true "Conteúdo da mensagem"
// @Success 201 {object} pkg.Message
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Security BearerAuth
// @Router /conversations/{conv_id}/messages [post]
func CreateMessage(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		convID, err := pkg.ParsePathParam(c, "conv_id")
		if err != nil {
			return
		}

		var req pkg.MessageCreateRequest
		if err := c.BindJSON(&req); pkg.HandleErr(c, err) {
			return
		}

		// Validação básica: conversão de string para int se necessário (se o ID for enviado como string no JSON)
		senderID := req.Sender_id
		if req.Content == "" || convID == 0 || senderID == 0 {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "Invalid data"))
			return
		}

		ins := squirrel.Insert("message").
			Columns("conversation_id", "sender_id", "content").
			Values(convID, senderID, req.Content).
			Suffix("RETURNING id, conversation_id, sender_id, content, created_at, read_at")

		query, args, err := ins.PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		msg, err := pkg.ScanMessage(conn.QueryRow(c.Request.Context(), query, args...))
		if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusCreated, msg)
	}
}

// UpdateMessage godoc
// @Summary Atualizar mensagem
// @Description Atualiza campos de uma mensagem pelo ID (via rota /messages/{msg_id})
// @Tags Mensagens
// @Accept json
// @Produce json
// @Param conv_id path int true "ID da conversa"
// @Param msg_id path int true "ID da mensagem"
// @Param message body pkg.MessageUpdateRequest true "Campos a atualizar"
// @Success 200 {object} pkg.Message
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Security BearerAuth
// @Router /conversations/{conv_id}/messages/{msg_id} [put]
func UpdateMessage(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		convID, err := pkg.ParsePathParam(c, "conv_id")
		if err != nil {
			return
		}

		msgID, err := pkg.ParsePathParam(c, "msg_id")
		if err != nil {
			return
		}

		var req pkg.MessageUpdateRequest
		if err := c.BindJSON(&req); pkg.HandleErr(c, err) {
			return
		}

		upd := squirrel.Update("message")
		anyField := false

		if req.Content != "" {
			upd = upd.Set("content", req.Content)
			anyField = true
		}

		if !anyField {
			c.JSON(http.StatusBadRequest, pkg.BadRequest(c.FullPath(), "No fields provided for update"))
			return
		}

		query, args, err := upd.
			Where(squirrel.Eq{"id": msgID}).
			Where(squirrel.Eq{"conversation_id": convID}).
			Suffix("RETURNING id, conversation_id, sender_id, content, created_at, read_at").
			PlaceholderFormat(squirrel.Dollar).ToSql()

		if pkg.HandleErr(c, err) {
			return
		}

		msg, err := pkg.ScanMessage(conn.QueryRow(c.Request.Context(), query, args...))
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Message not found"))
			return
		} else if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusOK, msg)
	}
}

// ReadMessage godoc
// @Summary Lê mensagem
// @Description Marca uma mensagem como lida e atualiza o campo `read_at` (só poder ser feito uma vez)
// @Tags Mensagens
// @Accept json
// @Produce json
// @Param conv_id path int true "ID da conversa"
// @Param msg_id path int true "ID da mensagem"
// @Success 200 "Leitura com sucesso"
// @Failure 400 {object} pkg.ErrorResponse
// @Failure 404 {object} pkg.ErrorResponse
// @Security BearerAuth
// @Router /conversations/{conv_id}/messages/{msg_id}/read [put]
func ReadMessage(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		convID, err := pkg.ParsePathParam(c, "conv_id")
		if err != nil {
			return
		}

		msgID, err := pkg.ParsePathParam(c, "msg_id")
		if err != nil {
			return
		}

		query, args, err := squirrel.Update("message").
			Set("read_at", time.Now().UTC()).
			Where(squirrel.Eq{"id": msgID}).
			Where(squirrel.Eq{"conversation_id": convID}).
			Where(squirrel.Eq{"read_at": nil}). // só pode ser feito uma vez
			Suffix("RETURNING id, conversation_id, sender_id, content, created_at, read_at").
			PlaceholderFormat(squirrel.Dollar).ToSql()
		if pkg.HandleErr(c, err) {
			return
		}

		msg, err := pkg.ScanMessage(conn.QueryRow(c.Request.Context(), query, args...))
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Message not found or already read"))
			return
		} else if pkg.HandleErr(c, err) {
			return
		}

		c.JSON(http.StatusOK, msg)
	}
}

// DeleteMessage godoc
// @Summary Excluir mensagem
// @Description Remove uma mensagem pelo ID
// @Tags Mensagens
// @Accept json
// @Produce json
// @Param conv_id path int true "ID da conversa"
// @Param msg_id path int true "ID da mensagem"
// @Success 204 "Deleção com sucesso"
// @Failure 404 {object} pkg.ErrorResponse
// @Security BearerAuth
// @Router /conversations/{conv_id}/messages/{msg_id} [delete]
func DeleteMessage(conn *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		convID, err := pkg.ParsePathParam(c, "conv_id")
		if err != nil {
			return
		}
		msgID, err := pkg.ParsePathParam(c, "msg_id")
		if err != nil {
			return
		}

		query, args, err := squirrel.Delete("message").
			Where(squirrel.Eq{"id": msgID}).
			Where(squirrel.Eq{"conversation_id": convID}).
			PlaceholderFormat(squirrel.Dollar).ToSql()

		if pkg.HandleErr(c, err) {
			return
		}

		result, err := conn.Exec(c.Request.Context(), query, args...)
		if pkg.HandleErr(c, err) {
			return
		} else if result.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, pkg.NotFound(c.FullPath(), "Message not found"))
			return
		}

		c.Status(http.StatusNoContent)
	}
}
