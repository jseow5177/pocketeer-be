package config

const (
	PathV1Prefix = "/api/v1/"

	PathHealthCheck = "/"

	PathCreateCategory = PathV1Prefix + "create_category"
	PathUpdateCategory = PathV1Prefix + "update_category"
	PathGetCategory    = PathV1Prefix + "get_category"
	PathGetCategories  = PathV1Prefix + "get_categories"

	PathCreateTransaction = PathV1Prefix + "create_transaction"
	PathUpdateTransaction = PathV1Prefix + "update_transaction"
	PathGetTransaction    = PathV1Prefix + "get_transaction"
	PathGetTransactions   = PathV1Prefix + "get_transactions"
	PathAggrTransactions  = PathV1Prefix + "aggr_transactions"

	PathGetCategoryBudgetsByMonth = PathV1Prefix + "get_category_budgets_by_month"
	PathGetAnnualBudgetBreakdown  = PathV1Prefix + "get_annual_budget_breakdown"
	PathSetBudget                 = PathV1Prefix + "set_budget"
)

const (
	OrderAsc  = "asc"
	OrderDesc = "desc"

	DefaultPagingLimit = 100
	MaxPagingLimit     = 500

	MinPagingPage = 1

	MaxTransactionNoteLength = 120

	AmountDecimalPlaces = 2
)
