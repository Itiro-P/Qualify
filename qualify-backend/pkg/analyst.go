package pkg

type Analyst struct {
	User
	Hourly_rate   float64  `json:"hourly_rate"`
	Total_reviews int      `json:"total_reviews"`
	Mean_rating   *float64 `json:"mean_rating"`
}

type AnalystCreateRequest struct {
	Hourly_rate float64 `json:"hourly_rate"`
}

type AnalystProfile struct {
	UserProfile
}

type AnalystSkill struct {
	Analyst_id int `json:"analyst_id,omitempty"`
	Skill_id   int `json:"skill_id"`
}

type AnalystCertification struct {
	Analyst_id       int `json:"analyst_id,omitempty"`
	Certification_id int `json:"certification_id"`
}

type AnalystsResponse struct {
	Analysts []Analyst `json:"analysts"`
	Count    int       `json:"count"`
}

type AnalystResponse struct {
	Analyst Analyst `json:"analyst"`
}

type AnalystProfileResponse struct {
	Analyst_profile AnalystProfile `json:"analyst_profile"`
}

type AnalystSkillResponse struct {
	Analyst_skill AnalystSkill `json:"analyst_skill"`
}

type AnalystSkillsResponse struct {
	Analyst_skills []AnalystSkill `json:"analyst_skills"`
	Count          int            `json:"count"`
}

type AnalystCertificationResponse struct {
	Analyst_certification AnalystCertification `json:"analyst_certification"`
}

type AnalystUpdateRequest struct {
	UserUpdateRequest
	Hourly_rate   *float64 `json:"hourly_rate,omitempty"`
	Total_reviews *int     `json:"total_reviews,omitempty"`
	Mean_rating   *float64 `json:"mean_rating,omitempty"`
}
