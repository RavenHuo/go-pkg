// Package encoder handles source encoding formats
package encode

// Encoder represents a format encoder
type Encoder interface {
	Encode(interface{}) ([]byte, error)
	Decode([]byte, interface{}) error
	Name() string
}
