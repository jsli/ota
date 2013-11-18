package models

import (
	"database/sql"
	"fmt"
	ota_constant "github.com/jsli/ota/radio/app/constant"
)

type RadioOtaRelease struct {
	Id          int64  `json:"id"`
	Model       string `json:"model"`
	Platform    string `json:"platform"`
	FingerPrint string `json:"fingerprint"`
	Md5         string `json:"md5sum"`
	Size        int64  `json:"size"`
	Flag        int    `json:"flag"`
	ReleaseNote string `json:"release_note"`
	ModifiedTs  int64  `json:"m_ts"`
	CreatedTs   int64  `json:"c_ts"`
	Versions    string `json:"versions"`
	Images      string `json:"image_list"` //divide by space, like [xxxx yyyy zzzz]
}

func (ror RadioOtaRelease) String() string {
	return fmt.Sprintf("RadioRelease(Id=%d, Model=%s, Platform=%s, FingerPrint=%s, Md5=%s, Size=%d, Flag=%d, ReleaseNote=%s, MT=%d, CT=%d, Version=%s, Images=%s)",
		ror.Id, ror.Model, ror.Platform, ror.FingerPrint, ror.Md5, ror.Size, ror.Flag, ror.ReleaseNote, ror.ModifiedTs, ror.CreatedTs, ror.Versions, ror.Images)
}

func (ror *RadioOtaRelease) Save(dal *Dal) (int64, error) {
	insert_sql := fmt.Sprintf("INSERT %s SET model=?, platform=?, fingerprint=?, md5=?, size=?, flag=?, release_note=?, modified_ts=?, created_ts=?, versions=?, images=?",
		ota_constant.TABLE_RADIO_OTA_RELEASE)
	stmt, eror := dal.DB.Prepare(insert_sql)

	if eror != nil {
		return -1, eror
	}
	res, eror := stmt.Exec(ror.Model, ror.Platform, ror.FingerPrint, ror.Md5, ror.Size, ror.Flag, ror.ReleaseNote, ror.ModifiedTs, ror.CreatedTs, ror.Versions, ror.Images)
	if eror != nil {
		return -1, eror
	}

	id, eror := res.LastInsertId()
	return id, eror
}

func (ror *RadioOtaRelease) Update(dal *Dal) (int64, error) {
	insert_sql := fmt.Sprintf("UPDATE %s SET model=?, platform=?, md5=?, size=?, flag=?, release_note=?, modified_ts=?, created_ts=?, versions=?, images=?",
		ota_constant.TABLE_RADIO_OTA_RELEASE)
	stmt, eror := dal.DB.Prepare(insert_sql)

	if eror != nil {
		return -1, eror
	}
	res, eror := stmt.Exec(ror.Model, ror.Platform, ror.Md5, ror.Size, ror.Flag, ror.ReleaseNote, ror.ModifiedTs, ror.CreatedTs, ror.Versions, ror.Images)
	if eror != nil {
		return -1, eror
	}

	id, eror := res.LastInsertId()
	return id, eror
}

func (ror *RadioOtaRelease) Delete(dal *Dal) (int64, error) {
	return DeleteRadioReleaseByFp(dal, ror.FingerPrint)
}

func FindRadioOtaReleaseByFp(dal *Dal, fp string) (*RadioOtaRelease, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE fingerprint='%s' LIMIT 1",
		ota_constant.TABLE_RADIO_OTA_RELEASE, fp)
	return FindRadioOtaRelease(dal, query)
}

func FindRadioOtaReleaseList(dal *Dal, flag int) ([]*RadioOtaRelease, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE flag=%d",
		ota_constant.TABLE_RADIO_OTA_RELEASE, flag)
	rows, err := dal.DB.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	rors := make([]*RadioOtaRelease, 0, 100)
	for rows.Next() {
		ror := RadioOtaRelease{}
		err := rows.Scan(&ror.Id, &ror.Model, &ror.Platform, &ror.FingerPrint, &ror.Md5, &ror.Size,
			&ror.Flag, &ror.ReleaseNote, &ror.ModifiedTs, &ror.CreatedTs, &ror.Versions, &ror.Images)

		if err != nil || ror.Id < 0 {
			continue
		}
		rors = append(rors, &ror)
	}
	return rors, nil
}

func FindRadioOtaRelease(dal *Dal, query string) (*RadioOtaRelease, error) {
	row := dal.DB.QueryRow(query)
	ror := RadioOtaRelease{}
	err := row.Scan(&ror.Id, &ror.Model, &ror.Platform, &ror.FingerPrint, &ror.Md5, &ror.Size,
		&ror.Flag, &ror.ReleaseNote, &ror.ModifiedTs, &ror.CreatedTs, &ror.Versions, &ror.Images)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &ror, nil
}

func DeleteRadioReleaseByFp(dal *Dal, fingerprint string) (int64, error) {
	delete_sql := fmt.Sprintf("DELETE FROM %s where fingerprint='%s'", ota_constant.TABLE_RADIO_OTA_RELEASE, fingerprint)
	return DeleteRadioRelease(dal, delete_sql)
}

func DeleteRadioRelease(dal *Dal, delete_sql string) (int64, error) {
	stmt, err := dal.DB.Prepare(delete_sql)

	if err != nil {
		return -1, err
	}
	res, err := stmt.Exec()
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	return id, err
}

type ReleaseCreationTask struct {
	Id            int64
	ReleaseId     int64
	Flag          int
	RetryCount    int
	UpdateRequest string
	Model         string
	Platform      string
	FingerPrint   string
	Data          []byte
	ModifiedTs    int64
	CreatedTs     int64
}

func (rct ReleaseCreationTask) String() string {
	return fmt.Sprintf("ReleaseCreationTask(Id=%d, ReleaseId=%d, UpdateRequest=%s, Flag=%d, RetryCount=%d, Model=%s, Platform=%s, FingerPrint=%s, Data=%s, MT=%d, CT=%d)",
		rct.Id, rct.ReleaseId, rct.UpdateRequest, rct.Flag, rct.RetryCount, rct.Model, rct.Platform, rct.FingerPrint, rct.Data, rct.ModifiedTs, rct.CreatedTs)
}

func (rct *ReleaseCreationTask) Save(dal *Dal) (int64, error) {
	insert_sql := fmt.Sprintf("INSERT %s SET release_id=?, flag=?, retry_count=?, update_request=?, model=?, platform=?, fingerprint=?, binary_data=?, modified_ts=?, created_ts=?",
		ota_constant.TABLE_RELEASE_CREATION_TASK)
	stmt, err := dal.DB.Prepare(insert_sql)
	if err != nil {
		return -1, err
	}
	res, err := stmt.Exec(rct.ReleaseId, rct.Flag, rct.RetryCount, rct.UpdateRequest, rct.Model, rct.Platform, rct.FingerPrint, rct.Data, rct.ModifiedTs, rct.CreatedTs)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	return id, err
}

func (rct *ReleaseCreationTask) Update(dal *Dal) (int64, error) {
	update_sql := fmt.Sprintf("UPDATE %s SET release_id=?, flag=?, retry_count=? WHERE fingerprint=?",
		ota_constant.TABLE_RELEASE_CREATION_TASK)
	stmt, err := dal.DB.Prepare(update_sql)

	if err != nil {
		return -1, err
	}
	res, err := stmt.Exec(rct.ReleaseId, rct.Flag, rct.RetryCount, rct.FingerPrint)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	return id, err
}

func (rct *ReleaseCreationTask) Delete(dal *Dal) (int64, error) {
	return DeleteReleaseCreationTaskFp(dal, rct.FingerPrint)
}

func PopOneCreationTask(dal *Dal) (*ReleaseCreationTask, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE flag=%d ORDER BY id ASC LIMIT 1",
		ota_constant.TABLE_RELEASE_CREATION_TASK, ota_constant.FLAG_INIT)
	return FindReleaseCreationTask(dal, query)
}

func FindReleaseCreationTaskByFp(dal *Dal, fingerprint string) (*ReleaseCreationTask, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE fingerprint='%s' LIMIT 1",
		ota_constant.TABLE_RELEASE_CREATION_TASK, fingerprint)
	return FindReleaseCreationTask(dal, query)
}

func FindReleaseCreationTask(dal *Dal, query string) (*ReleaseCreationTask, error) {
	row := dal.DB.QueryRow(query)
	rct := ReleaseCreationTask{}
	err := row.Scan(&rct.Id, &rct.ReleaseId, &rct.Flag, &rct.RetryCount, &rct.UpdateRequest, &rct.Model, &rct.Platform, &rct.FingerPrint, &rct.Data, &rct.ModifiedTs, &rct.CreatedTs)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &rct, nil
}

func DeleteReleaseCreationTaskFp(dal *Dal, fingerprint string) (int64, error) {
	delete_sql := fmt.Sprintf("DELETE FROM %s where fingerprint='%s'", ota_constant.TABLE_RELEASE_CREATION_TASK, fingerprint)
	return DeleteReleaseCreationTask(dal, delete_sql)
}

func DeleteReleaseCreationTask(dal *Dal, delete_sql string) (int64, error) {
	stmt, err := dal.DB.Prepare(delete_sql)

	if err != nil {
		return -1, err
	}
	res, err := stmt.Exec()
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	return id, err
}
