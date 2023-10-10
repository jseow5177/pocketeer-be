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
	return nil, nil
}
