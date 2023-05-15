package mongo

import (
	"context"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/pkg/errutil"
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

func (mc *MongoColl) wrapError(err error) error {
	if err == mongo.ErrNoDocuments {
		return errutil.ErrNotFound
	}
	return err
}

func (mc *MongoColl) create(ctx context.Context, doc interface{}) (string, error) {
	r, err := mc.coll.InsertOne(ctx, doc)
	if err != nil {
		return "", mc.wrapError(err)
	}

	id := r.InsertedID.(primitive.ObjectID)

	return id.Hex(), nil
}

func (mc *MongoColl) update(ctx context.Context, filter, update interface{}) error {
	var (
		f = mongoutil.BuildFilter(filter)
		u = mongoutil.BuildUpdate(update)
	)

	_, err := mc.coll.UpdateOne(ctx, f, u)
	if err != nil {
		return mc.wrapError(err)
	}

	return nil
}

func (mc *MongoColl) get(ctx context.Context, filter interface{}, model interface{}) error {
	f := mongoutil.BuildFilter(filter)

	if err := mc.coll.FindOne(ctx, f).Decode(model); err != nil {
		return mc.wrapError(err)
	}

	return nil
}

func (mc *MongoColl) getMany(ctx context.Context, filter interface{}, filterOpts filter.FilterOptions, model interface{}) ([]interface{}, error) {
	f := mongoutil.BuildFilter(filter)
	opts := mongoutil.BuildFilterOptions(filterOpts)

	cursor, err := mc.coll.Find(ctx, f, opts)
	if err != nil {
		return nil, mc.wrapError(err)
	}

	t := reflect.TypeOf(model).Elem()

	res := make([]interface{}, 0)
	for cursor.Next(ctx) {
		m := reflect.New(t).Interface()
		if err = cursor.Decode(m); err != nil {
			return nil, mc.wrapError(err)
		}
		res = append(res, m)
	}

	return res, nil
}
