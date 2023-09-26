package snapshot

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
)

type snapshotUseCase struct {
	snapshotRepo repo.SnapshotRepo
	accountRepo  repo.AccountRepo
}

func NewSnapshotUseCase(
	snapshotRepo repo.SnapshotRepo,
	accountRepo repo.AccountRepo,
) UseCase {
	return &snapshotUseCase{
		snapshotRepo,
		accountRepo,
	}
}

func (uc *snapshotUseCase) GetAccountSnapshots(ctx context.Context, req *GetAccountSnapshotsRequest) (*GetAccountSnapshotsResponse, error) {
	// var (
	// 	ac  *entity.Account
	// 	err error
	// )
	// if req.AccountID != nil {
	// 	ac, err = uc.accountRepo.Get(ctx, req.ToAccountFilter())
	// 	if err != nil {
	// 		log.Ctx(ctx).Error().Msgf("fail to get account from repo, err: %v", err)
	// 		return nil, err
	// 	}
	// }

	// sps, err := uc.snapshotRepo.GetMany(ctx, req.ToSnapshotFilter())
	// if err != nil {
	// 	log.Ctx(ctx).Error().Msgf("fail to get account snapshots from repo, err: %v", err)
	// 	return nil, err
	// }

	return nil, nil
}
