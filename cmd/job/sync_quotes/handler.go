package syncquotes

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/api"
	"github.com/jseow5177/pockteer-be/dep/repo"
)

type SyncQuotesHandler struct {
	cmd *SyncQuotesCmd
}

type SyncQuotesRequest struct{}

type SyncQuotesResponse struct{}

func NewSyncQuotesHandler(
	quoteRepo repo.QuoteRepo,
	holdingRepo repo.HoldingRepo,
	securityAPI api.SecurityAPI,
) *SyncQuotesHandler {
	return &SyncQuotesHandler{
		cmd: &SyncQuotesCmd{
			quoteRepo:   quoteRepo,
			holdingRepo: holdingRepo,
			securityAPI: securityAPI,
		},
	}
}

func (h *SyncQuotesHandler) SyncQuotes(ctx context.Context, req *SyncQuotesRequest, res *SyncQuotesResponse) error {
	return h.cmd.Run(ctx)
}
