package pkg

type Certification struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Year        int    `json:"year"`
	Description string `json:"description"`
	Institution string `json:"institution"`
}

type CertificationUpdateRequest struct {
	Name        *string `json:"name,omitempty"`
	Year        *int    `json:"year,omitempty"`
	Description *string `json:"description,omitempty"`
	Institution *string `json:"institution,omitempty"`
}

type CertificationsResponse struct {
	Certifications []Certification `json:"certifications"`
	Count          int             `json:"count"`
}

type CertificationResponse struct {
	Certification Certification `json:"certification"`
}
