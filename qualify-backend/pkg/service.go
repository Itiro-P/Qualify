package pkg

type Service struct {
	ID                 int    `json:"id"`
	Proposal_letter_id int    `json:"proposal_letter_id"`
	Title              string `json:"title"`
	Content            string `json:"content"`
	Hourly_rate        int    `json:"hourly_rate"`
	Status             string `json:"status"`
	Time_created       string `json:"time_created"`
}
