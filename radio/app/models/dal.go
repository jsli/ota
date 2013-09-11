package models

import (
	"database/sql"
	"fmt"
	"github.com/jsli/ota/radio/app/log"
)

const (
	DNS    = "root:lijinsong@/ota?charset=utf8"
	DRIVER = "mysql"
)

type Dal struct {
	Link *sql.DB
}

func NewDal(driver string, dns string) (*Dal, error) {
	tag := "NewDal"
	db, err := sql.Open(driver, dns)
	if err != nil {
		log.Log(tag, fmt.Sprintf("Open %s : %s error : %s", driver, dns, err))
		return nil, err
	}
	return &Dal{db}, nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (dal *Dal) Close() {
	dal.Link.Close()
}
