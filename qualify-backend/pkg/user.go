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

type UserProfile struct {
	User_id   int    `json:"user_id"`
	Biography string `json:"biography"`
}

type UserSkill struct {
	User_id  int `json:"user_id"`
	Skill_id int `json:"skill_id"`
}

type UserResponse struct {
	User User `json:"user"`
}

type UsersResponse struct {
	Users []User `json:"users"`
	Count int    `json:"count"`
}

type UserSkillsResponse struct {
	User_skills []UserSkill `json:"user_skills"`
	Count       int         `json:"count"`
}

type UserProfileResponse struct {
	User_profile UserProfile `json:"user_profile"`
}

type UserSkillResponse struct {
	User_skill UserSkill `json:"user_skill"`
}
