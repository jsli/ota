package models

import (
	"fmt"
)

type Product struct {
	Id          int
	Model       string
	Platform    string
	Vendor      string
	Description string
}

func FindProduct(dal *Dal, model string) (*Product, error) {
	rows := dal.Link.QueryRow(fmt.Sprintf("SELECT id, model, platform, vendor FROM products where model = '%s'", model))
	p := Product{}
	err := rows.Scan(&p.Id, &p.Model, &p.Platform, &p.Vendor)
	return &p, err
}
