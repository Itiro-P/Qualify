package pkg

import "time"

type MessageCreateRequest struct {
	Sender_id int    `json:"sender_id"`
	Content   string `json:"content"`
}

type MessageUpdateRequest struct {
	Content string `json:"content,omitempty"`
}

type Message struct {
	MessageCreateRequest
	Id              int        `json:"id"`
	Conversation_id int        `json:"conversation_id"`
	Created_at      time.Time  `json:"created_at"`
	Read_at         *time.Time `json:"read_at,omitempty"`
}

type MessageResponse struct {
	Messages []Message `json:"messages"`
	Count    int       `json:"count"`
}
