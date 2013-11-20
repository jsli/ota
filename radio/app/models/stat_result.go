package models

import ()

type StatSelfResult struct {
	IP      string `json:"ip"`
	Counter int    `json:"counter"`
}

type StatAllResult struct {
	TotalCount int            `json:"total_count"`
	Detail     map[string]int `json:"detail"`
}
