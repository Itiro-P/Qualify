package pkg

import "time"

type Service struct {
	Id                 int       `json:"id,omitempty"`
	Proposal_letter_id int       `json:"proposal_letter_id"`
	Title              string    `json:"title"`
	Content            string    `json:"content"`
	Hourly_rate        float64   `json:"hourly_rate"`
	Status             string    `json:"status"`
	Time_created       time.Time `json:"time_created"`
}

type ServiceUpdateRequest struct {
	Title       *string  `json:"title,omitempty"`
	Content     *string  `json:"content,omitempty"`
	Hourly_rate *float64 `json:"hourly_rate,omitempty"`
	Status      *string  `json:"status,omitempty"`
}

type ServiceResponse struct {
	Service Service `json:"service"`
}

type ServicesResponse struct {
	Services []Service `json:"services"`
	Count    int       `json:"count"`
}
