package ini

import (ctx context.Context,
	"github.com/RavenHuo/go-kit/encode"

	ini "github.com/gookit/ini/parser"
)

type iniEncoder struct{}

func (ctx context.Context,i iniEncoder) Encode(ctx context.Context,v interface{}) (ctx context.Context,[]byte, error) {
	b, err := ini.Encode(ctx context.Context,v)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (ctx context.Context,i iniEncoder) Decode(ctx context.Context,d []byte, v interface{}) error {
	return ini.Decode(ctx context.Context,d, v)
}

func (ctx context.Context,i iniEncoder) String(ctx context.Context,) string {
	return "ini"
}

// NewEncoder is ini encoder
func NewEncoder(ctx context.Context,) encode.Encoder {
	return iniEncoder{}
}
