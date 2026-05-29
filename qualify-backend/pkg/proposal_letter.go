package pkg

import "time"

type ProposalLetter struct {
	Id                   int       `json:"id,omitempty"`
	Title                string    `json:"title"`
	Content              string    `json:"content"`
	Client_id            int       `json:"client_id"`
	Analyst_id           int       `json:"analyst_id"`
	Proposed_hourly_rate float64   `json:"proposed_hourly_rate"`
	Time_created         time.Time `json:"time_created"`
}

type ProposalLetterUpdateRequest struct {
	Title                *string  `json:"title,omitempty"`
	Content              *string  `json:"content,omitempty"`
	Proposed_hourly_rate *float64 `json:"proposed_hourly_rate,omitempty"`
}

type ProposalLetterResponse struct {
	Proposal_letter ProposalLetter `json:"proposal_letter"`
}

type ProposalLettersResponse struct {
	Proposal_letters []ProposalLetter `json:"proposal_letters"`
	Count            int              `json:"count"`
}
