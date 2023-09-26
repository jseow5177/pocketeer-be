package presenter

import (
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/snapshot"
)

type Snapshot struct {
	Period *string `json:"period,omitempty"`
	Value  *string `json:"value,omitempty"`
}

func (m *Snapshot) GetPeriod() string {
	if m != nil && m.Period != nil {
		return *m.Period
	}
	return ""
}

func (m *Snapshot) GetValue() string {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return ""
}

type GetAccountSnapshotsRequest struct {
	AccountID *string `json:"account_id,omitempty"`
	Unit      *uint32 `json:"unit,omitempty"`
	Interval  *uint32 `json:"interval,omitempty"`
}

func (m *GetAccountSnapshotsRequest) GetAccountID() string {
	if m != nil && m.AccountID != nil {
		return *m.AccountID
	}
	return ""
}

func (m *GetAccountSnapshotsRequest) GetUnit() uint32 {
	if m != nil && m.Unit != nil {
		return *m.Unit
	}
	return 0
}

func (m *GetAccountSnapshotsRequest) GetInterval() uint32 {
	if m != nil && m.Interval != nil {
		return *m.Interval
	}
	return 0
}

func (m *GetAccountSnapshotsRequest) ToUseCaseReq(userID string) *snapshot.GetAccountSnapshotsRequest {
	return &snapshot.GetAccountSnapshotsRequest{
		UserID:    goutil.String(userID),
		AccountID: m.AccountID,
		Unit:      m.Unit,
		Interval:  m.Interval,
	}
}

type GetAccountSnapshotsResponse struct {
	Snapshots []*Snapshot `json:"snapshots,omitempty"`
}

func (m *GetAccountSnapshotsResponse) GetSnapshots() []*Snapshot {
	if m != nil && m.Snapshots != nil {
		return m.Snapshots
	}
	return nil
}

func (m *GetAccountSnapshotsResponse) Set(useCaseRes *snapshot.GetAccountSnapshotsResponse) {
	m.Snapshots = toSnapshots(useCaseRes.Snapshots)
}
