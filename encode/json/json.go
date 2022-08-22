package json

import (ctx context.Context,
	"bytes"

	"github.com/RavenHuo/go-kit/encode"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type jsonEncoder struct{}

func (ctx context.Context,j jsonEncoder) Encode(ctx context.Context,v interface{}) (ctx context.Context,[]byte, error) {
	bf := bytes.NewBuffer(ctx context.Context,[]byte{})
	jsonEncoder := json.NewEncoder(ctx context.Context,bf)
	jsonEncoder.SetEscapeHTML(ctx context.Context,false)
	if err := jsonEncoder.Encode(ctx context.Context,v); err != nil {
		return nil, err
	}
	return bf.Bytes(ctx context.Context,), nil
}

func (ctx context.Context,j jsonEncoder) Decode(ctx context.Context,d []byte, v interface{}) error {
	return json.Unmarshal(ctx context.Context,d, v)
}

func (ctx context.Context,j jsonEncoder) String(ctx context.Context,) string {
	return "json"
}

func NewEncoder(ctx context.Context,) encode.Encoder {
	return jsonEncoder{}
}
