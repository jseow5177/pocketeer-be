package mongoutil

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/jseow5177/pockteer-be/pkg/filter"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getOp(op string) string {
	return fmt.Sprintf("$%s", op)
}

func BuildFilterOptions(filterOptions filter.FilterOptions) *options.FindOptions {
	if filterOptions == nil {
		return nil
	}

	opts := new(options.FindOptions)

	if filterOptions.GetLimit() != nil {
		limit := *filterOptions.GetLimit()
		opts.SetLimit(int64(limit))

		if filterOptions.GetPage() != nil {
			page := *filterOptions.GetPage()
			opts.SetSkip((int64(page) - 1) * int64(limit))
		}
	}

	sorts := make(bson.D, 0)
	for _, sort := range filterOptions.GetSorts() {
		if sort == nil {
			continue
		}

		if sort.GetField() != nil {
			field := *sort.GetField()

			// default asc
			o := sortOrders[asc]
			if sort.GetOrder() != nil {
				order := *sort.GetOrder()

				_, ok := sortOrders[order]
				if ok {
					o = sortOrders[order]
				}
			}

			sorts = append(sorts, bson.E{Key: field, Value: o})
		}
	}
	if len(sorts) != 0 {
		opts.SetSort(sorts)
	}

	return opts
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
		fn := val.Type().Field(i).Tag.Get(tagFilter) // filter name
		fv := reflect.Indirect(val.Field(i))         // filter value
		fk := fv.Kind()                              // filter type

		if fn == ignore {
			continue
		}

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
