package api

import (
	"context"

	"github.com/jseow5177/pockteer-be/entity"
)

type SecurityAPI interface {
	SearchSecurities(ctx context.Context, keyword string) ([]*entity.Security, error)
	GetLatestQuote(ctx context.Context, symbol string) (*entity.Quote, error)
}
