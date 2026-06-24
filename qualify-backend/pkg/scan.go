package pkg

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

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

// Função auxiliar que escaneia Conversas
func ScanConversation(row pgx.Row) (Conversation, error) {
	var conv Conversation
	err := row.Scan(
		&conv.Id,
		&conv.Service_id,
		&conv.Proposal_id,
		&conv.Client_id,
		&conv.Analyst_id,
		&conv.Created_at,
	)
	return conv, err
}

// Função auxiliar que escaneia Mensagens
func ScanMessage(row pgx.Row) (Message, error) {
	var m Message
	err := row.Scan(
		&m.Id,
		&m.Conversation_id,
		&m.Sender_id,
		&m.Content,
		&m.Created_at,
		&m.Read_at)
	return m, err
}

func ScanRows[T any](c *gin.Context, rows pgx.Rows, scan func(pgx.Row) (T, error)) ([]T, error) {
	var results []T
	for rows.Next() {
		item, err := scan(rows)
		if HandleErr(c, err) {
			return nil, rows.Err()
		}
		results = append(results, item)
	}
	return results, nil
}
