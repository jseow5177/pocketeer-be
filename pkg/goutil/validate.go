package goutil

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

var (
	ErrInvalidFloat = errors.New("invalid float")
)

func FormatFloat(f float64, decimalPlaces int) string {
	ft := fmt.Sprintf("%%.%df", decimalPlaces)
	return fmt.Sprintf(ft, f)
}

func IsFloat(str string, maxDecimalPoints int) error {
	regex := regexp.MustCompile(`^\d+(\.\d{0,` + strconv.Itoa(maxDecimalPoints) + `})?$`)
	if !regex.MatchString(str) {
		return ErrInvalidFloat
	}

	if _, err := strconv.ParseFloat(str, 64); err != nil {
		return err
	}

	return nil
}

func Zero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return v.IsNil()
	case reflect.Array:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && Zero(v.Index(i))
		}
		return z
	case reflect.Struct:
		z := true
		for i := 0; i < v.NumField(); i++ {
			z = z && Zero(v.Field(i))
		}
		return z
	}
	z := reflect.Zero(v.Type())
	return z.Interface() == v.Interface()
}

func All(args ...interface{}) bool {
	for _, arg := range args {
		if arg == nil {
			return false
		}
		val := reflect.ValueOf(arg)
		if !val.IsValid() {
			return false
		}
		if Zero(val) {
			return false
		}
	}
	return true
}

func Any(args ...interface{}) bool {
	for _, arg := range args {
		if arg == nil {
			continue
		}
		val := reflect.ValueOf(arg)
		if !val.IsValid() {
			continue
		}
		if !Zero(val) {
			return true
		}
	}
	return false
}

func One(args ...interface{}) bool {
	count := 0
	for _, arg := range args {
		if arg == nil {
			continue
		}
		val := reflect.ValueOf(arg)
		if !val.IsValid() {
			continue
		}
		if !Zero(val) {
			count++
			if count > 1 {
				return false
			}
		}
	}
	return count == 1
}
