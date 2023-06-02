package model

import (
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Budget struct {
	BudgetID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID       *string            `bson:"user_id,omitempty"`
	CategoryID   *string            `bson:"category_id,omitempty"`
	IsDefault    *bool              `bson:"is_default,omitempty"`
	BudgetType   *uint32            `bson:"budget_type,omitempty"`
	Year         *uint32            `bson:"year,omitempty"`
	Month        *uint32            `bson:"month,omitempty"`
	BudgetAmount *int64             `bson:"budget_amount,omitempty"`
}

func ToBudgetModels(bs []*entity.Budget) []*Budget {
	budgets := make([]*Budget, len(bs))

	for idx, b := range bs {
		budgets[idx] = ToBudgetModel(b)
	}

	return budgets
}

func ToBudgetModel(b *entity.Budget) *Budget {
	objID := primitive.NewObjectID()

	if primitive.IsValidObjectID(b.GetBudgetID()) {
		objID, _ = primitive.ObjectIDFromHex(b.GetBudgetID())
	}

	return &Budget{
		BudgetID:     objID,
		UserID:       b.UserID,
		CategoryID:   b.CategoryID,
		IsDefault:    b.IsDefault,
		BudgetType:   b.BudgetType,
		Year:         b.Year,
		Month:        b.Month,
		BudgetAmount: b.BudgetAmount,
	}
}

func ToBudgetEntity(b *Budget) *entity.Budget {
	return &entity.Budget{
		BudgetID:     goutil.String(b.GetBudgetID()),
		UserID:       b.UserID,
		CategoryID:   b.CategoryID,
		IsDefault:    b.IsDefault,
		BudgetType:   b.BudgetType,
		Year:         b.Year,
		Month:        b.Month,
		BudgetAmount: b.BudgetAmount,
	}
}

func (b *Budget) GetBudgetID() string {
	if b != nil {
		b.BudgetID.Hex()
	}
	return ""
}

func (b *Budget) GetUserID() string {
	if b != nil {
		return *b.UserID
	}
	return ""
}

func (b *Budget) GetCategoryID() string {
	if b != nil {
		return *b.CategoryID
	}
	return ""
}

func (b *Budget) GetIsDefault() bool {
	if b != nil {
		return *b.IsDefault
	}
	return false
}

func (b *Budget) GetBudgetType() uint32 {
	if b != nil {
		return *b.BudgetType
	}
	return 0
}

func (b *Budget) GetYear() uint32 {
	if b != nil {
		return *b.Year
	}
	return 0
}

func (b *Budget) GetMonth() uint32 {
	if b != nil {
		return *b.Month
	}
	return 0
}

func (b *Budget) GetBudgetAmount() int64 {
	if b != nil {
		return *b.BudgetAmount
	}
	return 0
}
