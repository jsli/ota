package models

import ()

type ReleaseResultData struct {
	Url         string `json:"url"`
	Md5         string `json:"md5"`
	Size        int64  `json:"size"`
	CreatedTime string `json:"created_time"`
}
