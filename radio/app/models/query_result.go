package models

import (
	"fmt"
)

type QueryResultEntry struct {
	//	Id int `json:"id"`
	Model   string `json:"model"`
	Type    string `json:"type"`
	Version string `json:"version"`
}

func (e QueryResultEntry) String() string {
	return fmt.Sprintf("QueryResultEntry(model=%s, type=%s, version=%s)",
		e.Model, e.Type, e.Version)
}

type QueryResult struct {
	Count int `json:"count"`
	Data []QueryResultEntry `json:"data"`
	Extra map[string]interface{} `json:"extra"`
}

func (r QueryResult) String() string {
	return fmt.Sprintf("QueryResult(Count=%d, Data=%s, Extra=%s)",
		r.Count, r.Data, r.Extra)
}