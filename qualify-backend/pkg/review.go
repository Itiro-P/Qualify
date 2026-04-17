package pkg

type Review struct {
	Id           int    `json:"id"`
	Client_id    int    `json:"client_id"`
	Analyst_id   int    `json:"analyst_id"`
	Rating       int    `json:"rating"`
	Comment      string `json:"comment"`
	Service_id   int    `json:"service_id"`
	Time_created string `json:"time_created"`
}
