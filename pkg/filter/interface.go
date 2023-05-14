package filter

type Sort interface {
	GetField() *string
	GetOrder() *string
}

type FilterOptions interface {
	GetLimit() *uint32
	GetPage() *uint32
	GetSorts() []Sort
}
