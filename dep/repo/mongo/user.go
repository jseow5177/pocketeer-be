package mongo

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/dep/repo/mongo/model"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/pkg/mongoutil"
	"go.mongodb.org/mongo-driver/mongo"
)

const userCollName = "user"

type userMongo struct {
	mColl *MongoColl
}

func NewUserMongo(mongo *Mongo) repo.UserRepo {
	return &userMongo{
		mColl: NewMongoColl(mongo, userCollName),
	}
}

func (m *userMongo) Create(ctx context.Context, u *entity.User) (string, error) {
	um := model.ToUserModel(u)
	id, err := m.mColl.create(ctx, um)
	if err != nil {
		return "", err
	}
	u.SetUserID(goutil.String(id))

	return id, nil
}

func (m *userMongo) Get(ctx context.Context, uf *repo.UserFilter) (*entity.User, error) {
	f := mongoutil.BuildFilter(uf)

	u := new(model.User)
	if err := m.mColl.get(ctx, &u, f); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repo.ErrUserNotFound
		}
		return nil, err
	}

	return model.ToUserEntity(u)
}
