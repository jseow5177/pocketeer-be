package lot

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
)

type lotUseCase struct {
	lotRepo repo.LotRepo
}

func NewLotUseCase(lotRepo repo.LotRepo) UseCase {
	return &lotUseCase{
		lotRepo,
	}
}

func (uc *lotUseCase) CreateLot(ctx context.Context, req *CreateLotRequest) (*CreateLotResponse, error) {
	return nil, nil
}
