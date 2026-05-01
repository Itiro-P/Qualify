package pkg

import "time"

type User struct {
	Id            int       `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	Time_created  time.Time `json:"time_created"`
	Country_code  string    `json:"country_code"`
	Country_name  string    `json:"country_name"`
	Country_state string    `json:"country_state"`
	City          string    `json:"city"`
	Timezone      string    `json:"timezone"`
}

type UserUpdateRequest struct {
	// Devemos deixar mudar o nome/email/telefone?
	Name          *string `json:"name,omitempty"`
	Email         *string `json:"email,omitempty"`
	Phone         *string `json:"phone,omitempty"`
	Country_code  *string `json:"country_code,omitempty"`
	Country_name  *string `json:"country_name,omitempty"`
	Country_state *string `json:"country_state,omitempty"`
	City          *string `json:"city,omitempty"`
	Timezone      *string `json:"timezone,omitempty"`
}

type UserProfile struct {
	User_id   int    `json:"user_id"`
	Biography string `json:"biography"`
}

type UserResponse struct {
	User User `json:"user"`
}

type UsersResponse struct {
	Users []User `json:"users"`
	Count int    `json:"count"`
}

type UserProfileResponse struct {
	User_profile UserProfile `json:"user_profile"`
}
