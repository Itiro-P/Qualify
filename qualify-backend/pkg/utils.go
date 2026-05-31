package pkg

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

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

// Função auxiliar que realiza o parse do parâmetro ID.
// Retorna -1 se ocorrer erros.
func ParseIdParam(c *gin.Context) (int, error) {
	_id := c.Param("id")
	id, errID := strconv.Atoi(_id)
	if errID != nil {
		c.JSON(http.StatusBadRequest, BadRequest(c.FullPath(), errID.Error()))
		return -1, errID
	}
	return id, nil
}

// Função auxiliar que realiza o parse da query.
// Retorna -1 se ocorrer erros.
func ParsePathQuery(c *gin.Context, query string) (int, error) {
	_id := c.Query(query)
	id, errID := strconv.Atoi(_id)
	if errID != nil {
		c.JSON(http.StatusBadRequest, BadRequest(c.FullPath(), errID.Error()))
		return -1, errID
	}
	return id, nil
}

// Função auxiliar que realiza o parse do path.
// Retorna -1 se ocorrer erros.
func ParsePathParam(c *gin.Context, param string) (int, error) {
	id, err := strconv.Atoi(c.Param(param))
	if err != nil {
		c.JSON(http.StatusBadRequest, BadRequest(c.FullPath(), err.Error()))
		return -1, err
	}
	return id, nil
}

// Função auxiliar que escaneia usuários
func ScanUser(row pgx.Row) (User, error) {
	var u User
	err := row.Scan(
		&u.Id,
		&u.Name,
		&u.Email,
		&u.Phone,
		&u.Time_created,
		&u.Country_code,
		&u.Country_name,
		&u.Country_state,
		&u.City,
		&u.Timezone,
	)
	return u, err
}

// Função auxiliar que escaneia analistas
func ScanAnalyst(row pgx.Row) (Analyst, error) {
	var a Analyst
	err := row.Scan(
		&a.Id,
		&a.Name,
		&a.Email,
		&a.Phone,
		&a.Time_created,
		&a.Country_code,
		&a.Country_name,
		&a.Country_state,
		&a.City,
		&a.Timezone,
		&a.Hourly_rate,
		&a.Total_reviews,
		&a.Mean_rating,
	)
	return a, err
}

// Função auxiliar que escaneia clientes
func ScanClient(row pgx.Row) (Client, error) {
	var c Client
	err := row.Scan(
		&c.Id,
		&c.Name,
		&c.Email,
		&c.Phone,
		&c.Time_created,
		&c.Country_code,
		&c.Country_name,
		&c.Country_state,
		&c.City,
		&c.Timezone,
		&c.Proposed_budget,
	)
	return c, err
}

// Função auxiliar que escaneia certificações
func ScanCertification(row pgx.Row) (Certification, error) {
	var c Certification
	err := row.Scan(
		&c.Id,
		&c.Name,
		&c.Year,
		&c.Description,
		&c.Institution)
	return c, err
}

// Função auxiliar que escaneia habilidades
func ScanSkill(row pgx.Row) (Skill, error) {
	var s Skill
	err := row.Scan(
		&s.Id,
		&s.Name)
	return s, err
}

// Função auxiliar que escaneia avaliações
func ScanReview(row pgx.Row) (Review, error) {
	var r Review
	err := row.Scan(
		&r.Id,
		&r.Analyst_id,
		&r.Client_id,
		&r.Service_id,
		&r.Rating,
		&r.Comment,
		&r.Time_created)
	return r, err
}

// Função auxiliar que escaneia cartas-proposta
func ScanProposalLetter(row pgx.Row) (ProposalLetter, error) {
	var p ProposalLetter
	err := row.Scan(
		&p.Id,
		&p.Client_id,
		&p.Analyst_id,
		&p.Proposed_hourly_rate,
		&p.Title,
		&p.Content,
		&p.Time_created)
	return p, err
}

// Função auxiliar que escaneia serviços
func ScanService(row pgx.Row) (Service, error) {
	var s Service
	err := row.Scan(
		&s.Id,
		&s.Title,
		&s.Content,
		&s.Proposal_letter_id,
		&s.Hourly_rate,
		&s.Status,
		&s.Time_created)
	return s, err
}

// Função auxiliar que escaneia perfis de usuários
func ScanProfile(row pgx.Row) (UserProfile, error) {
	var p UserProfile
	err := row.Scan(
		&p.User_id,
		&p.Biography,
		&p.Picture)
	return p, err
}

// Função auxiliar que escaneia perfis de analistas
func ScanAnalystProfile(row pgx.Row) (AnalystProfile, error) {
	var p AnalystProfile
	err := row.Scan(
		&p.User_id,
		&p.Biography,
		&p.Picture)
	return p, err
}

// Função auxiliar que escaneia perfis de clientes
func ScanClientProfile(row pgx.Row) (ClientProfile, error) {
	var p ClientProfile
	err := row.Scan(
		&p.User_id,
		&p.Biography,
		&p.Picture)
	return p, err
}
