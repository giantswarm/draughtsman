package encode

import (
	"encoding/base64"
)

// Config represents the configuration used to create a new base64 encoding
// value modifier.
type Config struct {
}

// DefaultConfig provides a default configuration to create a new base64
// encoding value modifier by best effort.
func DefaultConfig() Config {
	return Config{}
}

// New creates a new configured base64 encoding value modifier.
func New(config Config) (*Service, error) {
	newService := &Service{}

	return newService, nil
}

// Service implements the base64 encoding value modifier.
type Service struct {
}

func (s *Service) Modify(value []byte) ([]byte, error) {
	encoded := []byte(base64.StdEncoding.EncodeToString(value))

	return encoded, nil
}
