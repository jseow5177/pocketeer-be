package entity

import "github.com/jseow5177/pockteer-be/pkg/goutil"

const (
	DefaultSecurityRegion   = "US"
	DefaultSecurityCurrency = "USD"
)

type SecurityType uint32

const (
	SecurityTypeOther SecurityType = iota
	SecurityTypeCommonStock
	SecurityTypeETF
)

type Security struct {
	SecurityID   *string
	Symbol       *string
	SecurityName *string
	SecurityType *uint32
	Region       *string
	Currency     *string
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

func NewSecurity(symbol string, opts ...SecurityOption) *Security {
	s := &Security{
		Symbol:       goutil.String(symbol),
		SecurityType: goutil.Uint32(uint32(SecurityTypeCommonStock)),
		Region:       goutil.String(DefaultSecurityRegion),
		Currency:     goutil.String(DefaultSecurityCurrency),
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

type Quote struct {
	LatestPrice   *float64
	Change        *float64 // change in price between LatestPrice and PreviousClose
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
