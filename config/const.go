package config

const (
	PathV1Prefix = "/api/v1/"

	PathHealthCheck       = "/"
	PathCreateCategory    = PathV1Prefix + "create_category"
	PathUpdateCategory    = PathV1Prefix + "update_category"
	PathGetCategory       = PathV1Prefix + "get_category"
	PathGetCategories     = PathV1Prefix + "get_categories"
	PathCreateTransaction = PathV1Prefix + "create_transaction"
	PathGetTransaction    = PathV1Prefix + "get_transaction"
)
