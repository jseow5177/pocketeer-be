package mongo

import (
	"context"
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
	client, err := mongo.Connect(
		ctx,
		options.Client().
			ApplyURI(cfg.String()).
			SetServerAPIOptions(
				options.ServerAPI(options.ServerAPIVersion1),
			),
	)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return &Mongo{
		client: client,
		db:     client.Database(cfg.Database),
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
	for _, r := range res.InsertedIDs {
		ids = append(ids, r.(primitive.ObjectID).Hex())
	}

	return ids, nil
}

func (mc *MongoColl) update(ctx context.Context, filter bson.D, update interface{}, opts ...*options.UpdateOptions) error {
	_, err := mc.coll.UpdateOne(ctx, filter, mongoutil.BuildUpdate(update), opts...)
	if err != nil {
		return err
	}

	return nil
}

func (mc *MongoColl) updateMany(ctx context.Context, filter bson.D, update interface{}, opts ...*options.UpdateOptions) error {
	_, err := mc.coll.UpdateMany(ctx, filter, mongoutil.BuildUpdate(update), opts...)
	if err != nil {
		return err
	}

	return nil
}

func (mc *MongoColl) deleteMany(ctx context.Context, filter interface{}) error {
	f := mongoutil.BuildFilter(filter)

	_, err := mc.coll.DeleteMany(ctx, f)
	if err != nil {
		return err
	}

	return nil
}

func (mc *MongoColl) get(ctx context.Context, model interface{}, filter bson.D) error {
	if err := mc.coll.FindOne(ctx, filter).Decode(model); err != nil {
		return err
	}

	return nil
}

func (mc *MongoColl) getMany(ctx context.Context, model interface{}, filterOpts filter.FilterOptions, filter bson.D) ([]interface{}, error) {
	opts := mongoutil.BuildFilterOptions(filterOpts)

	cursor, err := mc.coll.Find(ctx, filter, opts)
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
