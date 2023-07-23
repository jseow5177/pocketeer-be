package entity

import (
	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/util"
)

type SecurityType uint32

const (
	SecurityTypeOther SecurityType = iota
	SecurityTypeCommonStock
	SecurityTypeETF
)

type SecurityUpdate struct {
	Quote *Quote
}

func (su *SecurityUpdate) GetQuote() *Quote {
	if su != nil && su.Quote != nil {
		return su.Quote
	}
	return nil
}

type SecurityUpdateOption func(su *SecurityUpdate)

func WithUpdateSecurityQuote(quote *Quote) SecurityUpdateOption {
	return func(su *SecurityUpdate) {
		su.Quote = quote
	}
}

func NewSecurityUpdate(opts ...SecurityUpdateOption) *SecurityUpdate {
	su := new(SecurityUpdate)
	for _, opt := range opts {
		opt(su)
	}
	return su
}

type Security struct {
	SecurityID   *string
	Symbol       *string
	SecurityName *string
	SecurityType *uint32
	Region       *string
	Currency     *string
	Quote        *Quote
}

type SecurityOption = func(s *Security)

func WithSecurityID(securityID *string) SecurityOption {
	return func(s *Security) {
		s.SecurityID = securityID
	}
}

func WithSecurityName(securityName *string) SecurityOption {
	return func(s *Security) {
		s.SecurityName = securityName
	}
}

func WithSecurityType(securityType *uint32) SecurityOption {
	return func(s *Security) {
		s.SecurityType = securityType
	}
}

func WithSecurityRegion(region *string) SecurityOption {
	return func(s *Security) {
		s.Region = region
	}
}

func WithSecurityCurrency(currency *string) SecurityOption {
	return func(s *Security) {
		s.Currency = currency
	}
}

func WithSecurityQuote(quote *Quote) SecurityOption {
	return func(s *Security) {
		s.Quote = quote
	}
}

func NewSecurity(symbol string, opts ...SecurityOption) *Security {
	s := &Security{
		Symbol:       goutil.String(symbol),
		SecurityName: goutil.String(""),
		SecurityType: goutil.Uint32(uint32(SecurityTypeCommonStock)),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Security) GetSecurityID() string {
	if s != nil && s.SecurityID != nil {
		return *s.SecurityID
	}
	return ""
}

func (s *Security) SetSecurityID(securityID *string) {
	s.SecurityID = securityID
}

func (s *Security) GetSymbol() string {
	if s != nil && s.Symbol != nil {
		return *s.Symbol
	}
	return ""
}

func (s *Security) GetSecurityName() string {
	if s != nil && s.SecurityName != nil {
		return *s.SecurityName
	}
	return ""
}

func (s *Security) GetSecurityType() uint32 {
	if s != nil && s.SecurityType != nil {
		return *s.SecurityType
	}
	return 0
}

func (s *Security) GetRegion() string {
	if s != nil && s.Region != nil {
		return *s.Region
	}
	return ""
}

func (s *Security) GetCurrency() string {
	if s != nil && s.Currency != nil {
		return *s.Currency
	}
	return ""
}

func (s *Security) GetQuote() *Quote {
	if s != nil && s.Quote != nil {
		return s.Quote
	}
	return nil
}

type Quote struct {
	LatestPrice   *float64
	Change        *float64 // LatestPrice - PreviousClose
	ChangePercent *float64
	PreviousClose *float64
	UpdateTime    *uint64
}

func (q *Quote) GetLatestPrice() float64 {
	if q != nil && q.LatestPrice != nil {
		return *q.LatestPrice
	}
	return 0
}

func (q *Quote) GetChange() float64 {
	if q != nil && q.Change != nil {
		return *q.Change
	}
	return 0
}

func (q *Quote) GetChangePercent() float64 {
	if q != nil && q.ChangePercent != nil {
		return *q.ChangePercent
	}
	return 0
}

func (q *Quote) GetPreviousClose() float64 {
	if q != nil && q.PreviousClose != nil {
		return *q.PreviousClose
	}
	return 0
}

func (q *Quote) GetUpdateTime() uint64 {
	if q != nil && q.UpdateTime != nil {
		return *q.UpdateTime
	}
	return 0
}

// TODO: Have better currency conversion logic
func (q *Quote) ToSGD() *Quote {
	var (
		lp  = util.RoundFloat(q.GetLatestPrice()*config.USDToSGD, config.StandardDP)
		pc  = util.RoundFloat(q.GetPreviousClose()*config.USDToSGD, config.StandardDP)
		ch  = util.RoundFloat(lp-pc, config.StandardDP)
		chp = util.RoundFloat((lp-pc)*100/pc, config.PreciseDP)
	)
	return &Quote{
		LatestPrice:   goutil.Float64(lp),
		PreviousClose: goutil.Float64(pc),
		Change:        goutil.Float64(ch),
		ChangePercent: goutil.Float64(chp),
		UpdateTime:    q.UpdateTime,
	}
}
