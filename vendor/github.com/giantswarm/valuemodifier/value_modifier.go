package valuemodifier

import (
	"encoding/json"

	microerror "github.com/giantswarm/microkit/error"
	"github.com/spf13/cast"
	yaml "gopkg.in/yaml.v1"
)

// Config represents the configuration used to create a new value modifier
// traverser.
type Config struct {
	// Dependencies.
	ValueModifiers []ValueModifier

	// Settings.
	IgnoreFields []string
}

// DefaultConfig provides a default configuration to create a new GPG decryption
// decoding value modifier by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		ValueModifiers: nil,

		// Settings.
		IgnoreFields: nil,
	}
}

// New creates a new configured GPG decryption value modifier.
func New(config Config) (*Service, error) {
	// Dependencies.
	if len(config.ValueModifiers) == 0 {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.ValueModifiers must not be empty")
	}

	newService := &Service{
		// Dependencies.
		valueModifiers: config.ValueModifiers,

		// Settings.
		ignoreFields: config.IgnoreFields,
	}

	return newService, nil
}

// Service implements the traversing value modifier.
type Service struct {
	// Dependencies.
	valueModifiers []ValueModifier

	// Settings.
	ignoreFields []string
}

func (s *Service) TraverseJSON(input []byte) ([]byte, error) {
	var m map[string]interface{}
	err := json.Unmarshal(input, &m)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	for k, v := range m {
		m[k] = toModifiedValueJSON(k, v, s.ignoreFields, s.valueModifiers...)
	}
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	return b, nil
}

func (s *Service) TraverseYAML(input []byte) ([]byte, error) {
	var m map[interface{}]interface{}
	err := yaml.Unmarshal(input, &m)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	for k, v := range m {
		m[k] = toModifiedValueYAML(k, v, s.ignoreFields, s.valueModifiers...)
	}
	b, err := yaml.Marshal(m)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	return b, nil
}

func toModifiedValueJSON(key string, val interface{}, ignoreFields []string, valueModifiers ...ValueModifier) interface{} {
	m1, ok := val.(map[string]interface{})
	if ok {
		for k, v := range m1 {
			m1[k] = toModifiedValueJSON(k, v, ignoreFields, valueModifiers...)
		}

		return m1
	}

	m2, ok := val.([]interface{})
	if ok {
		for i, v := range m2 {
			m2[i] = toModifiedValueJSON("", v, ignoreFields, valueModifiers...)
		}

		return m2
	}

	s := cast.ToString(val)
	if s != "" {
		for _, f := range ignoreFields {
			if f == key {
				return s
			}
		}
		for _, m := range valueModifiers {
			o, err := m.Modify([]byte(s))
			if err != nil {
				panic(err)
			}
			s = string(o)
		}
	}

	return s
}

func toModifiedValueYAML(key interface{}, val interface{}, ignoreFields []string, valueModifiers ...ValueModifier) interface{} {
	m1, ok := val.(map[interface{}]interface{})
	if ok {
		for k, v := range m1 {
			m1[k] = toModifiedValueYAML(k, v, ignoreFields, valueModifiers...)
		}

		return m1
	}

	m2, ok := val.([]interface{})
	if ok {
		for i, v := range m2 {
			m2[i] = toModifiedValueYAML("", v, ignoreFields, valueModifiers...)
		}

		return m2
	}

	s := cast.ToString(val)
	if s != "" {
		var m3 map[interface{}]interface{}
		err := yaml.Unmarshal([]byte(s), &m3)
		if err != nil || m3 == nil {
			for _, f := range ignoreFields {
				if f == key {
					return s
				}
			}
			for _, m := range valueModifiers {
				o, err := m.Modify([]byte(s))
				if err != nil {
					panic(err)
				}
				s = string(o)
			}
		} else {
			var m4 map[string]interface{}
			err := json.Unmarshal([]byte(s), &m4)
			if err != nil || m4 == nil {
				for k, v := range m3 {
					m3[k] = toModifiedValueYAML(k, v, ignoreFields, valueModifiers...)
				}

				b, err := yaml.Marshal(m3)
				if err != nil {
					panic(err)
				}

				return string(b)
			} else {
				for k, v := range m4 {
					m4[k] = toModifiedValueJSON(k, v, ignoreFields, valueModifiers...)
				}

				b, err := json.MarshalIndent(m4, "", "  ")
				if err != nil {
					panic(err)
				}

				return string(b)
			}
		}
	}

	return s
}
