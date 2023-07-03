package iex

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
)

var securityTypes = map[string]entity.SecurityType{
	"cs": entity.SecurityTypeCommonStock,
	"et": entity.SecurityTypeETF,
}

type security struct {
	Symbol       *string `json:"symbol,omitempty"`
	SecurityName *string `json:"securityName,omitempty"`
	SecurityType *string `json:"securityType,omitempty"`
	Region       *string `json:"region,omitempty"`
	Currency     *string `json:"currency,omitempty"`
}

func (s *security) GetSymbol() string {
	if s != nil && s.Symbol != nil {
		return *s.Symbol
	}
	return ""
}

func (s *security) GetSecurityName() string {
	if s != nil && s.SecurityName != nil {
		return *s.SecurityName
	}
	return ""
}

func (s *security) GetSecurityType() string {
	if s != nil && s.SecurityType != nil {
		return *s.SecurityType
	}
	return ""
}

func (s *security) GetRegion() string {
	if s != nil && s.Region != nil {
		return *s.Region
	}
	return ""
}

func (s *security) GetCurrency() string {
	if s != nil && s.Currency != nil {
		return *s.Currency
	}
	return ""
}

type iexMgr struct {
	baseURL          string
	token            string
	supportedRegions []string
}

func NewIEXMgr(cfg *config.IEX) *iexMgr {
	return &iexMgr{
		baseURL:          cfg.BaseURL,
		token:            cfg.Token,
		supportedRegions: cfg.SupportedRegions,
	}
}

// Doc: https://iexcloud.io/docs/core/autocomplete-search
func (mgr *iexMgr) SearchSecurities(ctx context.Context, sf *api.SecurityFilter) ([]*entity.Security, error) {
	url := fmt.Sprintf("%s/search/%s", mgr.baseURL, sf.GetKeyword())

	code, data, err := httputil.SendGetRequest(url, mgr.getTokenParam(), nil)
	if err != nil {
		return nil, err
	}

	if code != http.StatusOK {
		return nil, fmt.Errorf("fail to search securities, code: %v", code)
	}

	iexSS := make([]*security, 0)
	if err = json.Unmarshal(data, &iexSS); err != nil {
		return nil, err
	}

	ss := make([]*entity.Security, 0)
	for _, s := range iexSS {
		securityType, ok := securityTypes[s.GetSecurityType()]
		if !ok {
			securityType = entity.SecurityTypeOther
		}

		if len(mgr.supportedRegions) != 0 && !goutil.ContainString(mgr.supportedRegions, s.GetRegion()) {
			continue
		}

		ss = append(ss, entity.NewSecurity(
			s.GetSymbol(),
			entity.WithSecurityCurrency(s.Currency),
			entity.WithSecurityRegion(s.Region),
			entity.WithSecurityType(goutil.Uint32(uint32(securityType))),
			entity.WithSecurityName(s.SecurityName),
		))
	}

	return ss, nil
}

// Doc: https://iexcloud.io/docs/core/QUOTE
func (mgr *iexMgr) GetLatestQuote(ctx context.Context, sf *api.SecurityFilter) (*entity.Quote, error) {
	return nil, nil
}

func (mgr *iexMgr) ListSymbols(ctx context.Context, sf *api.SecurityFilter) ([]*entity.Security, error) {
	return nil, nil
}

func (mgr *iexMgr) getTokenParam() map[string]string {
	return map[string]string{
		"token": mgr.token,
	}
}
