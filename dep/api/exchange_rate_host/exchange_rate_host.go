package exchangeratehost

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/dep/api"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/httputil"
	"github.com/jseow5177/pockteer-be/util"
)

const layout = "2006-01-02"

type exchangeRate struct {
	Code  *string  `json:"code,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func (er *exchangeRate) GetCode() string {
	if er != nil && er.Code != nil {
		return *er.Code
	}
	return ""
}

func (er *exchangeRate) GetValue() float64 {
	if er != nil && er.Value != nil {
		return *er.Value
	}
	return 0
}

type response struct {
	Data    map[string]*exchangeRate `json:"data,omitempty"`
	Message *string                  `json:"message,omitempty"`
	Errors  map[string][]string      `json:"errors,omitempty"`
}

func (r *response) GetData() map[string]*exchangeRate {
	if r != nil && r.Data != nil {
		return r.Data
	}
	return nil
}

func (r *response) GetErrors() map[string][]string {
	if r != nil && r.Errors != nil {
		return r.Errors
	}
	return nil
}

func (r *response) GetMessage() string {
	if r != nil && r.Message != nil {
		return *r.Message
	}
	return ""
}

type exchangeRateHostMgr struct {
	apiKey  string
	baseURL string
}

func NewExchangeRateHostMgr(cfg *config.ExchangeRateHost) api.ExchangeRateAPI {
	return &exchangeRateHostMgr{
		apiKey:  cfg.APIKey,
		baseURL: cfg.BaseURL,
	}
}

// Doc: https://currencyapi.com/docs/historical
func (mgr *exchangeRateHostMgr) GetExchangeRates(ctx context.Context, erf *api.ExchangeRateFilter) ([]*entity.ExchangeRate, error) {
	params := map[string]string{
		"base_currency": erf.GetBase(),
		"currencies":    strings.Join(erf.GetCurrencies(), ","),
		"apikey":        mgr.apiKey,
	}

	date, err := util.ParseDate(erf.GetDate())
	if err != nil {
		return nil, err
	}
	params["date"] = date.Format(layout)

	url := fmt.Sprintf("%s/historical", mgr.baseURL)
	code, data, err := httputil.SendGetRequest(url, params, nil)
	if err != nil {
		return nil, err
	}

	resp := new(response)
	if err = json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	if code != http.StatusOK {
		return nil, fmt.Errorf("fail to get exchange rate, code: %v, errors: %v, message: %v",
			code, resp.GetErrors(), resp.GetMessage())
	}

	exchangeRates := make([]*entity.ExchangeRate, 0)
	for _, er := range resp.Data {
		ts := date.UnixMilli()

		exchangeRates = append(exchangeRates, entity.NewExchangeRate(
			erf.GetBase(),
			er.GetCode(),
			er.GetValue(),
			uint64(ts),
		))
	}

	return exchangeRates, nil
}
