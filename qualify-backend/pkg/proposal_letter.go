package pkg

import "time"

type ProposalLetter struct {
	Id                   int       `json:"id"`
	Title                string    `json:"title"`
	Content              string    `json:"content"`
	Client_id            int       `json:"client_id"`
	Analyst_id           int       `json:"analyst_id"`
	Proposed_hourly_rate float64   `json:"proposed_hourly_rate"`
	Time_created         time.Time `json:"time_created"`
}
