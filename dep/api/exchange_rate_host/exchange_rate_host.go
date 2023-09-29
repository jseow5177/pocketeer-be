package exchangeratehost

import (
	"context"
	"encoding/json"
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

type response struct {
	Data    map[string]map[string]float64 `json:"data,omitempty"`
	Message *string                       `json:"message,omitempty"`
	Errors  map[string][]string           `json:"errors,omitempty"`
}

func (er *response) GetData() map[string]map[string]float64 {
	if er != nil && er.Data != nil {
		return er.Data
	}
	return nil
}

func (er *response) GetErrors() map[string][]string {
	if er != nil && er.Errors != nil {
		return er.Errors
	}
	return nil
}

func (er *response) GetMessage() string {
	if er != nil && er.Message != nil {
		return *er.Message
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

// Doc: https://freecurrencyapi.com/docs/
func (mgr *exchangeRateHostMgr) GetExchangeRates(ctx context.Context, erf *api.ExchangeRateFilter) ([]*entity.ExchangeRate, error) {
	params := map[string]string{
		"base_currency": erf.GetBase(),
		"currencies":    strings.Join(erf.GetCurrencies(), ","),
		"apikey":        mgr.apiKey,
	}

	toDate, err := util.ParseDate(erf.GetToDate())
	if err != nil {
		return nil, err
	}
	params["date_to"] = toDate.Format(layout)

	fromDate, err := util.ParseDate(erf.GetFromDate())
	if err != nil {
		return nil, err
	}
	params["date_from"] = fromDate.Format(layout)

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
	for date, data := range resp.Data {
		t, err := time.Parse(layout, date)
		if err != nil {
			return nil, err
		}
		ts := t.UnixMilli()

		for to, rate := range data {
			exchangeRates = append(exchangeRates, entity.NewExchangeRate(
				erf.GetBase(),
				to,
				rate,
				uint64(ts),
			))
		}
	}

	return exchangeRates, nil
}
