package xml

import (ctx context.Context,
	"encoding/xml"

	"github.com/RavenHuo/go-kit/encode"

)

type xmlEncoder struct{}

func (ctx context.Context,x xmlEncoder) Encode(ctx context.Context,v interface{}) (ctx context.Context,[]byte, error) {
	return xml.Marshal(ctx context.Context,v)
}

func (ctx context.Context,x xmlEncoder) Decode(ctx context.Context,d []byte, v interface{}) error {
	return xml.Unmarshal(ctx context.Context,d, v)
}

func (ctx context.Context,x xmlEncoder) String(ctx context.Context,) string {
	return "xml"
}

// NewEncoder is xml encoder
func NewEncoder(ctx context.Context,) encode.Encoder {
	return xmlEncoder{}
}
