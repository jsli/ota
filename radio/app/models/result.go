package models

import (
	ota_constant "github.com/jsli/ota/radio/app/constant"
)

const (
	API_VERSION = "1.0"
)

type ExtraInfo struct {
	ApiVersion   string      `json:"api_version"`
	ErrorCode    int         `json:"error_code"`
	ErrorMessage interface{} `json:"error_message"`
}

type Result struct {
	Extra ExtraInfo `json:"extra"`
}

type QueryData struct {
	Available map[string]interface{} `json:"available"`
	Current   map[string]interface{} `json:"current"`
}

type QueryResult struct {
	Data QueryData `json:"data"`
	Result
}

func NewQueryResult() *QueryResult {
	result := QueryResult{}
	result.Data = QueryData{}
	result.Extra = ExtraInfo{ApiVersion: API_VERSION, ErrorCode: ota_constant.ERROR_CODE_NOERR, ErrorMessage: nil}
	return &result
}

type RadioOtaReleaseData struct {
	Url         string `json:"url"`
	Md5         string `json:"md5"`
	Size        int64  `json:"size"`
	CreatedTime string `json:"created_time"`
}

type RadioOtaReleaseResult struct {
	Data RadioOtaReleaseData `json:"data"`
	Result
}

func NewRadioOtaReleaseResult() *RadioOtaReleaseResult {
	result := RadioOtaReleaseResult{}
	result.Data = RadioOtaReleaseData{}
	result.Extra = ExtraInfo{ApiVersion: API_VERSION, ErrorCode: ota_constant.ERROR_CODE_NOERR, ErrorMessage: nil}
	return &result
}

type MaintainResult struct {
	Result
}

func NewMaintainResult() *MaintainResult {
	result := MaintainResult{}
	result.Extra = ExtraInfo{ApiVersion: API_VERSION, ErrorCode: ota_constant.ERROR_CODE_NOERR, ErrorMessage: "Maintain"}
	return &result
}
