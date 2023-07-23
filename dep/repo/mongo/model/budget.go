package model

import (
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Budget struct {
	UserID           *string            `bson:"user_id,omitempty"`
	BudgetID         primitive.ObjectID `bson:"_id,omitempty"`
	BudgetType       *uint32            `bson:"budget_type,omitempty"`
	BudgetName       *string            `bson:"budget_name,omitempty"`
	CategoryIDs      []string           `bson:"category_ids,omitempty"`
	Status           *uint32            `bson:"status,omitempty"`
	BudgetBreakdowns []*BudgetBreakdown `bson:"budget_breakdowns,omitempty"`
	UpdateTime       *uint64            `bson:"update_time,omitempty"`
}

type BudgetBreakdown struct {
	Year   *int     `bson:"year,omitempty"`
	Month  *int     `bson:"month,omitempty"`
	Amount *float64 `bson:"amount,omitempty"`
}

func ToBudgetEntity(b *Budget) *entity.Budget {
	return &entity.Budget{
		BudgetID:     goutil.String(b.GetBudgetID()),
		UserID:       b.UserID,
		BudgetType:   b.BudgetType,
		BudgetName:   b.BudgetName,
		CategoryIDs:  b.CategoryIDs,
		Status:       b.Status,
		BreakdownMap: ToBudgetBreakdownMap(b.BudgetBreakdowns),
		UpdateTime:   b.UpdateTime,
	}
}

func ToBudgetBreakdownMap(breakdowns []*BudgetBreakdown) map[entity.DateInfo]*entity.BudgetBreakdown {
	_map := make(map[entity.DateInfo]*entity.BudgetBreakdown)

	for _, bd := range breakdowns {
		dateInfo := entity.DateInfo{
			Year:  bd.GetYear(),
			Month: bd.GetMonth(),
		}

		_map[dateInfo] = &entity.BudgetBreakdown{
			Year:   bd.Year,
			Month:  bd.Month,
			Amount: bd.Amount,
		}
	}

	return _map
}

func ToBudgetModel(e *entity.Budget) *Budget {
	objID := primitive.NewObjectID()
	if primitive.IsValidObjectID(e.GetBudgetID()) {
		objID, _ = primitive.ObjectIDFromHex(e.GetBudgetID())
	}

	b := &Budget{
		BudgetID:         objID,
		UserID:           e.UserID,
		BudgetType:       e.BudgetType,
		BudgetName:       e.BudgetName,
		CategoryIDs:      e.CategoryIDs,
		Status:           e.Status,
		BudgetBreakdowns: ToModelBreakdowns(e.BreakdownMap),
		UpdateTime:       e.UpdateTime,
	}

	return b
}

func ToModelBreakdowns(breakdownMap map[entity.DateInfo]*entity.BudgetBreakdown) []*BudgetBreakdown {
	budgetBreakdowns := make([]*BudgetBreakdown, 0)
	for _, bd := range breakdownMap {
		budgetBreakdowns = append(
			budgetBreakdowns,
			&BudgetBreakdown{
				Year:   bd.Year,
				Month:  bd.Month,
				Amount: bd.Amount,
			},
		)
	}

	return budgetBreakdowns
}

// Getters ----------------------------------------------------------------------
func (b *Budget) GetBudgetID() string {
	if b != nil {
		return b.BudgetID.Hex()
	}
	return ""
}

func (b *Budget) GetUserID() string {
	if b != nil {
		return *b.UserID
	}
	return ""
}

func (b *Budget) GetBudgetType() uint32 {
	if b != nil {
		return *b.BudgetType
	}
	return 0
}

func (b *Budget) GetBudgetName() string {
	if b != nil {
		return *b.BudgetName
	}
	return ""
}

func (b *Budget) GetCategoryIDs() []string {
	if b != nil {
		return b.CategoryIDs
	}
	return []string{}
}

func (b *Budget) GetBudgetBreakdowns() []*BudgetBreakdown {
	if b != nil {
		return b.BudgetBreakdowns
	}
	return []*BudgetBreakdown{}
}

func (bd *BudgetBreakdown) GetYear() int {
	if bd != nil {
		return *bd.Year
	}
	return 0
}

func (bd *BudgetBreakdown) GetMonth() int {
	if bd != nil {
		return *bd.Month
	}
	return 0
}

func (bd *BudgetBreakdown) GetAmount() float64 {
	if bd != nil {
		return *bd.Amount
	}
	return 0
}
