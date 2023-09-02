package entity

import (
	"errors"
	"math"
	"time"

	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/util"
)

var (
	ErrBudgetNotAllowed        = errors.New("budget not allowed under category")
	ErrBudgetDateEmpty         = errors.New("date cannot be empty")
	ErrBudgetConflict          = errors.New("conflict in budget")
	ErrBudgetMustRepeatAllTime = errors.New("budget must repeat all time")
)

type BudgetRepeat uint32

const (
	BudgetRepeatAllTime BudgetRepeat = iota
	BudgetRepeatNowToFuture
	BudgetRepeatNow
)

var BudgetRepeats = map[uint32]string{
	uint32(BudgetRepeatAllTime):     "for all periods past and future",
	uint32(BudgetRepeatNow):         "for a period only",
	uint32(BudgetRepeatNowToFuture): "for a period and all future periods",
}

type BudgetType uint32

const (
	BudgetTypeMonth BudgetType = iota
	BudgetTypeYear
)

var BudgetTypes = map[uint32]string{
	uint32(BudgetTypeMonth): "month",
	uint32(BudgetTypeYear):  "year",
}

type getDateRangeFn func(date, timezone string) (startDate, endDate uint64, err error)

var dateRangeFuncs = map[uint32]getDateRangeFn{
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

func (bu *BudgetUpdate) SetBudgetType(budgetType *uint32) {
	bu.BudgetType = budgetType
}

func (bu *BudgetUpdate) GetAmount() float64 {
	if bu != nil && bu.Amount != nil {
		return *bu.Amount
	}
	return 0
}

func (bu *BudgetUpdate) SetAmount(amount *float64) {
	bu.Amount = amount
}

func (bu *BudgetUpdate) GetUpdateTime() uint64 {
	if bu != nil && bu.UpdateTime != nil {
		return *bu.UpdateTime
	}
	return 0
}

func (bu *BudgetUpdate) SetUpdateTime(updateTime *uint64) {
	bu.UpdateTime = updateTime
}

func (bu *BudgetUpdate) GetStartDate() uint64 {
	if bu != nil && bu.StartDate != nil {
		return *bu.StartDate
	}
	return 0
}

func (bu *BudgetUpdate) SetStartDate(startDate *uint64) {
	bu.StartDate = startDate
}

func (bu *BudgetUpdate) GetEndDate() uint64 {
	if bu != nil && bu.EndDate != nil {
		return *bu.EndDate
	}
	return 0
}

func (bu *BudgetUpdate) SetEndDate(endDate *uint64) {
	bu.EndDate = endDate
}

func NewBudgetUpdate(opts ...BudgetUpdateOption) *BudgetUpdate {
	au := new(BudgetUpdate)
	for _, opt := range opts {
		opt(au)
	}
	return au
}

type BudgetUpdateOption func(bu *BudgetUpdate)

func WithUpdateBudgetType(budgetType *uint32) BudgetUpdateOption {
	return func(bu *BudgetUpdate) {
		bu.SetBudgetType(budgetType)
	}
}

func WithUpdateBudgetAmount(amount *float64) BudgetUpdateOption {
	return func(bu *BudgetUpdate) {
		bu.SetAmount(amount)
	}
}

func WithUpdateBudgetStartDate(startDate *uint64) BudgetUpdateOption {
	return func(bu *BudgetUpdate) {
		bu.SetStartDate(startDate)
	}
}

func WithUpdateBudgetEndDate(endDate *uint64) BudgetUpdateOption {
	return func(bu *BudgetUpdate) {
		bu.SetEndDate(endDate)
	}
}

func WithUpdateBudgetUpdateTime(updateTime *uint64) BudgetUpdateOption {
	return func(bu *BudgetUpdate) {
		bu.SetUpdateTime(updateTime)
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
		b.SetBudgetID(budgetID)
	}
}

func WithBudgetType(budgetType *uint32) BudgetOption {
	return func(b *Budget) {
		b.SetBudgetType(budgetType)
	}
}

func WithBudgetStatus(budgetStatus *uint32) BudgetOption {
	return func(b *Budget) {
		b.SetBudgetStatus(budgetStatus)
	}
}

func WithBudgetAmount(amount *float64) BudgetOption {
	return func(b *Budget) {
		b.SetAmount(amount)
	}
}

func WithBudgetStartDate(startDate *uint64) BudgetOption {
	return func(b *Budget) {
		b.SetStartDate(startDate)
	}
}

func WithBudgetEndDate(endDate *uint64) BudgetOption {
	return func(b *Budget) {
		b.SetEndDate(endDate)
	}
}

func WithBudgetCreateTime(createTime *uint64) BudgetOption {
	return func(b *Budget) {
		b.SetCreateTime(createTime)
	}
}

func WithBudgetUpdateTime(updateTime *uint64) BudgetOption {
	return func(b *Budget) {
		b.SetUpdateTime(updateTime)
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
		StartDate:    goutil.Uint64(0),
		EndDate:      goutil.Uint64(0),
	}
	for _, opt := range opts {
		opt(b)
	}

	if err := b.checkOpts(); err != nil {
		return nil, err
	}

	return b, nil
}

func (b *Budget) Update(bu *BudgetUpdate) (*BudgetUpdate, error) {
	var (
		hasUpdate    bool
		budgetUpdate = new(BudgetUpdate)
	)

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
		return nil, nil
	}

	now := goutil.Uint64(uint64(time.Now().UnixMilli()))
	b.UpdateTime = now

	if err := b.checkOpts(); err != nil {
		return nil, err
	}

	budgetUpdate.UpdateTime = now

	return budgetUpdate, nil
}

func (b *Budget) checkOpts() error {
	b.Amount = goutil.Float64(math.Abs(b.GetAmount()))

	return nil
}

func GetBudgetStartEnd(date string, budgetType, budgetRepeat uint32) (startDate, endDate uint64, err error) {
	fn := dateRangeFuncs[budgetType]
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

func (b *Budget) SetUserID(userID *string) {
	b.UserID = userID
}

func (b *Budget) GetCategoryID() string {
	if b != nil && b.CategoryID != nil {
		return *b.CategoryID
	}
	return ""
}

func (b *Budget) SetCategoryID(categoryID *string) {
	b.CategoryID = categoryID
}

func (b *Budget) GetBudgetType() uint32 {
	if b != nil && b.BudgetType != nil {
		return *b.BudgetType
	}
	return 0
}

func (b *Budget) SetBudgetType(budgetType *uint32) {
	b.BudgetType = budgetType
}

func (b *Budget) GetBudgetStatus() uint32 {
	if b != nil && b.BudgetStatus != nil {
		return *b.BudgetStatus
	}
	return 0
}

func (b *Budget) SetBudgetStatus(budgetStatus *uint32) {
	b.BudgetStatus = budgetStatus
}

func (b *Budget) GetAmount() float64 {
	if b != nil && b.Amount != nil {
		return *b.Amount
	}
	return 0
}

func (b *Budget) SetAmount(amount *float64) {
	b.Amount = amount
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

func (b *Budget) SetCreateTime(createTime *uint64) {
	b.CreateTime = createTime
}

func (b *Budget) GetUpdateTime() uint64 {
	if b != nil && b.UpdateTime != nil {
		return *b.UpdateTime
	}
	return 0
}

func (b *Budget) SetUpdateTime(updateTime *uint64) {
	b.UpdateTime = updateTime
}

func (b *Budget) GetStartDate() uint64 {
	if b != nil && b.StartDate != nil {
		return *b.StartDate
	}
	return 0
}

func (b *Budget) SetStartDate(startDate *uint64) {
	b.StartDate = startDate
}

func (b *Budget) GetEndDate() uint64 {
	if b != nil && b.EndDate != nil {
		return *b.EndDate
	}
	return 0
}

func (b *Budget) SetEndDate(endDate *uint64) {
	b.EndDate = endDate
}

func (b *Budget) CanBudgetUnderCategory(c *Category) error {
	if !c.CanAddBudget() {
		return ErrBudgetNotAllowed
	}
	return nil
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

func (b *Budget) IsRepeatAllTime() bool {
	return b.GetStartDate() == 0 && b.GetEndDate() == 0
}
