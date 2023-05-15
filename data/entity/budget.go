package entity

type Budget struct {
	CatID           *string
	CatName         *string
	BudgetType      *uint32
	TransactionType *uint32
	Year            *uint32
	Month           *uint32
	BudgetAmount    *int64
	Used            *int64
}

func (e *Budget) GetCatID() string {
	if e != nil && e.CatID != nil {
		return *e.CatID
	}
	return ""
}

func (e *Budget) GetCatName() string {
	if e != nil && e.CatName != nil {
		return *e.CatName
	}
	return ""
}

func (e *Budget) GetBudgetType() uint32 {
	if e != nil && e.BudgetType != nil {
		return *e.BudgetType
	}
	return 0
}

func (e *Budget) GetTransactionType() uint32 {
	if e != nil && e.TransactionType != nil {
		return *e.TransactionType
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

func (e *Budget) GetUsed() int64 {
	if e != nil && e.Used != nil {
		return *e.Used
	}
	return 0
}
