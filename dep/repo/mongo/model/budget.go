package model

import (
	"github.com/jseow5177/pockteer-be/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Budget struct {
	BudgetID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID       *string            `bson:"user_id,omitempty"`
	CategoryID   *string            `bson:"category_id,omitempty"`
	BudgetType   *uint32            `bson:"budget_type,omitempty"`
	Amount       *float64           `bson:"amount,omitempty"`
	BudgetStatus *uint32            `bson:"budget_status,omitempty"`
	CreateTime   *uint64            `bson:"create_time,omitempty"`
	UpdateTime   *uint64            `bson:"update_time,omitempty"`
}

func ToBudgetModelFromEntity(b *entity.Budget) *Budget {
	if b == nil {
		return nil
	}

	objID := primitive.NilObjectID
	if primitive.IsValidObjectID(b.GetBudgetID()) {
		objID, _ = primitive.ObjectIDFromHex(b.GetBudgetID())
	}

	return &Budget{
		BudgetID:     objID,
		UserID:       b.UserID,
		CategoryID:   b.CategoryID,
		BudgetType:   b.BudgetType,
		BudgetStatus: b.BudgetStatus,
		Amount:       b.Amount,
		CreateTime:   b.CreateTime,
		UpdateTime:   b.UpdateTime,
	}
}

func (b *Budget) GetBudgetID() string {
	if b != nil {
		return b.BudgetID.Hex()
	}
	return ""
}

func (b *Budget) GetUserID() string {
	if b != nil && b.UserID != nil {
		return *b.UserID
	}
	return ""
}

func (b *Budget) GetCategoryID() string {
	if b != nil && b.CategoryID != nil {
		return *b.CategoryID
	}
	return ""
}

func (b *Budget) GetBudgetType() uint32 {
	if b != nil && b.BudgetType != nil {
		return *b.BudgetType
	}
	return 0
}

func (b *Budget) GetBudgetStatus() uint32 {
	if b != nil && b.BudgetStatus != nil {
		return *b.BudgetStatus
	}
	return 0
}

func (b *Budget) GetAmount() float64 {
	if b != nil && b.Amount != nil {
		return *b.Amount
	}
	return 0
}

func (b *Budget) GetCreateTime() uint64 {
	if b != nil && b.CreateTime != nil {
		return *b.CreateTime
	}
	return 0
}

func (b *Budget) GetUpdateTime() uint64 {
	if b != nil && b.UpdateTime != nil {
		return *b.UpdateTime
	}
	return 0
}
