package config

const (
	pathPrefix = "/api/v1"

	PathHealthCheck        = "/"
	PathCreateCategory     = pathPrefix + "/category/create"
	PathUpdateCategory     = pathPrefix + "/category/update"
	PathGetCategory        = pathPrefix + "/category/get"
	PathGetCategories      = pathPrefix + "/category/list"
	PathGetBudgets         = pathPrefix + "/budget/list"
	PathGetBudgetBreakdown = pathPrefix + "/budget/get_breakdown"
	PathSetBudget          = pathPrefix + "/budget/set"
)
