package filter

type BoolOp string

var (
	And BoolOp = "and"
	Or  BoolOp = "or"
)

type Sort interface {
	GetField() *string
	GetOrder() *string
}

type FilterOptions interface {
	GetLimit() *uint32
	GetPage() *uint32
	GetSorts() []Sort
}

type Query interface {
	GetQueries() []Query
	GetFilters() []interface{}
	GetOp() BoolOp
}
