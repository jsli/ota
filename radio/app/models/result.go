package models

import (
	ota_constant "github.com/jsli/ota/radio/app/constant"
)

type ExtraSetter interface {
	ApiVersionSetter
	ErrorSetter
}

type ApiVersionSetter interface {
	SetApiVersion(version string)
}

type ErrorSetter interface {
	SetErrorCode(code int)
	SetErrorMessage(message interface{})
}

type ExtraInfo struct {
	ApiVersion   string      `json:"api_version"`
	ErrorCode    int         `json:"error_code"`
	ErrorMessage interface{} `json:"error_message"`
}

/*
 * implement ApiVersionInterface
 */
func (eInfo *ExtraInfo) SetApiVersion(version string) {
	eInfo.ApiVersion = version
}

/*
 * implement ErrorInterface
 */
func (eInfo *ExtraInfo) SetErrorCode(code int) {
	eInfo.ErrorCode = code
}

func (eInfo *ExtraInfo) SetErrorMessage(message interface{}) {
	eInfo.ErrorMessage = message
}

type DataSetter interface {
	SetData(data interface{})
}

type Result struct {
	Data  interface{} `json:"data"`
	Extra ExtraSetter `json:"extra"`
}

func (result *Result) SetData(data interface{}) {
	result.Data = data
}

func (result *Result) SetErrorCode(code int) {
	result.Extra.SetErrorCode(code)
}

func (result *Result) SetErrorMessage(message interface{}) {
	result.Extra.SetErrorMessage(message)
}

func NewResult() *Result {
	result := &Result{}
	result.Data = nil
	result.Extra = new(ExtraInfo)
	result.Extra.SetApiVersion(ota_constant.CURRENT_API_VERSION)
	return result
}
