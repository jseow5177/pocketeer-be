package security

import (
	"context"
	"fmt"
	"strings"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/filter"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

type UseCase interface {
	SearchSecurities(ctx context.Context, req *SearchSecuritiesRequest) (*SearchSecuritiesResponse, error)
}

type SearchSecuritiesRequest struct {
	Symbol *string
}

func (m *SearchSecuritiesRequest) GetSymbol() string {
	if m != nil && m.Symbol != nil {
		return *m.Symbol
	}
	return ""
}

func (m *SearchSecuritiesRequest) ToSecurityFilter() *repo.SecurityFilter {
	return repo.NewSecurityFilter(
		repo.WithSecuritySymbolRegex(goutil.String(fmt.Sprintf("^%s.*", strings.ToUpper(m.GetSymbol())))),
		repo.WithSecurityPaging(&repo.Paging{
			Sorts: []filter.Sort{
				&repo.Sort{
					Field: goutil.String("symbol"),
					Order: goutil.String(config.OrderAsc),
				},
			},
		}),
	)
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
