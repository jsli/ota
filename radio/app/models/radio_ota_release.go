package models

import (
	"database/sql"
	"fmt"
	ota_constant "github.com/jsli/ota/radio/app/constant"
)

type RadioOtaRelease struct {
	Id          int64
	FingerPrint string
	Md5         string
	Size        int64
	Flag        int
	Detail      string
}

func (ror RadioOtaRelease) String() string {
	return fmt.Sprintf("RadioRelease(Id=%d, FingerPrint=%s, Md5=%s, Size=%d, Flag=%d, detail=%s)",
		ror.Id, ror.FingerPrint, ror.Md5, ror.Size, ror.Flag, ror.Detail)
}

func (ror *RadioOtaRelease) Save(dal *Dal) (int64, error) {
	insert_sql := fmt.Sprintf("INSERT %s SET fingerprint=?, md5=?, size=?, flag=?, detail=?", ota_constant.TABLE_RADIO_OTA_RELEASE)
	stmt, eror := dal.DB.Prepare(insert_sql)

	if eror != nil {
		return -1, eror
	}
	res, eror := stmt.Exec(ror.FingerPrint, ror.Md5, ror.Size, ror.Flag, ror.Detail)
	if eror != nil {
		return -1, eror
	}

	id, eror := res.LastInsertId()
	return id, eror
}

func FindRadioOtaRelease(dal *Dal, query string) (*RadioOtaRelease, error) {
	row := dal.DB.QueryRow(query)
	ror := RadioOtaRelease{}
	err := row.Scan(&ror.Id, &ror.FingerPrint, &ror.Md5, &ror.Size, &ror.Flag, &ror.Detail)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &ror, nil
}

type ReleaseCreationTask struct {
	Id            int64
	ReleaseId     int64
	Flag          int
	UpdateRequest string
	FingerPrint   string
	Data          []byte
}

func (rct ReleaseCreationTask) String() string {
	return fmt.Sprintf("ReleaseCreationTask(Id=%d, ReleaseId=%d, UpdateRequest=%s, Flag=%d, FingerPrint=%s, Data=%s)",
		rct.Id, rct.ReleaseId, rct.UpdateRequest, rct.Flag, rct.FingerPrint, rct.Data)
}

func (rct *ReleaseCreationTask) Save(dal *Dal) (int64, error) {
	insert_sql := fmt.Sprintf("INSERT %s SET release_id=?, flag=?, update_request=?, finger_print=?, binary_data=?", ota_constant.TABLE_RELEASE_CREATION_TASK)
	stmt, eror := dal.DB.Prepare(insert_sql)

	if eror != nil {
		return -1, eror
	}
	res, eror := stmt.Exec(rct.ReleaseId, rct.Flag, rct.UpdateRequest, rct.FingerPrint, rct.Data)
	if eror != nil {
		return -1, eror
	}

	id, eror := res.LastInsertId()
	return id, eror
}

func (rct *ReleaseCreationTask) Update(dal *Dal) (int64, error) {
	update_sql := fmt.Sprintf("UPDATE %s SET release_id=?, flag=?, update_request=?, finger_print=?, binary_data=?", ota_constant.TABLE_RELEASE_CREATION_TASK)
	stmt, eror := dal.DB.Prepare(update_sql)

	if eror != nil {
		return -1, eror
	}
	res, eror := stmt.Exec(rct.ReleaseId, rct.Flag, rct.UpdateRequest, rct.FingerPrint, rct.Data)
	if eror != nil {
		return -1, eror
	}

	id, eror := res.LastInsertId()
	return id, eror
}

func PopOneCreationTask(dal *Dal) (*ReleaseCreationTask, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE flag=%d ORDER BY id ASC LIMIT 1",
		ota_constant.TABLE_RELEASE_CREATION_TASK, ota_constant.FLAG_INIT)
	return FindReleaseCreationTask(dal, query)
}

func FindReleaseCreationTask(dal *Dal, query string) (*ReleaseCreationTask, error) {
	row := dal.DB.QueryRow(query)
	rct := ReleaseCreationTask{}
	err := row.Scan(&rct.Id, &rct.ReleaseId, &rct.Flag, &rct.UpdateRequest, &rct.FingerPrint, &rct.Data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &rct, nil
}
