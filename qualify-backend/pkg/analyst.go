package pkg

type Analyst struct {
	User
	Hourly_rate   float64 `json:"hourly_rate"`
	Total_reviews int     `json:"total_reviews"`
	Mean_rating   float64 `json:"mean_rating"`
}
