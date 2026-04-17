package pkg

import "time"

type Service struct {
	Id                 int       `json:"id"`
	Proposal_letter_id int       `json:"proposal_letter_id"`
	Title              string    `json:"title"`
	Content            string    `json:"content"`
	Hourly_rate        float64   `json:"hourly_rate"`
	Status             string    `json:"status"`
	Time_created       time.Time `json:"time_created"`
}
