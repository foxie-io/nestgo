package ng

import "fmt"

// PayloadKeyer is an interface for defining keys used to store and retrieve payloads in the context.
type PayloadKeyer interface {
	PayloadKey() string
}

// TypeKey is a key type based on generic type T
type TypeKey[T any] struct{}

func (p TypeKey[T]) PayloadKey() string {
	return fmt.Sprintf("%T", p)
}

type PayloadKey string

func (p PayloadKey) PayloadKey() string {
	return "__" + string(p) + "__"
}
