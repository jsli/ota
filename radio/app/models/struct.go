package models

import (
	"errors"
	"fmt"
	"github.com/jsli/ota/radio/app/constant"
	"github.com/jsli/ota/radio/app/utils"
	"github.com/jsli/ota/radio/app/log"
	"github.com/robfig/revel"
)

const (
	SINGLE_CP  = "single_cp.bin"
	SINGLE_DSP = "single_dsp.bin"
	DSDS_CP    = "dsds_cp.bin"
	DSDS_DSP   = "dsds_dsp.bin"
)

var OFFSET_MAP = map[string]int64{
	SINGLE_CP: 0, SINGLE_DSP: 8388608,
	DSDS_CP: 10485760, DSDS_DSP: 18874368,
}

type CpAtomic struct {
	Model         string
	Type          string
	Version       string
	VersionScalar int64
	Flag          int
}

func (ca CpAtomic) String() string {
	return fmt.Sprintf("CpAtomic(Model=%s, Type=%s, version=%s, flag=%d)",
		ca.Model, ca.Type, ca.Version, ca.Flag)
}

func (ca *CpAtomic) Validate(v *revel.Validation, dal *Dal) {
	v.Check(ca.Model,
		revel.Required{},
		LegalModelValidator{},
	).Message("Illegal model")

	v.Check(ca.Type,
		revel.Required{},
		LegalTypeValidator{},
	).Message("Illegal type")

	v.Check(ca.Version,
		revel.Required{},
		revel.Match{constant.CP_VERSION_REX},
		DupCpValidator{dal, ca},
	).Message("Illegal version")
}

func (ca *CpAtomic) Save(dal *Dal) error {
	tag := "CpAtomic.Save"
	stmt, err := dal.Link.Prepare("INSERT cp_source_file SET model=?, type=?, version=?,version_scalar=?,flag=?")
	if err != nil {
		log.Log(tag, fmt.Sprintf("Prepare error : %s", err))
		return err
	}
	_, err = stmt.Exec(ca.Model, ca.Type, ca.Version, ca.VersionScalar, ca.Flag)
	if err != nil {
		log.Log(tag, fmt.Sprintf("Exec error : %s", err))
		return err
	}

	//	id, err := res.LastInsertId()
	return nil
}

type LegalModelValidator struct {
}

func (legal LegalModelValidator) IsSatisfied(obj interface{}) bool {
	return utils.IsAvailableModel(obj.(string))
}

func (legal LegalModelValidator) DefaultMessage() string {
	return fmt.Sprintf("Illegal model, should be one of %s", constant.MODEL_LIST)
}

type LegalTypeValidator struct {
}

func (v LegalTypeValidator) IsSatisfied(obj interface{}) bool {
	return utils.IsAvailableType(obj.(string))
}

func (v LegalTypeValidator) DefaultMessage() string {
	return fmt.Sprintf("Illegal type, should be one of %s", constant.TYPE_LIST)
}

type DupCpValidator struct {
	dal *Dal
	ca  *CpAtomic
}

func (v DupCpValidator) IsSatisfied(obj interface{}) bool {
	tag := "DupCpValidator.IsSatisfied"
	ca, err := FindCpAtomicOne(v.dal, v.ca.Model, v.ca.Type, v.ca.Version)
	log.Log(tag, fmt.Sprintf("%s", err))
	if ca != nil {
		return false
	}
	return true
}

func (v DupCpValidator) DefaultMessage() string {
	return fmt.Sprintf("Duplicated cp file")
}

func FindCpAtomicOne(dal *Dal, model string, _type string, version string) (*CpAtomic, error) {
	var id int = -1
	row := dal.Link.QueryRow(fmt.Sprintf("SELECT id, version_scalar, flag FROM cp_source_file where model='%s' and type='%s' and version='%s'",
		model, _type, version))
	ca := CpAtomic{Model: model, Type: _type, Version: version}
	err := row.Scan(&id, &ca.VersionScalar, &ca.Flag)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("query cp source file failed: %s", err))
	}
	if id >= 0 {
		return &ca, errors.New(fmt.Sprintf("dupliated cp source file : %d", id))
	}
	return nil, errors.New(fmt.Sprintf("query cp source file failed: %s-%s-%s", model, _type, version))
}

func GetOffset(name string) int64 {
	if offset, ok := OFFSET_MAP[name]; ok {
		return offset
	}
	return -1
}
