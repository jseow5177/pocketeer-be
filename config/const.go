package config

const (
	PathV1Prefix = "/api/v1/"

	PathHealthCheck = "/"

	PathSignUp  = PathV1Prefix + "sign_up"
	PathLogin   = PathV1Prefix + "log_in"
	PathGetUser = PathV1Prefix + "get_user"

	PathCreateAccount = PathV1Prefix + "create_account"
	PathGetAccount    = PathV1Prefix + "get_account"

	PathCreateCategory = PathV1Prefix + "create_category"
	PathUpdateCategory = PathV1Prefix + "update_category"
	PathGetCategory    = PathV1Prefix + "get_category"
	PathGetCategories  = PathV1Prefix + "get_categories"

	PathCreateTransaction = PathV1Prefix + "create_transaction"
	PathUpdateTransaction = PathV1Prefix + "update_transaction"
	PathGetTransaction    = PathV1Prefix + "get_transaction"
	PathGetTransactions   = PathV1Prefix + "get_transactions"
	PathAggrTransactions  = PathV1Prefix + "aggr_transactions"

	PathGetBudget  = PathV1Prefix + "get_budget"
	PathGetBudgets = PathV1Prefix + "get_budgets"
	PathSetBudget  = PathV1Prefix + "set_budget"
)

const (
	OrderAsc  = "asc"
	OrderDesc = "desc"

	DefaultPagingLimit = 100
	MaxPagingLimit     = 500

	MinPagingPage = 1

	MaxTransactionNoteLength = 120
	MaxAccountNoteLength     = 60

	AmountDecimalPlaces = 2

	PasswordMinLength = 8
	UsernameMaxLength = 60
	SaltByteSize      = 24
)
