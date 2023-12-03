package entity

import (
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/util"
)

type MetricUnit string

const (
	MetricUnitPercent MetricUnit = "%"
)

type MetricType uint32

const (
	MetricTypeInvalid MetricType = iota
	MetricTypeNetWorth
	MetricTypeSavings
)

var MetricTypes = map[uint32]string{
	uint32(MetricTypeNetWorth): "net worth",
	uint32(MetricTypeSavings):  "savings",
}

type MetricID uint32

const (
	MetricIDInvalid MetricID = iota
	MetricIDDebtRatio
	MetricIDSavingsRatio
	MetricIDInvestmentsToNetWorthRatio
)

var MetricIDs = map[uint32]string{
	uint32(MetricIDDebtRatio):                  "Debt Ratio",
	uint32(MetricIDSavingsRatio):               "Savings Ratio",
	uint32(MetricIDInvestmentsToNetWorthRatio): "Investments to Net Worth Ratio",
}

type MetricStatus uint32

const (
	MetricStatusHealthy MetricStatus = iota
	MetricStatusUnhealthy
)

var MetricThresholds = map[uint32]string{
	uint32(MetricIDDebtRatio):                  "50% or less",
	uint32(MetricIDSavingsRatio):               "20% or more",
	uint32(MetricIDInvestmentsToNetWorthRatio): "50% or more",
}

var MetricHealthChecks = map[uint32]func(value float64) bool{
	uint32(MetricIDDebtRatio): func(value float64) bool {
		return value < float64(50)
	},
	uint32(MetricIDSavingsRatio): func(value float64) bool {
		return value > float64(20)
	},
	uint32(MetricIDInvestmentsToNetWorthRatio): func(value float64) bool {
		return value > float64(50)
	},
}

type Metric struct {
	ID        *uint32
	Name      *string
	Value     *float64
	Type      *uint32
	Unit      *string
	Status    *uint32
	Threshold *string
}

type MetricOption func(mt *Metric)

func WithMetricValue(value *float64) MetricOption {
	return func(mt *Metric) {
		if value != nil {
			mt.SetValue(value)
		}
	}
}

func WithMetricUnit(unit *string) MetricOption {
	return func(mt *Metric) {
		if unit != nil {
			mt.SetUnit(unit)
		}
	}
}

func WithMetricName(name *string) MetricOption {
	return func(mt *Metric) {
		if name != nil {
			mt.SetName(name)
		}
	}
}

func NewMetric(id, metricType uint32, opts ...MetricOption) *Metric {
	mt := &Metric{
		ID:   goutil.Uint32(id),
		Type: goutil.Uint32(metricType),
	}
	for _, opt := range opts {
		opt(mt)
	}

	status := uint32(MetricStatusHealthy)
	if !MetricHealthChecks[id](mt.GetValue()) {
		status = uint32(MetricStatusUnhealthy)
	}
	mt.SetStatus(goutil.Uint32(status))

	mt.SetThreshold(goutil.String(MetricThresholds[id]))

	return mt
}

func (m *Metric) GetID() uint32 {
	if m != nil && m.ID != nil {
		return *m.ID
	}
	return 0
}

func (m *Metric) SetID(id *uint32) {
	m.ID = id
}

func (m *Metric) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *Metric) SetName(name *string) {
	m.Name = name
}

func (m *Metric) GetValue() float64 {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return 0
}

func (m *Metric) SetValue(value *float64) {
	m.Value = value

	if value != nil {
		v := util.RoundFloatToStandardDP(*value)
		m.Value = goutil.Float64(v)
	}
}

func (m *Metric) GetType() uint32 {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return 0
}

func (m *Metric) SetType(metricType *uint32) {
	m.Type = metricType
}

func (m *Metric) GetUnit() string {
	if m != nil && m.Unit != nil {
		return *m.Unit
	}
	return ""
}

func (m *Metric) SetUnit(unit *string) {
	m.Unit = unit
}

func (m *Metric) GetStatus() uint32 {
	if m != nil && m.Status != nil {
		return *m.Status
	}
	return 0
}

func (m *Metric) SetStatus(status *uint32) {
	m.Status = status
}

func (m *Metric) GetThreshold() string {
	if m != nil && m.Threshold != nil {
		return *m.Threshold
	}
	return ""
}

func (m *Metric) SetThreshold(threshold *string) {
	m.Threshold = threshold
}
