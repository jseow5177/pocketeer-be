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

var sortOrders = map[string]int{
	"asc":  1,
	"desc": -1,
}

var supportedOps = map[string]string{
	"eq":         "equal",
	"ne":         "not equal",
	"gt":         "greater than",
	"gte":        "greater than equal",
	"lt":         "less than",
	"lte":        "less than equal",
	"in":         "in",
	"nin":        "not in",
	"bitsAllSet": "bit all set",
}

func GetOp(op string) string {
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
			o := sortOrders["asc"]
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

func BuildFilter(filter interface{}) bson.D {
	if filter == nil {
		return nil
	}

	val := reflect.ValueOf(filter)
	if val.IsNil() {
		return nil
	}
	val = val.Elem()

	conds := make(bson.A, 0)
	for i := 0; i < val.NumField(); i++ {
		ft := val.Type().Field(i).Tag.Get("filter") // filter tag
		fv := reflect.Indirect(val.Field(i))        // filter value
		fk := fv.Kind()                             // filter type

		if ft == "-" {
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
		parts := strings.SplitN(ft, "__", 2)

		// operator
		var op string
		if len(parts) > 1 {
			op = parts[1]
			if _, ok := supportedOps[op]; !ok {
				continue
			}
		} else {
			// default to eq
			op = "eq"
		}
		op = GetOp(op)

		// filter value
		v := fv.Interface()

		// handle _id field
		fn := strcase.ToSnake(parts[0])
		if fn == "_id" {
			if fk == reflect.Slice {
				ids := make([]primitive.ObjectID, 0)
				// loop through each element and convert them to ObjectIDs
				for i := 0; i < fv.Len(); i++ {
					e := fv.Index(i).Interface()
					s := fmt.Sprint(e)
					if primitive.IsValidObjectID(s) {
						id, _ := primitive.ObjectIDFromHex(s)
						ids = append(ids, id)
					}
				}
				v = reflect.ValueOf(ids).Interface()
			} else {
				s := fmt.Sprint(v)
				if primitive.IsValidObjectID(s) {
					v, _ = primitive.ObjectIDFromHex(s)
				}
			}
		}

		cond := bson.D{{Key: fn, Value: bson.D{{Key: op, Value: v}}}}
		conds = append(conds, cond)
	}

	if len(conds) == 0 {
		return nil
	}

	return bson.D{{Key: GetOp("and"), Value: conds}}
}
