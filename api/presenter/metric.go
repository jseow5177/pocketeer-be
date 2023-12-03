package presenter

import (
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/usecase/metric"
)

type Metric struct {
	ID        *uint32 `json:"id,omitempty"`
	Name      *string `json:"name,omitempty"`
	Value     *string `json:"value,omitempty"`
	Type      *uint32 `json:"type,omitempty"`
	Unit      *string `json:"unit,omitempty"`
	Status    *uint32 `json:"status,omitempty"`
	Threshold *string `json:"threshold,omitempty"`
}

func (m *Metric) GetID() uint32 {
	if m != nil && m.ID != nil {
		return *m.ID
	}
	return 0
}

func (m *Metric) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *Metric) GetValue() string {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return ""
}

func (m *Metric) GetType() uint32 {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return 0
}

func (m *Metric) GetUnit() string {
	if m != nil && m.Unit != nil {
		return *m.Unit
	}
	return ""
}

func (m *Metric) GetStatus() uint32 {
	if m != nil && m.Status != nil {
		return *m.Status
	}
	return 0
}

func (m *Metric) GetThreshold() string {
	if m != nil && m.Threshold != nil {
		return *m.Threshold
	}
	return ""
}

type GetMetricsRequest struct {
	AppMeta    *AppMeta `json:"app_meta,omitempty"`
	MetricType *uint32  `json:"metric_type,omitempty"`
}

func (m *GetMetricsRequest) GetAppMeta() *AppMeta {
	if m != nil && m.AppMeta != nil {
		return m.AppMeta
	}
	return nil
}

func (m *GetMetricsRequest) GetMetricType() uint32 {
	if m != nil && m.MetricType != nil {
		return *m.MetricType
	}
	return 0
}

func (m *GetMetricsRequest) ToUseCaseReq(user *entity.User) *metric.GetMetricsRequest {
	return &metric.GetMetricsRequest{
		AppMeta:    m.AppMeta.toAppMeta(),
		User:       user,
		MetricType: m.MetricType,
	}
}

type GetMetricsResponse struct {
	Metrics []*Metric `json:"metrics,omitempty"`
}

func (m *GetMetricsResponse) Set(useCaseRes *metric.GetMetricsResponse) {
	m.Metrics = toMetrics(useCaseRes.Metrics)
}
