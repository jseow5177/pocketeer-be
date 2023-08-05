package entity

import (
	"github.com/jseow5177/pockteer-be/pkg/goutil"
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
