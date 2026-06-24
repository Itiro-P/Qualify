package pkg

import "time"

type ConversationCreateRequest struct {
	Service_id  *int `json:"service_id,omitempty"`
	Proposal_id *int `json:"proposal_id,omitempty"`
	Analyst_id  int  `json:"analyst_id"`
	Client_id   int  `json:"client_id"`
}

type ConversationUpdateRequest struct {
	Service_id  *int `json:"service_id,omitempty"`
	Proposal_id *int `json:"proposal_id,omitempty"`
	Analyst_id  *int `json:"analyst_id,omitempty"`
	Client_id   *int `json:"client_id,omitempty"`
}

type Conversation struct {
	ConversationCreateRequest
	Id         int       `json:"id"`
	Created_at time.Time `json:"created_at"`
}

type ConversationResponse struct {
	Conversations []Conversation `json:"conversations"`
	Count         int            `json:"count"`
}
