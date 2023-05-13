package presenter

type Paging struct {
	Limit *int `json:"limit"`
	Page  *int `json:"page"`
}

func (p *Paging) GetLimit() int {
	if p != nil && p.Limit != nil {
		return *p.Limit
	}
	return 0
}

func (p *Paging) GetPage() int {
	if p != nil && p.Page != nil {
		return *p.Page
	}
	return 0
}

type UInt64Filter struct {
	Gte *uint64 `json:"gte"`
	Lte *uint64 `json:"lte"`
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
