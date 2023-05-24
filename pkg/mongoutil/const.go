package mongoutil

const (
	tagBson   = "bson"
	tagFilter = "filter"
	tagPaging = "paging"
)

const (
	ignore = "-"
	sep    = "__"

	_id = "_id"
	eq  = "eq"
	ne  = "ne"
	gt  = "gt"
	gte = "gte"
	lt  = "lt"
	lte = "lte"
	in  = "in"
	nin = "nin"
	set = "set"

	asc  = "asc"
	desc = "desc"

	and = "and"
)

var supportedOps = map[string]string{
	eq:  "equal",
	ne:  "not equal",
	gt:  "greater than",
	gte: "greater than equal",
	lt:  "less than",
	lte: "less than equal",
	in:  "in",
	nin: "not in",
}

var sortOrders = map[string]int{
	asc:  1,
	desc: -1,
}
