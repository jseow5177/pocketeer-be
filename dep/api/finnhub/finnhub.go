package finnhub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/dep/api"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/pkg/httputil"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
)

type quote struct {
	// Current price
	C *float64 `json:"c,omitempty"`
	// Previous close price
	Pc *float64 `json:"pc,omitempty"`
	// Change
	D *float64 `json:"d,omitempty"`
	// Percent change
	Dp *float64 `json:"dp,omitempty"`
	// Update time
	T *uint64 `json:"t,omitempty"`
}

func (q *quote) GetC() float64 {
	if q != nil && q.C != nil {
		return *q.C
	}
	return 0
}

func (q *quote) GetPc() float64 {
	if q != nil && q.Pc != nil {
		return *q.Pc
	}
	return 0
}

func (q *quote) GetD() float64 {
	if q != nil && q.D != nil {
		return *q.D
	}
	return 0
}

func (q *quote) GetDp() float64 {
	if q != nil && q.Dp != nil {
		return *q.Dp
	}
	return 0
}

func (q *quote) GetT() uint64 {
	if q != nil && q.T != nil {
		return *q.T
	}
	return 0
}

var securityTypes = map[string]entity.SecurityType{
	"Common Stock": entity.SecurityTypeCommonStock,
	"ETP":          entity.SecurityTypeETF,
}

type finnhubMgr struct {
	baseURL string
	token   string

	client *finnhub.DefaultApiService
}

func NewFinnHubMgr(cfg *config.FinnHub) api.SecurityAPI {
	fhCfg := finnhub.NewConfiguration()
	fhCfg.AddDefaultHeader("X-Finnhub-Token", cfg.Token)

	return &finnhubMgr{
		client:  finnhub.NewAPIClient(fhCfg).DefaultApi,
		baseURL: cfg.BaseURL,
		token:   cfg.Token,
	}
}

// Doc: https://finnhub.io/docs/api/symbol-search
func (mgr *finnhubMgr) SearchSecurities(ctx context.Context, sf *api.SecurityFilter) ([]*entity.Security, error) {
	l, _, err := mgr.client.SymbolSearch(ctx).Q(sf.GetSymbol()).Execute()
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
		))
	}

	return ss, nil
}

// Doc: https://finnhub.io/docs/api/quote
func (mgr *finnhubMgr) GetLatestQuote(ctx context.Context, sf *api.SecurityFilter) (*entity.Quote, error) {
	url := fmt.Sprintf("%s/quote", mgr.baseURL)

	queryParams := map[string]string{
		"token":  mgr.token,
		"symbol": sf.GetSymbol(),
	}

	code, data, err := httputil.SendGetRequest(url, queryParams, nil)
	if err != nil {
		return nil, err
	}

	if code != http.StatusOK {
		return nil, fmt.Errorf("fail to get latest quote, code: %v", code)
	}

	q := new(quote)
	if err = json.Unmarshal(data, &q); err != nil {
		return nil, err
	}

	// to milli
	t := q.GetT() * 1000

	return entity.NewQuote(
		sf.GetSymbol(),
		entity.WithQuoteLatestPrice(q.C),
		entity.WithQuoteChange(q.D),
		entity.WithQuoteChangePercent(q.Dp),
		entity.WithQuotePreviousClose(q.Pc),
		entity.WithQuoteUpdateTime(goutil.Uint64(t)),
		entity.WithQuoteCurrency(goutil.String(string(entity.CurrencyUSD))), // only USD for now
	), nil
}

// Doc: https://finnhub.io/docs/api/stock-symbols
func (mgr *finnhubMgr) ListSymbols(ctx context.Context, sf *api.SecurityFilter) ([]*entity.Security, error) {
	res, _, err := mgr.client.StockSymbols(ctx).Exchange(sf.GetExchange()).Execute()
	if err != nil {
		return nil, fmt.Errorf("fail to list stock symbols, err: %v", err)
	}

	ss := make([]*entity.Security, 0, len(res))
	for _, r := range res {
		securityType, ok := securityTypes[r.GetType()]
		if !ok {
			securityType = entity.SecurityTypeOther
		}
		ss = append(ss, entity.NewSecurity(
			r.GetSymbol(),
			entity.WithSecurityName(r.Description),
			entity.WithSecurityType(goutil.Uint32(uint32(securityType))),
			entity.WithSecurityCurrency(r.Currency),
			entity.WithSecurityRegion(sf.Exchange),
		))
	}

	return ss, nil
}
