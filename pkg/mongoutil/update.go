package mongoutil

import (
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
	"go.mongodb.org/mongo-driver/bson"
)

func BuildUpdate(update interface{}) bson.D {
	if update == nil {
		return nil
	}

	val := reflect.ValueOf(update)
	if val.IsNil() {
		return nil
	}
	val = val.Elem()

	d := make(bson.D, 0)
	for i := 0; i < val.NumField(); i++ {
		fn := val.Type().Field(i).Tag.Get(tagBson) // field name
		fv := reflect.Indirect(val.Field(i))       // field value
		fk := fv.Kind()                            // field type

		// handle nil pointer
		if fk == reflect.Invalid {
			continue
		}

		parts := strings.SplitN(fn, ",", 2)

		if len(parts) > 1 {
			f := parts[0]
			// cannot update _id field
			if f != _id {
				f = strcase.ToSnake(f)
				d = append(d, bson.E{Key: f, Value: fv.Interface()})
			}
		}
	}

	op := getOp(set)

	return bson.D{{Key: op, Value: d}}
}
