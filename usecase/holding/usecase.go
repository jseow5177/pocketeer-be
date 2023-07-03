package holding

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
)

type holdingUseCase struct {
	txMgr       repo.TxMgr
	holdingRepo repo.HoldingRepo
}

func NewHoldingUseCase(txMgr repo.TxMgr, holdingRepo repo.HoldingRepo) UseCase {
	return &holdingUseCase{
		txMgr,
		holdingRepo,
	}
}

func (uc *holdingUseCase) CreateHolding(ctx context.Context, req *CreateHoldingRequest) (*CreateHoldingResponse, error) {
	return nil, nil
}
