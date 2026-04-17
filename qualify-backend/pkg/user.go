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
