package pkg

import "fmt"

// Adiciona % no início e final de uma string.
func PutPercent(str string) string {
	if str == "" {
		return ""
	}
	return "%" + str + "%"
}

func mergeMaps(maps ...map[string]bool) map[string]bool {
	result := map[string]bool{}
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

type Pagination struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

func (p *Pagination) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 || p.PageSize > 100 {
		p.PageSize = 20
	}
}

func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

type PaginatedResponse struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Count    int `json:"count"`
}

type SortOptions struct {
	SortBy string `form:"sort_by"`
	Order  string `form:"order"`
}

type RatingFilter struct {
	Rating    *int `form:"rating"`
	MinRating *int `form:"min_rating"`
	MaxRating *int `form:"max_rating"`
}

func (s *SortOptions) ValidateSort(allowed map[string]bool) string {
	if !allowed[s.SortBy] {
		return ""
	}
	if s.Order != "ASC" && s.Order != "DESC" {
		s.Order = "ASC"
	}
	return fmt.Sprintf("%s %s", s.SortBy, s.Order)
}

type UserFilter struct {
	Pagination
	Name         string `form:"name"`
	Email        string `form:"email"`
	Country      string `form:"country"`
	CountryCode  string `form:"country_code"`
	CountryState string `form:"country_state"`
	City         string `form:"city"`
	Timezone     string `form:"timezone"`
}

type AnalystFilter struct {
	RatingFilter
	UserFilter
	SortOptions
	Skills          string   `form:"skills"`
	MinHourlyRate   *float64 `form:"min_hourly_rate"`
	MaxHourlyRate   *float64 `form:"max_hourly_rate"`
	MinTotalReviews *int     `form:"min_total_reviews"`
	MaxTotalReviews *int     `form:"max_total_reviews"`
}

type ClientFilter struct {
	UserFilter
	SortOptions
	MinProposedBudget *float64 `form:"min_proposed_budget"`
	MaxProposedBudget *float64 `form:"max_proposed_budget"`
}

type CertificationFilter struct {
	Pagination
	SortOptions
	Name        string `form:"name"`
	Institution string `form:"institution"`
	Year        *int   `form:"year"`
	FromYear    *int   `form:"from_year"`
	ToYear      *int   `form:"to_year"`
}

type SkillFilter struct {
	Pagination
	SortOptions
	Name string `form:"name"`
}

type ProposalFilter struct {
	Pagination
	SortOptions
	AnalystId             *int     `form:"analyst_id"`
	ClientId              *int     `form:"client_id"`
	Title                 string   `form:"title"`
	Content               string   `form:"content"`
	MinProposedHourlyRate *float64 `form:"min_proposed_hourly_rate"`
	MaxProposedHourlyRate *float64 `form:"max_proposed_hourly_rate"`
}

type ServiceFilter struct {
	Pagination
	SortOptions
	ProposalId    *int     `form:"proposal_id"`
	Title         string   `form:"title"`
	Content       string   `form:"content"`
	MinHourlyRate *float64 `form:"min_hourly_rate"`
	MaxHourlyRate *float64 `form:"max_hourly_rate"`
	Status        string   `form:"status"`
}

type ReviewFilter struct {
	Pagination
	RatingFilter
	SortOptions
	AnalystId *int   `form:"analyst_id"`
	ClientId  *int   `form:"client_id"`
	ServiceId *int   `form:"service_id"`
	Comment   string `form:"comment"`
}

var UserSortFields = map[string]bool{
	"name":          true,
	"country_name":  true,
	"country_state": true,
	"city":          true,
	"time_created":  true,
}

var AnalystSortFields = mergeMaps(UserSortFields, map[string]bool{
	"hourly_rate":   true,
	"total_reviews": true,
	"mean_rating":   true,
})

var ClientSortFields = mergeMaps(UserSortFields, map[string]bool{
	"proposed_budget": true,
})

var CertificationSortFields = map[string]bool{
	"name":        true,
	"year":        true,
	"institution": true,
}

var SkillSortFields = map[string]bool{
	"name": true,
}

var ProposalSortFields = map[string]bool{
	"title":                true,
	"proposed_hourly_rate": true,
	"time_created":         true,
}

var ServiceSortFields = map[string]bool{
	"title":        true,
	"hourly_rate":  true,
	"status":       true,
	"time_created": true,
}

var ReviewSortFields = map[string]bool{
	"rating":       true,
	"time_created": true,
}
