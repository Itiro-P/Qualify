package pkg

import "time"

type Review struct {
	Id           int       `json:"id"`
	Rating       int       `json:"rating"`
	Comment      string    `json:"comment"`
	Service_id   int       `json:"service_id"`
	Time_created time.Time `json:"time_created"`
}
