package entity

import (
	"errors"
	"math"
	"time"

	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/util"
)

var (
	ErrBudgetNotAllowed = errors.New("budget not allowed under category")
	ErrBudgetDateEmpty  = errors.New("date cannot be empty")
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

type getDateRangeFn func(date, timezone string) (startDate, endDate uint64, err error)

var DateRangeFuncs = map[uint32]getDateRangeFn{
	uint32(BudgetTypeMonth): util.GetMonthRangeAsDate,
	uint32(BudgetTypeYear):  util.GetYearRangeAsDate,
}

type BudgetStatus uint32

const (
	BudgetStatusInvalid BudgetStatus = iota
	BudgetStatusNormal
	BudgetStatusDeleted
)

type BudgetUpdate struct {
	BudgetType *uint32
	Amount     *float64
	StartDate  *uint64
	EndDate    *uint64
	UpdateTime *uint64
}

func (bu *BudgetUpdate) GetBudgetType() uint32 {
	if bu != nil && bu.BudgetType != nil {
		return *bu.BudgetType
	}
	return 0
}

func (bu *BudgetUpdate) GetAmount() float64 {
	if bu != nil && bu.Amount != nil {
		return *bu.Amount
	}
	return 0
}

func (bu *BudgetUpdate) GetUpdateTime() uint64 {
	if bu != nil && bu.UpdateTime != nil {
		return *bu.UpdateTime
	}
	return 0
}

func (bu *BudgetUpdate) GetStartDate() uint64 {
	if bu != nil && bu.StartDate != nil {
		return *bu.StartDate
	}
	return 0
}

func (bu *BudgetUpdate) GetEndDate() uint64 {
	if bu != nil && bu.EndDate != nil {
		return *bu.EndDate
	}
	return 0
}

type BudgetUpdateOption func(bu *BudgetUpdate)

func WithUpdateBudgetType(budgetType *uint32) BudgetUpdateOption {
	return func(bu *BudgetUpdate) {
		bu.BudgetType = budgetType
	}
}

func WithUpdateBudgetAmount(amount *float64) BudgetUpdateOption {
	return func(bu *BudgetUpdate) {
		bu.Amount = amount
	}
}

func WithUpdateBudgetStartDate(startDate *uint64) BudgetUpdateOption {
	return func(bu *BudgetUpdate) {
		bu.StartDate = startDate
	}
}

func WithUpdateBudgetEndDate(endDate *uint64) BudgetUpdateOption {
	return func(bu *BudgetUpdate) {
		bu.EndDate = endDate
	}
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

	UsedAmount *float64
}

type BudgetOption = func(b *Budget)

func WithBudgetID(budgetID *string) BudgetOption {
	return func(b *Budget) {
		b.BudgetID = budgetID
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

func (b *Budget) Update(bu *BudgetUpdate) (budgetUpdate *BudgetUpdate, hasUpdate bool, err error) {
	budgetUpdate = new(BudgetUpdate)

	if bu.Amount != nil && bu.GetAmount() != b.GetAmount() {
		hasUpdate = true
		b.Amount = bu.Amount

		defer func() {
			budgetUpdate.Amount = b.Amount
		}()
	}

	if bu.BudgetType != nil && bu.GetBudgetType() != b.GetBudgetType() {
		hasUpdate = true
		b.BudgetType = bu.BudgetType

		defer func() {
			budgetUpdate.BudgetType = b.BudgetType
		}()
	}

	if bu.StartDate != nil && bu.GetStartDate() != b.GetStartDate() {
		hasUpdate = true
		b.StartDate = bu.StartDate

		defer func() {
			budgetUpdate.StartDate = b.StartDate
		}()
	}

	if bu.EndDate != nil && bu.GetEndDate() != b.GetEndDate() {
		hasUpdate = true
		b.EndDate = bu.EndDate

		defer func() {
			budgetUpdate.EndDate = b.EndDate
		}()
	}

	if !hasUpdate {
		return
	}

	now := goutil.Uint64(uint64(time.Now().UnixMilli()))
	b.UpdateTime = now

	if err = b.checkOpts(); err != nil {
		return nil, false, err
	}

	budgetUpdate.UpdateTime = now

	return
}

func (b *Budget) checkOpts() error {
	b.Amount = goutil.Float64(math.Abs(b.GetAmount()))

	return nil
}

func GetBudgetDateRange(date string, budgetType, budgetRepeat uint32) (startDate, endDate uint64, err error) {
	fn := DateRangeFuncs[budgetType]
	if fn == nil {
		return 0, 0, ErrInvalidBudgetType
	}

	switch budgetRepeat {
	case uint32(BudgetRepeatNow):
		startDate, endDate, err = fn(date, "")
	case uint32(BudgetRepeatNowToFuture):
		startDate, _, err = fn(date, "")
	case uint32(BudgetRepeatAllTime):
	}
	if err != nil {
		return 0, 0, err
	}

	return
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

func (b *Budget) GetStartDate() uint64 {
	if b != nil && b.StartDate != nil {
		return *b.StartDate
	}
	return 0
}

func (b *Budget) GetEndDate() uint64 {
	if b != nil && b.EndDate != nil {
		return *b.EndDate
	}
	return 0
}

func (b *Budget) CanBudgetUnderCategory(c *Category) (bool, error) {
	if !c.CanAddBudget() {
		return false, ErrBudgetNotAllowed
	}
	return true, nil
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
