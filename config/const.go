package config

const (
	PathV1Prefix = "/api/v1"

	PathHealthCheck    = "/"
	PathCreateCategory = PathV1Prefix + "/category/create"
	PathUpdateCategory = PathV1Prefix + "/category/update"
	PathGetCategory    = PathV1Prefix + "/category/get"
	PathGetCategories  = PathV1Prefix + "/category/list"
)
