package exchangeratehost

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/dep/api"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/httputil"
	"github.com/jseow5177/pockteer-be/util"
)

const layout = "2006-01-02"

var (
	ErrInvalidDate = errors.New("invalid date")
)

type exchangeRate struct {
	Base  *string            `json:"base,omitempty"`
	Date  *string            `json:"date,omitempty"`
	Rates map[string]float64 `json:"rates,omitempty"`
}

func (er *exchangeRate) GetBase() string {
	if er != nil && er.Base != nil {
		return *er.Base
	}
	return ""
}

func (er *exchangeRate) GetDate() string {
	if er != nil && er.Date != nil {
		return *er.Date
	}
	return ""
}

func (er *exchangeRate) GetRates() map[string]float64 {
	if er != nil && er.Rates != nil {
		return er.Rates
	}
	return nil
}

type exchangeRateHostMgr struct {
	baseURL string
}

func NewExchangeRateHostMgr(cfg *config.ExchangeRateHost) api.ExchangeRateAPI {
	return &exchangeRateHostMgr{
		baseURL: cfg.BaseURL,
	}
}

// Doc: https://exchangerate.host/#/docs
func (mgr *exchangeRateHostMgr) GetExchangeRates(ctx context.Context, erf *api.ExchangeRateFilter) ([]*entity.ExchangeRate, error) {
	ep := "latest"
	if erf.GetDate() != "" {
		t, err := util.ParseDate(erf.GetDate())
		if err != nil {
			return nil, err
		}
		ep = t.Format(layout)
	}

	url := fmt.Sprintf("%s/%s", mgr.baseURL, ep)

	params := map[string]string{
		"base":    erf.GetBase(),
		"symbols": strings.Join(erf.GetSymbols(), ","),
	}

	code, data, err := httputil.SendGetRequest(url, params, nil)
	if err != nil {
		return nil, err
	}

	if code != http.StatusOK {
		return nil, fmt.Errorf("fail to get exchange rate, code: %v", code)
	}

	er := new(exchangeRate)
	if err = json.Unmarshal(data, &er); err != nil {
		return nil, err
	}

	t, err := time.Parse(layout, er.GetDate())
	if err != nil {
		return nil, err
	}
	ts := t.UnixMilli()

	ers := make([]*entity.ExchangeRate, 0)
	for c, r := range er.Rates {
		ers = append(ers, entity.NewExchangeRate(er.GetBase(), c, r, uint64(ts)))
	}

	return ers, nil
}
