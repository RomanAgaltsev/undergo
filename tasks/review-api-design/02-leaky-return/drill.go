// Package drill — C6/02 leaky-return.
package drill

// Config is the exported service configuration.
type Config struct {
	Timeout int
}

// result is the internal outcome of a Load.
type result struct {
	Code int
	Body string
}

// Service serves load requests.
type Service struct {
	cfg *Config
}

// Load fetches a resource and returns the result.
func (s *Service) Load() *result {
	return &result{Code: 200, Body: "ok"}
}

// Config returns the service configuration.
func (s *Service) Config() *Config {
	return s.cfg
}
