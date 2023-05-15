package config

const (
	PathV1Prefix = "/api/v1/"

	PathHealthCheck       = "/"
	PathCreateCategory    = PathV1Prefix + "create_category"
	PathUpdateCategory    = PathV1Prefix + "update_category"
	PathGetCategory       = PathV1Prefix + "get_category"
	PathGetCategories     = PathV1Prefix + "get_categories"
	PathCreateTransaction = PathV1Prefix + "create_transaction"
	PathUpdateTransaction = PathV1Prefix + "update_transaction"
	PathGetTransaction    = PathV1Prefix + "get_transaction"
	PathGetTransactions   = PathV1Prefix + "get_transactions"
	PathGetMonthBudgets         = PathV1Prefix + "get_month_budgets"
	PathGetFullYearBudget = PathV1Prefix + "get_full_year_budget"
	PathSetBudget          = PathV1Prefix + "set_budget"
)

const (
	AmountDecimalPlaces = 2

	OrderAsc  = "asc"
	OrderDesc = "desc"

	DefaultPagingLimit = 100
	MaxPagingLimit     = 500

	MinPagingPage = 1
)
