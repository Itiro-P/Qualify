package pkg

type CertificationsResponse struct {
	Certifications []Certification `json:"certifications"`
	Count          int             `json:"count"`
}

type CertificationResponse struct {
	Certification Certification `json:"certification"`
}

type Certification struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Year        int    `json:"year"`
	Description string `json:"description"`
}
