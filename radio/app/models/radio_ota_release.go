package models

import (
	"database/sql"
	"fmt"
	"github.com/jsli/ota/radio/app/constant"
)

type RadioOtaRelease struct {
	Id      int64
	Mode    string
	Version string
	Md5     string
	Size    int64
	Flag    int
}

func (ror RadioOtaRelease) String() string {
	return fmt.Sprintf("RadioRelease(Id=%d, Mode=%s, Version=%s, Md5=%s, Size=%d, Flag=%d)",
		ror.Id, ror.Mode, ror.Version, ror.Md5, ror.Size, ror.Flag)
}

func (ror *RadioOtaRelease) Save(dal *Dal) (int64, error) {
	insert_sql := fmt.Sprintf("INSERT %s SET mode=?, version=?, md5=?, size=?, flag=?", constant.TABLE_RADIO_OTA_RELEASE)
	stmt, eror := dal.Link.Prepare(insert_sql)

	if eror != nil {
		return -1, eror
	}
	res, eror := stmt.Exec(ror.Mode, ror.Version, ror.Md5, ror.Size, ror.Flag)
	if eror != nil {
		return -1, eror
	}

	id, eror := res.LastInsertId()
	return id, eror
}

func FindRadioOtaRelease(dal *Dal, query string) (*RadioOtaRelease, error) {
	row := dal.Link.QueryRow(query)
	ror := RadioOtaRelease{}
	err := row.Scan(&ror.Id, &ror.Mode, &ror.Version, &ror.Md5, &ror.Size, &ror.Flag)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &ror, nil
}
