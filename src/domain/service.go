package domain

// Service ...
type Service struct {
	ID        string `json:",omitempty"`
	Name      string `json:",omitempty"`
	Port      int    `json:",omitempty"`
	Address   string `json:",omitempty"`
	Endpoints []Endpoint
}

// AddEndpoint ...
func (s *Service) AddEndpoint(endpoint Endpoint) {
	s.Endpoints = append(s.Endpoints, endpoint)
}

// Endpoint ...
type Endpoint struct {
	Address string `json:",omitempty"`
}
