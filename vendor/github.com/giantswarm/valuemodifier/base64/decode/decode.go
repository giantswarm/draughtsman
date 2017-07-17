package decode

import (
	"encoding/base64"

	microerror "github.com/giantswarm/microkit/error"
)

// Config represents the configuration used to create a new base64 decoding
// value modifier.
type Config struct {
}

// DefaultConfig provides a default configuration to create a new base64
// decoding value modifier by best effort.
func DefaultConfig() Config {
	return Config{}
}

// New creates a new configured base64 decoding value modifier.
func New(config Config) (*Service, error) {
	newService := &Service{}

	return newService, nil
}

// Service implements the base64 decoding value modifier.
type Service struct {
}

func (s *Service) Modify(value []byte) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(string(value))
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	return decoded, nil
}
