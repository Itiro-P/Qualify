package pkg

type ProposalLetter struct {
	ID                   int     `json:"id"`
	Title                string  `json:"title"`
	Content              string  `json:"content"`
	Client_id            string  `json:"client_id"`
	Analyst_id           string  `json:"analyst_id"`
	Proposed_hourly_rate float64 `json:"proposed_hourly_rate"`
	Time_created         string  `json:"time_created"`
}
