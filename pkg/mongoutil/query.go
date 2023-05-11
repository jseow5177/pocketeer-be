package mongoutil

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	filterTag = "filter"

	_id = "_id"
	sep = "__"
	eq  = "eq"
	ne  = "ne"
	gt  = "gt"
	gte = "gte"
	lt  = "lt"
	lte = "lte"
	in  = "in"
	nin = "nin"
)

var supportedOps = map[string]string{
	eq:  "equal",
	ne:  "not equal",
	gt:  "greater than",
	gte: "greater than equal",
	lt:  "less than",
	lte: "less than equal",
	in:  "in",
	nin: "not in",
}

func getOp(op string) string {
	return fmt.Sprintf("$%s", op)
}

func BuildFilter(filter interface{}) bson.M {
	if filter == nil {
		return nil
	}

	val := reflect.ValueOf(filter)
	if val.IsNil() {
		return nil
	}
	val = val.Elem()

	conds := make(bson.M)
	for i := 0; i < val.NumField(); i++ {
		fn := val.Type().Field(i).Tag.Get(filterTag) // filter name
		fv := reflect.Indirect(val.Field(i))         // filter value
		fk := fv.Kind()                              // filter type

		// handle nil pointer
		if fk == reflect.Invalid {
			continue
		}

		// handle empty slice
		if fk == reflect.Slice && fv.Len() == 0 {
			continue
		}

		// field and operator
		parts := strings.SplitN(fn, sep, 2)

		// operator
		var op string
		if len(parts) > 1 {
			op = parts[1]
			if _, ok := supportedOps[op]; !ok {
				continue
			}
		} else {
			// default to eq
			op = eq
		}
		op = getOp(op)

		// one condition
		cond := make(bson.M)
		cond[op] = fv.Interface()

		// handle _id field
		f := strcase.ToSnake(parts[0])
		if f == _id {
			id := fmt.Sprint(cond[op])
			if primitive.IsValidObjectID(id) {
				objID, _ := primitive.ObjectIDFromHex(id)
				cond[op] = objID
			}
		}

		conds[f] = cond
	}

	return conds
}
