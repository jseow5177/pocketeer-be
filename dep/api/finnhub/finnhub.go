package finnhub

import (
	"context"
	"fmt"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
)

var securityTypes = map[string]entity.SecurityType{
	"Common Stock": entity.SecurityTypeCommonStock,
	"ETP":          entity.SecurityTypeETF,
}

type finnhubMgr struct {
	client *finnhub.DefaultApiService
}

func NewFinnHubMgr(cfg *config.FinnHub) *finnhubMgr {
	fhCfg := finnhub.NewConfiguration()
	fhCfg.AddDefaultHeader("X-Finnhub-Token", cfg.Token)

	return &finnhubMgr{
		client: finnhub.NewAPIClient(fhCfg).DefaultApi,
	}
}

// Doc: https://finnhub.io/docs/api/symbol-search
func (mgr *finnhubMgr) SearchSecurities(ctx context.Context, keyword string) ([]*entity.Security, error) {
	l, _, err := mgr.client.SymbolSearch(ctx).Q(keyword).Execute()
	if err != nil {
		return nil, fmt.Errorf("fail to search securities, err: %v", err)
	}
	res := l.GetResult()

	ss := make([]*entity.Security, 0)
	for _, r := range res {
		securityType, ok := securityTypes[r.GetType()]
		if !ok {
			securityType = entity.SecurityTypeOther
		}
		ss = append(ss, entity.NewSecurity(
			r.GetSymbol(),
			entity.WithSecurityName(r.Description),
			entity.WithSecurityType(goutil.Uint32(uint32(securityType))),
			entity.WithSecurityRegion(goutil.String("")),
			entity.WithSecurityCurrency(goutil.String("")),
		))
	}

	return ss, nil
}

// Doc: https://finnhub.io/docs/api/quote
func (mgr *finnhubMgr) GetLatestQuote(ctx context.Context, symbol string) (*entity.Quote, error) {
	return nil, nil
}
