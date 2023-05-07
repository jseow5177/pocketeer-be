package model

import (
	"github.com/jseow5177/pockteer-be/data/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	CatID      primitive.ObjectID `bson:"_id,omitempty"`
	CatName    *string            `bson:"cat_name,omitempty"`
	CatType    *uint32            `bson:"cat_type,omitempty"`
	CreateTime *uint64            `bson:"create_time,omitempty"`
	UpdateTime *uint64            `bson:"update_time,omitempty"`
}

func ToCategoryModel(c *entity.Category) *Category {
	objID := primitive.NewObjectID()

	if primitive.IsValidObjectID(c.GetCatID()) {
		objID, _ = primitive.ObjectIDFromHex(c.GetCatID())
	}

	return &Category{
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
		CatName:    c.CatName,
		CatType:    c.CatType,
		CreateTime: c.CreateTime,
		UpdateTime: c.UpdateTime,
	}
}

func (m *Category) GetCatID() string {
	if m != nil {
		return m.CatID.Hex()
	}
	return ""
}

func (m *Category) GetCatName() string {
	if m != nil && m.CatName != nil {
		return *m.CatName
	}
	return ""
}

func (m *Category) GetCatType() uint32 {
	if m != nil && m.CatType != nil {
		return *m.CatType
	}
	return 0
}

func (m *Category) GetCreateTime() uint64 {
	if m != nil && m.CreateTime != nil {
		return *m.CreateTime
	}
	return 0
}

func (m *Category) GetUpdateTime() uint64 {
	if m != nil && m.UpdateTime != nil {
		return *m.UpdateTime
	}
	return 0
}
