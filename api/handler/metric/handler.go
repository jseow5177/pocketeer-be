package metric

import "github.com/jseow5177/pockteer-be/usecase/metric"

type metricHandler struct {
	metricUseCase metric.UseCase
}

func NewMetricHandler(metricUseCase metric.UseCase) *metricHandler {
	return &metricHandler{
		metricUseCase,
	}
}
