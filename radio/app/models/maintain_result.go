package models

import (
	ota_constant "github.com/jsli/ota/radio/app/constant"
)

type MaintainResult struct {
	Result
}

func NewMaintainResult() *MaintainResult {
	result := MaintainResult{}
	result.Extra = ExtraInfo{ApiVersion: API_VERSION, ErrorCode: ota_constant.ERROR_CODE_MAINTAIN, ErrorMessage: "Maintain"}
	return &result
}
