package xml

import (
	"encoding/xml"

	"github.com/RavenHuo/go-kit/encode"

)

type xmlEncoder struct{}

func (x xmlEncoder) Encode(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

func (x xmlEncoder) Decode(d []byte, v interface{}) error {
	return xml.Unmarshal(d, v)
}

func (x xmlEncoder) String() string {
	return "xml"
}

// NewEncoder is xml encoder
func NewEncoder() encode.Encoder {
	return xmlEncoder{}
}
