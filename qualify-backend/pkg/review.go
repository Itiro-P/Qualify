package pkg

import "time"

type ReviewCreateRequest struct {
	Analyst_id int    `json:"analyst_id"`
	Client_id  int    `json:"client_id"`
	Service_id int    `json:"service_id"`
	Rating     int    `json:"rating"`
	Comment    string `json:"comment"`
}

type Review struct {
	ReviewCreateRequest
	Id           int       `json:"id"`
	Time_created time.Time `json:"time_created"`
}

// Devemos permitir atualizar a avaliação? Talvez só o comentário, ou nem isso?
type ReviewUpdateRequest struct {
	Rating  *int    `json:"rating,omitempty"`
	Comment *string `json:"comment,omitempty"`
}

type ReviewResponse struct {
	Review Review `json:"review"`
}

type ReviewsResponse struct {
	Reviews   []Review `json:"reviews"`
	Count     int      `json:"count"`
	Page      int      `json:"page"`
	Page_size int      `json:"page_size"`
}
