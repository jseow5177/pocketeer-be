package common

type AggrType uint32

const (
	AggrTypeSum AggrType = 0
)

type Aggr struct {
	Field *string
	Type  *uint32
}

func (m *Aggr) GetField() string {
	if m != nil && m.Field != nil {
		return *m.Field
	}
	return ""
}

func (m *Aggr) GetType() uint32 {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return 0
}
