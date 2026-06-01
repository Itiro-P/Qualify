package pkg

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func newError(path, httpError, code, message string) ErrorResponse {
	return ErrorResponse{
		Error:     httpError,
		Message:   message,
		Code:      code,
		Timestamp: time.Now(),
		Path:      path,
	}
}

func NotFound(path, message string) ErrorResponse {
	return newError(path, "Not Found", "NOT_FOUND", message)
}

func UnprocessableEntity(path, message string) ErrorResponse {
	return newError(path, "Unprocessable Entity", "UNPROCESSABLE_ENTITY", message)
}

func Unauthorized(path, message string) ErrorResponse {
	return newError(path, "Unauthorized", "UNAUTHORIZED", message)
}

func Forbidden(path, message string) ErrorResponse {
	return newError(path, "Forbidden", "FORBIDDEN", message)
}

func Conflict(path, message string) ErrorResponse {
	return newError(path, "Conflict", "CONFLICT", message)
}

func BadRequest(path, message string) ErrorResponse {
	return newError(path, "Bad Request", "BAD_REQUEST", message)
}

func Internal(path, message string) ErrorResponse {
	return newError(path, "Internal Server Error", "INTERNAL_ERROR", message)
}

func ValidationFailed(path, message string, validationErrors map[string]string) ErrorResponse {
	r := newError(path, "Validation Error", "VALIDATION_ERROR", message)
	r.ValidationErrors = validationErrors
	return r
}

// Função auxiliar para tratar erros
// Retorna `true` caso seja necessário encerrar a execução
func HandleErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, NotFound(c.FullPath(), err.Error()))
		return true
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		// Integridade
		case "23505": // unique_violation
			c.JSON(http.StatusConflict, Conflict(c.FullPath(), pgErr.Detail))
		case "23503": // foreign_key_violation
			c.JSON(http.StatusUnprocessableEntity, UnprocessableEntity(c.FullPath(), pgErr.Detail))
		case "23502": // not_null_violation
			c.JSON(http.StatusBadRequest, BadRequest(c.FullPath(), pgErr.Detail))
		case "23514": // check_violation
			c.JSON(http.StatusBadRequest, BadRequest(c.FullPath(), pgErr.Detail))

		// Conexão / recursos
		case "53300": // too_many_connections
			c.JSON(http.StatusServiceUnavailable, Internal(c.FullPath(), "too many connections"))
		case "57014": // query_canceled (ex: context timeout)
			c.JSON(http.StatusGatewayTimeout, Internal(c.FullPath(), "query timed out"))

		// Autenticação / permissão
		case "28000", "28P01": // invalid_authorization / wrong password
			c.JSON(http.StatusUnauthorized, Unauthorized(c.FullPath(), pgErr.Message))
		case "42501": // insufficient_privilege
			c.JSON(http.StatusForbidden, Forbidden(c.FullPath(), pgErr.Message))

		default:
			c.JSON(http.StatusInternalServerError, Internal(c.FullPath(), pgErr.Message))
		}
		return true
	}
	c.JSON(http.StatusInternalServerError, Internal(c.FullPath(), err.Error()))
	return true
}
