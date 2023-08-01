package entity

import (
	"errors"
	"time"

	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/util"
)

var (
	ErrBudgetNotAllowed      = errors.New("budget not allowed under category")
	ErrUnsupportedBudgetType = errors.New("unsupported budget type")
	ErrBudgetDateEmpty       = errors.New("date cannot be empty")
)

type BudgetRepeat uint32

const (
	BudgetRepeatNow BudgetRepeat = iota
	BudgetRepeatNowToFuture
	BudgetRepeatAllTime
)

var BudgetRepeats = map[uint32]string{
	uint32(BudgetRepeatNow):         "for a period only",
	uint32(BudgetRepeatNowToFuture): "for a period and all future periods",
	uint32(BudgetRepeatAllTime):     "for all periods past and future",
}

type BudgetType uint32

const (
	BudgetTypeInvalid BudgetType = iota
	BudgetTypeMonth
	BudgetTypeYear
)

var BudgetTypes = map[uint32]string{
	uint32(BudgetTypeMonth): "month",
	uint32(BudgetTypeYear):  "year",
}

type BudgetStatus uint32

const (
	BudgetStatusInvalid BudgetStatus = iota
	BudgetStatusNormal
	BudgetStatusDeleted
)

type IBudget interface {
	CheckOpts() error
}

type MonthBudget struct {
	Budget *Budget
}

func (b *MonthBudget) CheckOpts() (err error) {
	return b.Budget.checkDate(util.GetMonthRangeAsDate)
}

type YearBudget struct {
	Budget *Budget
}

func (b *YearBudget) CheckOpts() (err error) {
	return b.Budget.checkDate(util.GetYearRangeAsDate)
}

type Budget struct {
	UserID       *string
	BudgetID     *string
	CategoryID   *string
	Amount       *float64
	BudgetType   *uint32
	BudgetStatus *uint32
	StartDate    *uint64
	EndDate      *uint64
	CreateTime   *uint64
	UpdateTime   *uint64

	BudgetDate   *string
	BudgetRepeat *uint32
	UsedAmount   *float64
}

type BudgetOption = func(b *Budget)

func WithBudgetID(budgetID *string) BudgetOption {
	return func(b *Budget) {
		b.BudgetID = budgetID
	}
}

func WithBudgetRepeat(budgetPeriod *uint32) BudgetOption {
	return func(b *Budget) {
		b.BudgetRepeat = budgetPeriod
	}
}

func WithBudgetType(budgetType *uint32) BudgetOption {
	return func(b *Budget) {
		b.BudgetType = budgetType
	}
}

func WithBudgetStatus(budgetStatus *uint32) BudgetOption {
	return func(b *Budget) {
		b.BudgetStatus = budgetStatus
	}
}

func WithBudgetAmount(amount *float64) BudgetOption {
	return func(b *Budget) {
		b.Amount = amount
	}
}

func WithBudgetStartDate(startDate *uint64) BudgetOption {
	return func(b *Budget) {
		b.StartDate = startDate
	}
}

func WithBudgetEndDate(endDate *uint64) BudgetOption {
	return func(b *Budget) {
		b.EndDate = endDate
	}
}

func WithBudgetDate(budgetTime *string) BudgetOption {
	return func(b *Budget) {
		b.BudgetDate = budgetTime
	}
}

func WithBudgetCreateTime(createTime *uint64) BudgetOption {
	return func(b *Budget) {
		b.CreateTime = createTime
	}
}

func WithBudgetUpdateTime(updateTime *uint64) BudgetOption {
	return func(b *Budget) {
		b.UpdateTime = updateTime
	}
}

func NewBudget(userID, categoryID string, opts ...BudgetOption) (*Budget, error) {
	now := uint64(time.Now().UnixMilli())
	b := &Budget{
		UserID:       goutil.String(userID),
		CategoryID:   goutil.String(categoryID),
		BudgetType:   goutil.Uint32(uint32(BudgetTypeMonth)),
		BudgetStatus: goutil.Uint32(uint32(BudgetStatusNormal)),
		BudgetRepeat: goutil.Uint32(uint32(BudgetRepeatNow)),
		Amount:       goutil.Float64(0),
		CreateTime:   goutil.Uint64(now),
		UpdateTime:   goutil.Uint64(now),
	}
	for _, opt := range opts {
		opt(b)
	}

	if err := b.checkOpts(); err != nil {
		return nil, err
	}

	return b, nil
}

func (b *Budget) checkOpts() error {
	ib, err := b.ToIBudget()
	if err != nil {
		return err
	}
	return ib.CheckOpts()
}

type getDateRangeFn func(date, timezone string) (startDate, endDate uint64, err error)

func (b *Budget) checkDate(getDateRange getDateRangeFn) error {
	// no-op if already set
	if b.StartDate != nil && b.EndDate != nil {
		return nil
	}

	if b.MustSetBudgetDate() && b.GetBudgetDate() == "" {
		return ErrBudgetDateEmpty
	}

	var (
		startDate uint64
		endDate   uint64
		err       error
	)
	switch b.GetBudgetRepeat() {
	case uint32(BudgetRepeatNow):
		startDate, endDate, err = getDateRange(b.GetBudgetDate(), "")
	case uint32(BudgetRepeatNowToFuture):
		startDate, _, err = getDateRange(b.GetBudgetDate(), "")
	case uint32(BudgetRepeatAllTime):
	}
	if err != nil {
		return err
	}

	b.StartDate = goutil.Uint64(startDate)
	b.EndDate = goutil.Uint64(endDate)

	return nil
}

func (b *Budget) ToIBudget() (IBudget, error) {
	switch b.GetBudgetType() {
	case uint32(BudgetTypeMonth):
		return &MonthBudget{b}, nil
	case uint32(BudgetTypeYear):
		return &YearBudget{b}, nil
	}
	return nil, ErrUnsupportedBudgetType
}

func (b *Budget) GetBudgetID() string {
	if b != nil && b.BudgetID != nil {
		return *b.BudgetID
	}
	return ""
}

func (b *Budget) SetBudgetID(budgetID *string) {
	b.BudgetID = budgetID
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

func (b *Budget) GetBudgetRepeat() uint32 {
	if b != nil && b.BudgetRepeat != nil {
		return *b.BudgetRepeat
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

func (b *Budget) GetUsedAmount() float64 {
	if b != nil && b.UsedAmount != nil {
		return *b.UsedAmount
	}
	return 0
}

func (b *Budget) SetUsedAmount(usedAmount *float64) {
	b.UsedAmount = usedAmount
}

func (b *Budget) GetBudgetDate() string {
	if b != nil && b.BudgetDate != nil {
		return *b.BudgetDate
	}
	return ""
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

func (b *Budget) CanBudgetUnderCategory(c *Category) (bool, error) {
	if !c.CanAddBudget() {
		return false, ErrBudgetNotAllowed
	}
	return true, nil
}

func (b *Budget) MustSetBudgetDate() bool {
	return b.GetBudgetRepeat() == uint32(BudgetRepeatNow) ||
		b.GetBudgetRepeat() == uint32(BudgetRepeatNowToFuture)
}

func (b *Budget) IsDeleted() bool {
	return b.GetBudgetStatus() == uint32(BudgetStatusDeleted)
}

func (b *Budget) IsMonth() bool {
	return b.GetBudgetType() == uint32(BudgetTypeMonth)
}

func (b *Budget) IsYear() bool {
	return b.GetBudgetType() == uint32(BudgetTypeYear)
}
