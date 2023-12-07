package metric

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var GetMetricsValidator = validator.MustForm(map[string]validator.Validator{
	"app_meta": entity.AppMetaValidator(),
	"metric_type": &validator.UInt32{
		Optional:   true,
		Validators: []validator.UInt32Func{entity.CheckMetricType},
	},
})

func (h *metricHandler) GetMetrics(ctx context.Context, req *presenter.GetMetricsRequest, res *presenter.GetMetricsResponse) error {
	user := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.metricUseCase.GetMetrics(ctx, req.ToUseCaseReq(user))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get metrics, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
