package pkg

type Skill struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type SkillResponse struct {
	Skill Skill `json:"skill"`
}

type SkillsResponse struct {
	Skills []Skill `json:"skills"`
	Count  int     `json:"count"`
}
