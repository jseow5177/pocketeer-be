package mongoutil

import (
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const maxAggrDepth = 2

type aggrOp string

const (
	AggrSum      aggrOp = "sum"
	AggrMultiply aggrOp = "multiply"
)

type AggrOpt struct {
	Field interface{} // field(s) to aggr on, or
	Aggr  *Aggr       // aggr on a sub-aggr
}

type Aggr struct {
	name    string   // aggr name
	op      aggrOp   // aggr op
	aggrOpt *AggrOpt // options to aggr on
}

func (ag *Aggr) GetName() string {
	return ag.name
}

func (ag *Aggr) buildPipe() bson.E {
	return bson.E{
		Key:   ag.name,
		Value: ag.buildVal(),
	}
}

func (ag *Aggr) buildVal() bson.D {
	be := bson.E{
		Key: Prefix(fmt.Sprint(ag.op)),
	}

	if ag.aggrOpt == nil {
		return nil
	}

	if ag.aggrOpt.Aggr != nil {
		be.Value = ag.aggrOpt.Aggr.buildVal()
	} else {
		fv := reflect.Indirect(reflect.ValueOf(ag.aggrOpt.Field))
		fk := fv.Kind()

		if fk == reflect.Invalid {
			be.Value = nil
		}

		if fk == reflect.Slice {
			ba := make(bson.A, 0)
			for i := 0; i < fv.Len(); i++ {
				elem := fv.Index(i)
				ba = append(ba, Prefix(fmt.Sprint(elem)))
			}
			be.Value = ba
		} else {
			be.Value = Prefix(fmt.Sprint(ag.aggrOpt.Field))
		}
	}

	return bson.D{be}
}

func (ag *Aggr) isTooDeep(depth int, maxDepth int) bool {
	if depth > maxDepth {
		return true
	}

	if ag.aggrOpt.Aggr != nil {
		return ag.aggrOpt.Aggr.isTooDeep(depth+1, maxDepth)
	}

	return false
}

func NewAggr(name string, op aggrOp, opt *AggrOpt) *Aggr {
	return &Aggr{
		name:    name,
		op:      op,
		aggrOpt: opt,
	}
}

func BuildAggrPipeline(filter bson.D, groupBy string, aggrs ...*Aggr) primitive.A {
	pipeline := make(bson.A, 0)

	if filter != nil {
		pipeline = append(pipeline, bson.D{{Key: Prefix("match"), Value: filter}})
	}

	be := bson.E{
		Key: Prefix("group"),
	}
	group := make(bson.D, 0)

	var val interface{}
	if groupBy != "" {
		val = Prefix(groupBy)
	}
	group = append(group, bson.E{Key: "_id", Value: val})

	for _, aggr := range aggrs {
		if aggr.isTooDeep(0, maxAggrDepth) {
			continue
		}
		group = append(group, aggr.buildPipe())
	}

	be.Value = group

	pipeline = append(pipeline, bson.D{be})

	return pipeline
}

func ToFloat64(i interface{}) float64 {
	var f float64
	if v, ok := i.(int32); ok {
		f = float64(v)
	}

	if v, ok := i.(float64); ok {
		f = v
	}

	return f
}
