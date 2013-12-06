package models

import (
//	ota_constant "github.com/jsli/ota/radio/app/constant"
)

type API map[string]string

type ApiResult struct {
	Radio API `json:"radio"`
}

func NewApiResult() *ApiResult {
	result := ApiResult{}
	return &result
}
