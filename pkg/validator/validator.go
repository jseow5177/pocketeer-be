package validator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type Validator interface {
	Validate(value interface{}) error
}

// ========== Used to unset zero values ========== //

type unsetValue struct{}

func (uv *unsetValue) Error() string {
	return ""
}

// ========== Form Field Validator ========== //

type Form struct {
	Optional   bool
	Validators map[string]Validator
}

func MustForm(validators map[string]Validator) Validator {
	return &Form{
		Validators: validators,
	}
}

func OptionalForm(validators map[string]Validator) Validator {
	return &Form{
		Optional:   true,
		Validators: validators,
	}
}

func (f *Form) Validate(form interface{}) error {
	if form == nil {
		if f.Optional {
			return nil
		}
		return errors.New("form is required")
	}

	// no check needed
	if f.Validators == nil || len(f.Validators) == 0 {
		return nil
	}

	// form must be a struct
	fv := reflect.Indirect(reflect.ValueOf(form))
	if fv.Kind() != reflect.Struct {
		return errors.New("form must be a struct")
	}

	// validate each struct field
	for i := 0; i < fv.NumField(); i++ {
		fieldVal := fv.Field(i)
		fieldName := getFieldName(fv.Type().Field(i))

		if v, ok := f.Validators[fieldName]; ok && v != nil {
			var err error

			if !fieldVal.CanInterface() {
				return errors.New("unsupported field value")
			}

			kind := fieldVal.Kind()
			if kind == reflect.Ptr {
				if fieldVal.IsNil() {
					err = v.Validate(nil)
				} else {
					err = v.Validate(reflect.Indirect(fieldVal).Interface())
				}
			} else {
				err = v.Validate(fieldVal.Interface())
			}

			if _, ok := err.(*unsetValue); ok {
				fieldVal.Set(reflect.Zero(fieldVal.Type()))
				err = nil
			}

			if err != nil {
				return fmt.Errorf("%s: %v", fieldName, err)
			}
		}
	}

	return nil
}

// ========== Slice Field Validator ========== //

type Slice struct {
	Optional  bool
	MinLen    uint32
	MaxLen    uint32
	Validator Validator
}

func (sv *Slice) Validate(value interface{}) error {
	if value == nil {
		if sv.Optional {
			return nil
		}
		return errors.New("slice field is required")
	}

	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Slice {
		return errors.New("unexpected non-slice type")
	}

	if val.IsNil() || val.Len() == 0 {
		if sv.Optional {
			return nil
		}
		return errors.New("slice cannot be empty")
	}

	valLen := uint32(val.Len())
	if valLen < sv.MinLen {
		return fmt.Errorf("requires at least %d elements, it has %d", sv.MinLen, valLen)
	}
	if sv.MaxLen > 0 && valLen > sv.MaxLen {
		return fmt.Errorf("allows at most %d elements, it has %d", sv.MaxLen, valLen)
	}

	if sv.Validator == nil {
		return nil
	}

	// run validator through each slice element
	for i := 0; i < val.Len(); i++ {
		var err error

		elem := val.Index(i)
		if !elem.CanInterface() {
			return fmt.Errorf("unsupported element in slice")
		}

		kind := elem.Kind()
		if kind == reflect.Ptr {
			if elem.IsNil() {
				err = sv.Validator.Validate(nil)
			} else {
				err = sv.Validator.Validate(reflect.Indirect(elem).Interface())
			}
		} else {
			err = sv.Validator.Validate(elem.Interface())
		}

		if err != nil {
			return fmt.Errorf("index %d: %v", i, err)
		}
	}

	return nil
}

// ========== String Field Validator ========== //

type StringFunc func(string) error

type String struct {
	Optional   bool
	UnsetZero  bool
	MinLen     uint32
	MaxLen     uint32
	Charset    string
	Regex      *regexp.Regexp
	Validators []StringFunc
}

func (stv *String) Validate(value interface{}) error {
	if value == nil {
		if stv.Optional {
			return nil
		}
		return errors.New("string field is required")
	}

	str, ok := value.(string)
	if !ok {
		return errors.New("unexpected non-string type")
	}

	if str == "" {
		if stv.Optional {
			if stv.UnsetZero {
				return &unsetValue{}
			}
			return nil
		}
		return errors.New("string field cannot be empty")
	}

	runes := []rune(str)
	runesLen := uint32(len(runes))
	if runesLen < stv.MinLen {
		return fmt.Errorf("requires at least %d chars, it has %d", stv.MinLen, runesLen)
	}
	if stv.MaxLen > 0 && runesLen > stv.MaxLen {
		return fmt.Errorf("allows at most %d chars, it has %d", stv.MaxLen, runesLen)
	}

	if len(stv.Charset) > 0 {
		for i := 0; i < len(runes); i++ {
			if !strings.ContainsRune(stv.Charset, runes[i]) {
				return fmt.Errorf("must contain only chars '%v'", stv.Charset)
			}
		}
	}

	for _, v := range stv.Validators {
		if v != nil {
			if err := v(str); err != nil {
				return err
			}
		}
	}

	return nil
}

// ========== Uint64 Field Validator ========== //

type Uint64Func func(uint64) error

type Uint64 struct {
	Optional   bool
	UnsetZero  bool
	Min        *uint64
	Max        *uint64
	Validators []Uint64Func
}

func (uv *Uint64) Validate(value interface{}) error {
	if value == nil {
		if uv.Optional {
			return nil
		}
		return errors.New("uint64 field is required")
	}

	ui, ok := value.(uint64)
	if !ok {
		return errors.New("unexpected non-uint64 type")
	}

	if ui == 0 {
		if uv.Optional {
			if uv.UnsetZero {
				return &unsetValue{}
			}
			return nil
		}
		return errors.New("uint64 field cannot be zero")
	}

	if uv.Min != nil && ui < *uv.Min {
		return fmt.Errorf("must be greater than or equal to %v", *uv.Min)
	}

	if uv.Max != nil && ui > *uv.Max {
		return fmt.Errorf("must be lesser than or equal to %v", *uv.Max)
	}

	for _, v := range uv.Validators {
		if v != nil {
			if err := v(ui); err != nil {
				return err
			}
		}
	}

	return nil
}

func getFieldName(structField reflect.StructField) string {
	jsonTag := structField.Tag.Get("json")
	switch jsonTag {
	case "-":
		return ""
	case "":
		return structField.Name
	default:
		parts := strings.Split(jsonTag, ",")
		if len(parts) == 0 {
			return structField.Name
		}

		jsonFieldName := parts[0]
		if jsonFieldName == "" {
			return structField.Name
		}

		return jsonFieldName
	}
}
