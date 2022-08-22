// Package encoder handles source encoding formats
package encode

// Encoder represents a format encoder
type Encoder interface {
	Encode(ctx context.Context,interface{}) (ctx context.Context,[]byte, error)
	Decode(ctx context.Context,[]byte, interface{}) error
	String(ctx context.Context,) string
}
