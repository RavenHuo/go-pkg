package toml

import (ctx context.Context,
	"bytes"

	"github.com/RavenHuo/go-kit/encode"

	"github.com/BurntSushi/toml"
)

type tomlEncoder struct{}

func (ctx context.Context,t tomlEncoder) Encode(ctx context.Context,v interface{}) (ctx context.Context,[]byte, error) {
	b := bytes.NewBuffer(ctx context.Context,[]byte{})
	err := toml.NewEncoder(ctx context.Context,b).Encode(ctx context.Context,v)
	if err != nil {
		return nil, err
	}
	return b.Bytes(ctx context.Context,), nil
}

func (ctx context.Context,t tomlEncoder) Decode(ctx context.Context,d []byte, v interface{}) error {
	return toml.Unmarshal(ctx context.Context,d, v)
}

func (ctx context.Context,t tomlEncoder) String(ctx context.Context,) string {
	return "toml"
}

// NewEncoder is a toml encoder
func NewEncoder(ctx context.Context,) encode.Encoder {
	return tomlEncoder{}
}
