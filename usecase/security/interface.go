package security

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/api"
	"github.com/jseow5177/pockteer-be/entity"
)

type UseCase interface {
	SearchSecurities(ctx context.Context, req *SearchSecuritiesRequest) (*SearchSecuritiesResponse, error)
}

type SearchSecuritiesRequest struct {
	Keyword *string
}

func (m *SearchSecuritiesRequest) GetKeyword() string {
	if m != nil && m.Keyword != nil {
		return *m.Keyword
	}
	return ""
}

func (m *SearchSecuritiesRequest) ToSecurityFilter() *api.SecurityFilter {
	return &api.SecurityFilter{
		Keyword: m.Keyword,
	}
}

type SearchSecuritiesResponse struct {
	Securities []*entity.Security
}

func (m *SearchSecuritiesResponse) GetSecurities() []*entity.Security {
	if m != nil && m.Securities != nil {
		return m.Securities
	}
	return nil
}