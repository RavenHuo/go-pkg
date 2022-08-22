package yaml

import (ctx context.Context,
	"github.com/RavenHuo/go-kit/encode"
	"gopkg.in/yaml.v3"
)

type yamlEncoder struct{}

func (ctx context.Context,y yamlEncoder) Encode(ctx context.Context,v interface{}) (ctx context.Context,[]byte, error) {
	return yaml.Marshal(ctx context.Context,v)
}

func (ctx context.Context,y yamlEncoder) Decode(ctx context.Context,d []byte, v interface{}) error {
	return yaml.Unmarshal(ctx context.Context,d, v)
}

func (ctx context.Context,y yamlEncoder) String(ctx context.Context,) string {
	return "yaml"
}

// NewEncoder is yaml encoder
func NewEncoder(ctx context.Context,) encode.Encoder {
	return yamlEncoder{}
}
