package models

import (
	ota_constant "github.com/jsli/ota/radio/app/constant"
)

type MaintainResult struct {
	Result
}

func NewMaintainResult() *MaintainResult {
	result := MaintainResult{}
	result.Extra = new(ExtraInfo)
	result.Extra.SetApiVersion(ota_constant.CURRENT_API_VERSION)
	result.Extra.SetErrorCode(ota_constant.ERROR_CODE_MAINTAIN)
	result.Extra.SetErrorMessage("Maintain")
	return &result
}
