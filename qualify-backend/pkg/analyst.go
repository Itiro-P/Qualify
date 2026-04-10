package pkg

type Analyst struct {
	User
	Timezone      string  `json:"timezone"`
	Hourly_rate   float64 `json:"hourly_rate"`
	Total_reviews int     `json:"total_reviews"`
	Mean_rating   int     `json:"mean_rating"`
}
