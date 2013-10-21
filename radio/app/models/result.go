package models

import ()

const (
	API_VERSION = "1.0"
)

type Result struct {
	Data  map[string]interface{} `json:"data"`
	Extra map[string]interface{} `json:"extra"`
}

type QueryResult struct {
	Result
}

func NewQueryResult() *QueryResult {
	result := QueryResult{}
	result.Data = make(map[string]interface{})
	result.Extra = make(map[string]interface{})
	result.Extra["api_version"] = API_VERSION
	return &result
}

type RadioOtaReleaseResult struct {
	Result
}

func NewRadioOtaReleaseResult() *RadioOtaReleaseResult {
	result := RadioOtaReleaseResult{}
	result.Data = make(map[string]interface{})
	result.Extra = make(map[string]interface{})
	result.Extra["api_version"] = API_VERSION
	return &result
}
