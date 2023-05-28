package presenter

type Paging struct {
	Limit *uint32 `json:"limit,omitempty"`
	Page  *uint32 `json:"page,omitempty"`
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
	Gte *uint64 `json:"gte,omitempty"`
	Lte *uint64 `json:"lte,omitempty"`
}

func (uv *UInt64Filter) GetGte() uint64 {
	if uv != nil && uv.Gte != nil {
		return *uv.Gte
	}
	return 0
}

func (uv *UInt64Filter) GetLte() uint64 {
	if uv != nil && uv.Lte != nil {
		return *uv.Lte
	}
	return 0
}
