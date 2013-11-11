package models

import ()

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
