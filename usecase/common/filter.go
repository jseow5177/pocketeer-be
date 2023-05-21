package common

type Paging struct {
	Limit *uint32
	Page  *uint32
}

func (p *Paging) GetLimit() uint32 {
	if p != nil && p.Limit != nil {
		return *p.Limit
	}
	return 0
}

func (p *Paging) GetPage() uint32 {
	if p != nil && p.Page != nil {
		return *p.Page
	}
	return 0
}

type UInt64Filter struct {
	Gte *uint64
	Lte *uint64
}

func (m *UInt64Filter) GetGte() uint64 {
	if m != nil && m.Gte != nil {
		return *m.Gte
	}
	return 0
}

func (m *UInt64Filter) GetLte() uint64 {
	if m != nil && m.Lte != nil {
		return *m.Lte
	}
	return 0
}
