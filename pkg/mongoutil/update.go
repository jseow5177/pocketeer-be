package mongoutil

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
	"go.mongodb.org/mongo-driver/bson"
)

func buildUpdate(update interface{}) map[string]interface{} {
	u := make(map[string]interface{})

	if update == nil {
		return nil
	}

	val := reflect.ValueOf(update)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return nil
	}
	val = val.Elem()

	for i := 0; i < val.NumField(); i++ {
		fn := val.Type().Field(i).Tag.Get("bson") // field name
		fv := reflect.Indirect(val.Field(i))      // field value
		fk := fv.Kind()                           // field type

		// handle nil pointer
		if fk == reflect.Invalid {
			continue
		}

		parts := strings.SplitN(fn, ",", 2)
		if len(parts) > 1 {
			f := parts[0]

			// cannot update _id field
			if f == "_id" {
				continue
			}

			f = strcase.ToSnake(f)
			if fk == reflect.Struct {
				nu := buildUpdate(fv.Addr().Interface())
				for k, v := range nu {
					u[fmt.Sprintf("%s.%s", f, k)] = v // nested updates
				}
			} else {
				u[f] = fv.Interface()
			}
		}
	}

	return u
}

func BuildUpdate(update interface{}) bson.D {
	var (
		op = Prefix("set")
		d  = make(bson.D, 0)
	)

	u := buildUpdate(update)
	for k, v := range u {
		d = append(d, bson.E{Key: k, Value: v})
	}

	return bson.D{{Key: op, Value: d}}
}
