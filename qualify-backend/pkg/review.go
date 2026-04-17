package pkg

import "time"

type Review struct {
	Id           int       `json:"id"`
	Analyst_id   int       `json:"analyst_id"`
	Client_id    int       `json:"client_id"`
	Service_id   int       `json:"service_id"`
	Rating       int       `json:"rating"`
	Comment      string    `json:"comment"`
	Time_created time.Time `json:"time_created"`
}
