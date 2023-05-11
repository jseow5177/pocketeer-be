package model

import (
	"github.com/jseow5177/pockteer-be/data/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	UserID     *string            `bson:"user_id,omitempty"`
	CatID      primitive.ObjectID `bson:"_id,omitempty"`
	CatName    *string            `bson:"cat_name,omitempty"`
	CatType    *uint32            `bson:"cat_type,omitempty"`
	CreateTime *uint64            `bson:"create_time,omitempty"`
	UpdateTime *uint64            `bson:"update_time,omitempty"`
}

func ToCategoryModel(c *entity.Category) *Category {
	objID := primitive.NilObjectID
	if primitive.IsValidObjectID(c.GetCatID()) {
		objID, _ = primitive.ObjectIDFromHex(c.GetCatID())
	}

	return &Category{
		UserID:     c.UserID,
		CatID:      objID,
		CatName:    c.CatName,
		CatType:    c.CatType,
		CreateTime: c.CreateTime,
		UpdateTime: c.UpdateTime,
	}
}

func ToCategoryEntity(c *Category) *entity.Category {
	return &entity.Category{
		CatID:      goutil.String(c.CatID.Hex()),
		UserID:     c.UserID,
		CatName:    c.CatName,
		CatType:    c.CatType,
		CreateTime: c.CreateTime,
		UpdateTime: c.UpdateTime,
	}
}

func (c *Category) GetCatID() string {
	if c != nil {
		return c.CatID.Hex()
	}
	return ""
}

func (c *Category) GetCatName() string {
	if c != nil && c.CatName != nil {
		return *c.CatName
	}
	return ""
}

func (c *Category) GetCatType() uint32 {
	if c != nil && c.CatType != nil {
		return *c.CatType
	}
	return 0
}

func (c *Category) GetCreateTime() uint64 {
	if c != nil && c.CreateTime != nil {
		return *c.CreateTime
	}
	return 0
}

func (c *Category) GetUpdateTime() uint64 {
	if c != nil && c.UpdateTime != nil {
		return *c.UpdateTime
	}
	return 0
}
