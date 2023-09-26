package snapshot

import (
	"github.com/jseow5177/pockteer-be/usecase/snapshot"
)

type snapshotHandler struct {
	snapshotUseCase snapshot.UseCase
}

func NewSnapshotHandler(snapshotUseCase snapshot.UseCase) *snapshotHandler {
	return &snapshotHandler{
		snapshotUseCase,
	}
}
