package mongoutil

import (
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const maxAggrDepth = 2

type Aggr struct {
	name  string      // aggr name
	op    string      // aggr op
	field interface{} // fields to aggr on
	aggr  *Aggr
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
		Key: Prefix(ag.op),
	}

	if ag.aggr != nil {
		be.Value = ag.buildVal()
	} else {
		fv := reflect.Indirect(reflect.ValueOf(ag.field))
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
			be.Value = Prefix(fmt.Sprint(ag.field))
		}
	}

	return bson.D{be}
}

func (ag *Aggr) getDepth(maxDepth int) int {
	if ag == nil || maxDepth == 0 {
		return 0
	}
	return 1 + ag.getDepth(maxDepth-1)
}

func NewAggr(name, op string, field interface{}, subAggr *Aggr) *Aggr {
	return &Aggr{
		name:  name,
		op:    op,
		field: field,
		aggr:  subAggr,
	}
}

func BuildAggrPipeline(filter interface{}, groupBy string, aggrs ...*Aggr) primitive.A {
	pipeline := make(bson.A, 0)

	if filter != nil {
		f := BuildFilter(filter)
		pipeline = append(pipeline, bson.D{{Key: Prefix("match"), Value: f}})
	}

	be := bson.E{
		Key: Prefix("group"),
	}
	group := make(bson.D, 0)

	if groupBy != "" {
		group = append(group, bson.E{Key: "_id", Value: Prefix(groupBy)})
	}

	for _, aggr := range aggrs {
		if aggr.getDepth(maxAggrDepth) == 0 {
			continue
		}
		group = append(group, aggr.buildPipe())
	}

	be.Value = group

	pipeline = append(pipeline, bson.D{be})

	return pipeline
}
