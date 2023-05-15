package model

import (
	"github.com/jseow5177/pockteer-be/data/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	UserID       *string            `bson:"user_id,omitempty"`
	CategoryID   primitive.ObjectID `bson:"_id,omitempty"`
	CategoryName *string            `bson:"category_name,omitempty"`
	CategoryType *uint32            `bson:"category_type,omitempty"`
	CreateTime   *uint64            `bson:"create_time,omitempty"`
	UpdateTime   *uint64            `bson:"update_time,omitempty"`
}

func ToCategoryModel(c *entity.Category) *Category {
	objID := primitive.NilObjectID
	if primitive.IsValidObjectID(c.GetCategoryID()) {
		objID, _ = primitive.ObjectIDFromHex(c.GetCategoryID())
	}

	return &Category{
		UserID:       c.UserID,
		CategoryID:   objID,
		CategoryName: c.CategoryName,
		CategoryType: c.CategoryType,
		CreateTime:   c.CreateTime,
		UpdateTime:   c.UpdateTime,
	}
}

func ToCategoryEntity(c *Category) *entity.Category {
	return &entity.Category{
		CategoryID:   goutil.String(c.CategoryID.Hex()),
		UserID:       c.UserID,
		CategoryName: c.CategoryName,
		CategoryType: c.CategoryType,
		CreateTime:   c.CreateTime,
		UpdateTime:   c.UpdateTime,
	}
}

func (c *Category) GetCategoryID() string {
	if c != nil {
		return c.CategoryID.Hex()
	}
	return ""
}

func (c *Category) GetCategoryName() string {
	if c != nil && c.CategoryName != nil {
		return *c.CategoryName
	}
	return ""
}

func (c *Category) GetCategoryType() uint32 {
	if c != nil && c.CategoryType != nil {
		return *c.CategoryType
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
