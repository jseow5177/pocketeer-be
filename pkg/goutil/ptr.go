package goutil

import (
	"reflect"

	"github.com/imdario/mergo"
)

func String(s string) *string {
	return &s
}

func Uint32(ui uint32) *uint32 {
	return &ui
}

func Uint64(ui uint64) *uint64 {
	return &ui
}

func Int(i int) *int {
	return &i
}

type ptrTransformer struct{}

func (t ptrTransformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	if typ.Kind() == reflect.Ptr {
		return func(dst, src reflect.Value) error {
			// If the src ptr is nil, don't overwrite the dst ptr.
			if src.IsNil() {
				return nil
			}
			// Otherwise, allocate a new value for the dst ptr and copy the value from the src ptr.
			if dst.IsNil() {
				dst.Set(reflect.New(typ.Elem()))
			}
			dst.Elem().Set(src.Elem())
			return nil
		}
	}
	return nil
}

func MergeWithPtrFields(dst interface{}, src interface{}) error {
	return mergo.Merge(dst, src, mergo.WithTransformers(new(ptrTransformer)))
}
