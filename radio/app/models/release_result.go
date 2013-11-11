package models

import (
	ota_constant "github.com/jsli/ota/radio/app/constant"
)

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
