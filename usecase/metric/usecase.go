package metric

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/account"
	"github.com/jseow5177/pockteer-be/usecase/transaction"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

var (
	ErrInvalidTransactionsSummary = errors.New("invalid number of transactions summary")
)

type metricUseCase struct {
	accountUseCase     account.UseCase
	transactionUseCase transaction.UseCase
}

func NewMetricUseCase(
	accountUseCase account.UseCase,
	transactionUseCase transaction.UseCase,
) UseCase {
	return &metricUseCase{
		accountUseCase,
		transactionUseCase,
	}
}

func (uc *metricUseCase) GetMetrics(ctx context.Context, req *GetMetricsRequest) (*GetMetricsResponse, error) {
	metrics := make([]*entity.Metric, 0)

	if req.MetricType == nil || req.GetMetricType() == uint32(entity.MetricTypeNetWorth) {
		getAccountsRes, err := uc.accountUseCase.GetAccounts(ctx, req.ToGetAccountsRequest())
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to get accounts, err: %v", err)
			return nil, err
		}

		metrics = append(metrics, entity.NewMetric(
			uint32(entity.MetricIDDebtRatio),
			uint32(entity.MetricTypeNetWorth),
			entity.WithMetricName(goutil.String(entity.MetricIDs[uint32(entity.MetricIDDebtRatio)])),
			entity.WithMetricValue(goutil.Float64(getAccountsRes.GetDebtRatio())),
			entity.WithMetricUnit(goutil.String(string(entity.MetricUnitPercent))),
		))

		metrics = append(metrics, entity.NewMetric(
			uint32(entity.MetricIDInvestmentsToNetWorthRatio),
			uint32(entity.MetricTypeNetWorth),
			entity.WithMetricName(goutil.String(entity.MetricIDs[uint32(entity.MetricIDInvestmentsToNetWorthRatio)])),
			entity.WithMetricValue(goutil.Float64(getAccountsRes.GetInvestmentsToNetWorthRatio())),
			entity.WithMetricUnit(goutil.String(string(entity.MetricUnitPercent))),
		))
	}

	if req.MetricType == nil || req.GetMetricType() == uint32(entity.MetricTypeSavings) {
		getTransactionsSummaryRes, err := uc.transactionUseCase.GetTransactionsSummary(ctx, req.ToGetTransactionsSummaryRequest())
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to get transactions summary, err: %v", err)
			return nil, err
		}

		if len(getTransactionsSummaryRes.Summary) != 1 {
			log.Ctx(ctx).Error().Msg("wrong number of transactions summary")
			return nil, ErrInvalidTransactionsSummary
		}

		summary := getTransactionsSummaryRes.Summary[0]

		// only calculate savings ratio if
		// 1. Income is more than 0
		// 2. Savings is more than 0
		var savingsRatio float64
		if summary.GetTotalIncome() > 0 && summary.GetSum() > 0 {
			savingsRatio = summary.GetSum() / summary.GetTotalIncome()
		}
		savingsRatio = util.RoundFloatToStandardDP(savingsRatio * 100)

		metrics = append(metrics, entity.NewMetric(
			uint32(entity.MetricIDSavingsRatio),
			uint32(entity.MetricTypeSavings),
			entity.WithMetricName(goutil.String(entity.MetricIDs[uint32(entity.MetricIDSavingsRatio)])),
			entity.WithMetricValue(goutil.Float64(savingsRatio)),
			entity.WithMetricUnit(goutil.String(string(entity.MetricUnitPercent))),
		))
	}

	return &GetMetricsResponse{
		Metrics: metrics,
	}, nil
}
