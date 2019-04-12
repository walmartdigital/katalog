package domain

// Service ...
type Service struct {
	ID        string `json:",omitempty"`
	Name      string `json:",omitempty"`
	Port      int    `json:",omitempty"`
	Address   string `json:",omitempty"`
	Instances []Instance
}

// AddInstance ...
func (s *Service) AddInstance(endpoint Instance) {
	s.Instances = append(s.Instances, endpoint)
}
