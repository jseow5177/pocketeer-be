package config

const (
	PathHealthCheck = "/"

	// User APIs
	PathV1Prefix                = "/api/v1/"
	PathSignUp                  = PathV1Prefix + "sign_up"
	PathLogin                   = PathV1Prefix + "log_in"
	PathGetUser                 = PathV1Prefix + "get_user"
	PathVerifyEmail             = PathV1Prefix + "verify_email"
	PathInitUser                = PathV1Prefix + "init_user"
	PathSendOTP                 = PathV1Prefix + "send_otp"
	PathUpdateUserMeta          = PathV1Prefix + "update_user_meta"
	PathCreateAccount           = PathV1Prefix + "create_account"
	PathGetAccount              = PathV1Prefix + "get_account"
	PathUpdateAccount           = PathV1Prefix + "update_account"
	PathGetAccounts             = PathV1Prefix + "get_accounts"
	PathDeleteAccount           = PathV1Prefix + "delete_account"
	PathGetAccountsSummary      = PathV1Prefix + "get_accounts_summary"
	PathCreateCategory          = PathV1Prefix + "create_category"
	PathUpdateCategory          = PathV1Prefix + "update_category"
	PathGetCategory             = PathV1Prefix + "get_category"
	PathGetCategoryBudget       = PathV1Prefix + "get_category_budget"
	PathGetCategories           = PathV1Prefix + "get_categories"
	PathGetCategoriesBudget     = PathV1Prefix + "get_categories_budget"
	PathDeleteCategory          = PathV1Prefix + "delete_category"
	PathSumCategoryTransactions = PathV1Prefix + "sum_category_transactions"
	PathCreateTransaction       = PathV1Prefix + "create_transaction"
	PathUpdateTransaction       = PathV1Prefix + "update_transaction"
	PathDeleteTransaction       = PathV1Prefix + "delete_transaction"
	PathGetTransaction          = PathV1Prefix + "get_transaction"
	PathGetTransactions         = PathV1Prefix + "get_transactions"
	PathGetTransactionGroups    = PathV1Prefix + "get_transaction_groups"
	PathSumTransactions         = PathV1Prefix + "sum_transactions"
	PathGetTransactionsSummary  = PathV1Prefix + "get_transactions_summary"
	PathGetBudget               = PathV1Prefix + "get_budget"
	PathUpdateBudget            = PathV1Prefix + "update_budget"
	PathGetBudgets              = PathV1Prefix + "get_budgets"
	PathCreateBudget            = PathV1Prefix + "create_budget"
	PathDeleteBudget            = PathV1Prefix + "delete_budget"
	PathSearchSecurities        = PathV1Prefix + "search_securities"
	PathCreateHolding           = PathV1Prefix + "create_holding"
	PathUpdateHolding           = PathV1Prefix + "update_holding"
	PathGetHolding              = PathV1Prefix + "get_holding"
	PathDeleteHolding           = PathV1Prefix + "delete_holding"
	PathCreateLot               = PathV1Prefix + "create_lot"
	PathDeleteLot               = PathV1Prefix + "delete_lot"
	PathUpdateLot               = PathV1Prefix + "update_lot"
	PathGetLot                  = PathV1Prefix + "get_lot"
	PathGetLots                 = PathV1Prefix + "get_lots"
	PathCreateFeedback          = PathV1Prefix + "create_feedback"
	PathGetExchangeRate         = PathV1Prefix + "get_exchange_rate"
	PathGetCurrencies           = PathV1Prefix + "get_currencies"
	PathGetMetrics              = PathV1Prefix + "get_metrics"

	// Admin APIs
	PathAdminV1Prefix   = "/api/admin/v1/"
	PathAdminSyncQuotes = PathAdminV1Prefix + "sync_quotes"
)

const (
	OrderAsc  = "asc"
	OrderDesc = "desc"

	DefaultPagingLimit = 100
	MaxPagingLimit     = 500

	MinPagingPage = 1

	MaxTransactionNoteLength = 120
	MaxAccountNoteLength     = 60

	PasswordMinLength = 8
	SaltByteSize      = 24

	StandardDP = 2
	PreciseDP  = 5

	LogLevelDebug = "DEBUG"
	LogLevelInfo  = "INFO"

	MinCurrencyDate = "20180101"
)
