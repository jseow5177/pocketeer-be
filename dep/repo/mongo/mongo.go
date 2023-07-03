package mongo

import (
	"context"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/pkg/filter"
	"github.com/jseow5177/pockteer-be/pkg/mongoutil"
)

type Mongo struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewMongo(ctx context.Context, cfg *config.Mongo) (*Mongo, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.String()))
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return &Mongo{
		client: client,
		db:     client.Database(cfg.DBName),
	}, nil
}

func (m *Mongo) Close(ctx context.Context) error {
	defer func() {
		m.db = nil
		m.client = nil
	}()
	return m.client.Disconnect(ctx)
}

func (m *Mongo) WithTx(ctx context.Context, txFn func(txCtx context.Context) error) error {
	session, err := m.client.StartSession()
	if err != nil {
		return err
	}

	if err = session.StartTransaction(); err != nil {
		return err
	}

	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err = txFn(sc); err != nil {
			return err
		}

		if err = session.CommitTransaction(sc); err != nil {
			return err
		}

		return nil
	}); err != nil {
		_ = session.AbortTransaction(ctx)
		return err
	}

	session.EndSession(ctx)

	return nil
}

type MongoColl struct {
	coll *mongo.Collection
}

func NewMongoColl(m *Mongo, collName string) *MongoColl {
	return &MongoColl{
		coll: m.db.Collection(collName),
	}
}

func (mc *MongoColl) create(ctx context.Context, doc interface{}) (string, error) {
	res, err := mc.coll.InsertOne(ctx, doc)
	if err != nil {
		return "", err
	}

	id := res.InsertedID.(primitive.ObjectID)

	return id.Hex(), nil
}

func (mc *MongoColl) createMany(ctx context.Context, docs []interface{}) ([]string, error) {
	res, err := mc.coll.InsertMany(ctx, docs)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0)
	for _, id := range res.InsertedIDs {
		ids = append(ids, id.(primitive.ObjectID).Hex())
	}

	return ids, nil
}

func (mc *MongoColl) update(ctx context.Context, filter, update interface{}) error {
	var (
		f = mongoutil.BuildFilter(filter)
		u = mongoutil.BuildUpdate(update)
	)

	_, err := mc.coll.UpdateOne(ctx, f, u)
	if err != nil {
		return err
	}

	return nil
}

func (mc *MongoColl) upsert(ctx context.Context, uniqueKey string, doc interface{}) (id string, err error) {
	keyValue := getUniqueKeyValue(doc, uniqueKey)

	filter := bson.M{uniqueKey: keyValue}
	update := bson.M{"$set": removeUniqueKeyField(doc, uniqueKey)}

	upsert := mongo.NewUpdateOneModel()
	upsert.SetFilter(filter)
	upsert.SetUpdate(update)
	upsert.SetUpsert(true)

	opts := options.Update().SetUpsert(true)
	result, err := mc.coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return "", err
	}

	if result.UpsertedID != nil {
		objID := result.UpsertedID.(primitive.ObjectID)
		id = objID.Hex()
	}

	return id, nil
}

func (mc *MongoColl) get(ctx context.Context, filter interface{}, model interface{}) error {
	f := mongoutil.BuildFilter(filter)

	if err := mc.coll.FindOne(ctx, f).Decode(model); err != nil {
		return err
	}

	return nil
}

func (mc *MongoColl) getMany(ctx context.Context, filter interface{}, filterOpts filter.FilterOptions, model interface{}) ([]interface{}, error) {
	f := mongoutil.BuildFilter(filter)
	opts := mongoutil.BuildFilterOptions(filterOpts)

	cursor, err := mc.coll.Find(ctx, f, opts)
	if err != nil {
		return nil, err
	}

	t := reflect.TypeOf(model).Elem()

	res := make([]interface{}, 0)
	for cursor.Next(ctx) {
		m := reflect.New(t).Interface()
		if err = cursor.Decode(m); err != nil {
			return nil, err
		}
		res = append(res, m)
	}

	return res, nil
}

func (mc *MongoColl) sum(ctx context.Context, sumBy, field string, filter interface{}) (map[string]float64, error) {
	f := mongoutil.BuildFilter(filter)

	// aggregation pipeline
	pipeline := make(bson.A, 0)
	if f != nil {
		pipeline = append(pipeline, bson.D{{Key: mongoutil.GetOp("match"), Value: f}})
	}

	pipeline = append(pipeline, bson.D{{Key: mongoutil.GetOp("group"), Value: bson.D{
		{Key: "_id", Value: mongoutil.GetOp(sumBy)},
		{Key: "sum", Value: bson.D{{Key: mongoutil.GetOp("sum"), Value: mongoutil.GetOp(field)}}},
	}}})

	cursor, err := mc.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	results := make([]bson.M, 0)
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	sumResults := make(map[string]float64)
	for _, result := range results {
		sumResults[fmt.Sprint(result["_id"])] = result["sum"].(float64)
	}

	return sumResults, nil
}

func getUniqueKeyValue(doc interface{}, uniqueKey string) interface{} {
	value := bson.M{}
	bsonBytes, _ := bson.Marshal(doc)
	_ = bson.Unmarshal(bsonBytes, &value)

	return value[uniqueKey]
}

func removeUniqueKeyField(doc interface{}, uniqueKey string) interface{} {
	value := bson.M{}
	bsonBytes, _ := bson.Marshal(doc)
	_ = bson.Unmarshal(bsonBytes, &value)

	delete(value, uniqueKey)

	return value
}
