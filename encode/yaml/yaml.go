package yaml

import (
	"github.com/RavenHuo/go-pkg/encode"
	"gopkg.in/yaml.v3"
)

type yamlEncoder struct{}

func (y yamlEncoder) Encode(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func (y yamlEncoder) Decode(d []byte, v interface{}) error {
	return yaml.Unmarshal(d, v)
}

func (y yamlEncoder) Name() string {
	return "yaml"
}

// NewEncoder is yaml encoder
func NewEncoder() encode.Encoder {
	return yamlEncoder{}
}
