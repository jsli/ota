package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

const (
	DNS    = "root:lijinsong@/ota?charset=utf8"
	DRIVER = "mysql"
)

type Dal struct {
	DB *sql.DB
}

func NewDal() (*Dal, error) {
	db, err := sql.Open(DRIVER, DNS)
	if err != nil {
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
	dal.DB.Close()
}
