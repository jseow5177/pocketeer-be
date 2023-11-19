package metric

import (
	"context"

	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/account"
	"github.com/jseow5177/pockteer-be/usecase/common"
	"github.com/jseow5177/pockteer-be/usecase/transaction"
)

type UseCase interface {
	GetMetrics(ctx context.Context, req *GetMetricsRequest) (*GetMetricsResponse, error)
}

type GetMetricsRequest struct {
	AppMeta    *common.AppMeta
	User       *entity.User
	MetricType *uint32
}

func (m *GetMetricsRequest) GetMetricType() uint32 {
	if m != nil && m.MetricType != nil {
		return *m.MetricType
	}
	return 0
}

func (m *GetMetricsRequest) GetUser() *entity.User {
	if m != nil && m.User != nil {
		return m.User
	}
	return nil
}

func (m *GetMetricsRequest) ToGetAccountsRequest() *account.GetAccountsRequest {
	return &account.GetAccountsRequest{
		UserID: m.User.UserID,
	}
}

func (m *GetMetricsRequest) ToGetTransactionsSummaryRequest() *transaction.GetTransactionsSummaryRequest {
	return &transaction.GetTransactionsSummaryRequest{
		AppMeta:  m.AppMeta,
		User:     m.User,
		Unit:     goutil.Uint32(uint32(entity.SnapshotUnitMonth)),
		Interval: goutil.Uint32(0),
	}
}

type GetMetricsResponse struct {
	Metrics []*entity.Metric
}
