package models

import ()

const (
	API_VERSION = "1.0"
)

type ErrorInterface interface {
	SetErrorCode(code int)
	SetErrorMessage(message interface{})
}

type ExtraInfo struct {
	ApiVersion   string      `json:"api_version"`
	ErrorCode    int         `json:"error_code"`
	ErrorMessage interface{} `json:"error_message"`
}

type Result struct {
	Extra ExtraInfo `json:"extra"`
}

func (result *Result) SetErrorCode(code int) {
	result.Extra.ErrorCode = code
}

func (result *Result) SetErrorMessage(message interface{}) {
	result.Extra.ErrorMessage = message
}
