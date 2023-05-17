package entity

type Budget struct {
	BudgetID     *string
	UserID       *string
	CategoryID   *string
	IsDefault    *bool
	BudgetType   *uint32
	Year         *uint32
	Month        *uint32
	BudgetAmount *int64
}

func NewEmptyMonthBudget(
	userID,
	categoryID string,
) *Budget {
	return nil
}

func (e *Budget) GetCategoryID() string {
	if e != nil && e.CategoryID != nil {
		return *e.CategoryID
	}
	return ""
}

func (e *Budget) GetBudgetType() uint32 {
	if e != nil && e.BudgetType != nil {
		return *e.BudgetType
	}
	return 0
}

func (e *Budget) GetYear() uint32 {
	if e != nil && e.Year != nil {
		return *e.Year
	}
	return 0
}

func (e *Budget) GetMonth() uint32 {
	if e != nil && e.Month != nil {
		return *e.Month
	}
	return 0
}

func (e *Budget) GetBudgetAmount() int64 {
	if e != nil && e.BudgetAmount != nil {
		return *e.BudgetAmount
	}
	return 0
}
