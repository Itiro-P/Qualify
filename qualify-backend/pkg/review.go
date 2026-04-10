package pkg

type Review struct {
	ID           int    `json:"id"`
	Client_id    int    `json:"client_id"`
	Analyst_id   int    `json:"analyst_id"`
	Rating       int    `json:"rating"`
	Comment      string `json:"comment"`
	Time_created string `json:"time_created"`
}
