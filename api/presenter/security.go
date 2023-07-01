package presenter

import (
	"github.com/jseow5177/pockteer-be/usecase/security"
)

type Security struct {
	Symbol       *string `json:"symbol,omitempty"`
	SecurityName *string `json:"security_name,omitempty"`
	SecurityType *uint32 `json:"security_type,omitempty"`
	Region       *string `json:"region,omitempty"`
	Currency     *string `json:"currency,omitempty"`
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

type SearchSecuritiesRequest struct {
	Keyword *string `json:"keyword,omitempty"`
}

func (m *SearchSecuritiesRequest) GetKeyword() string {
	if m != nil && m.Keyword != nil {
		return *m.Keyword
	}
	return ""
}

func (m *SearchSecuritiesRequest) ToUseCaseReq() *security.SearchSecuritiesRequest {
	return &security.SearchSecuritiesRequest{
		Keyword: m.Keyword,
	}
}

type SearchSecuritiesResponse struct {
	Securities []*Security `json:"securities,omitempty"`
}

func (m *SearchSecuritiesResponse) GetSecurities() []*Security {
	if m != nil && m.Securities != nil {
		return m.Securities
	}
	return nil
}

func (m *SearchSecuritiesResponse) Set(useCaseRes *security.SearchSecuritiesResponse) {
	ss := make([]*Security, 0)
	for _, s := range useCaseRes.Securities {
		ss = append(ss, toSecurity(s))
	}
	m.Securities = ss
}
