package pkg

type Client struct {
	User
	Proposed_budget float64 `json:"proposed_budget"`
}

type ClientCreateRequest struct {
	Proposed_budget float64 `json:"proposed_budget"`
}

type ClientUpdateRequest struct {
	UserUpdateRequest
	Proposed_budget *float64 `json:"proposed_budget,omitempty"`
}

type ClientProfile struct {
	UserProfile
}

type ClientsResponse struct {
	Clients []Client `json:"clients"`
	Count   int      `json:"count"`
}

type ClientResponse struct {
	Client Client `json:"client"`
}

type ClientProfileResponse struct {
	Client_profile ClientProfile `json:"client_profile"`
}
