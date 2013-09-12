package policy

import (
	"fmt"
	cp_constant "github.com/jsli/cp_release/constant"
	cp_policy "github.com/jsli/cp_release/policy"
	"github.com/jsli/cp_release/release"
	ota_constant "github.com/jsli/ota/radio/app/constant"
	"github.com/jsli/ota/radio/app/models"
	"strings"
)

func ProvideRadioRelease(dal *models.Dal, parsedParams *ParsedParams, result *models.RadioOtaReleaseResult, versions string) error {
	var mode, version string
	version_list := strings.Split(versions, "-")

	single_info, ok := parsedParams.CpMap[ota_constant.TYPE_SINGLE]
	if ok {
		arbi := single_info.ImageMap[ota_constant.ID_ARBI]
		mode = arbi.Mode
		version = version_list[0]
	}

	dsds_info, ok := parsedParams.CpMap[ota_constant.TYPE_DSDS]
	if ok {
		arbi := dsds_info.ImageMap[ota_constant.ID_ARB2]
		mode = fmt.Sprintf("%s-%s", mode, arbi.Mode)
		version = fmt.Sprintf("%s-%s", version, version_list[1])
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE mode='%s' AND version='%s' AND flag=%d LIMIT 1",
		ota_constant.TABLE_RADIO_OTA_RELEASE, mode, version, cp_constant.AVAILABLE_FLAG)
	release, err := models.FindRadioOtaRelease(dal, query)
	if err != nil {
		return err
	}

	if release != nil {
		result.Data["url"] = fmt.Sprintf("%s/%s/%s", mode, version, ota_constant.RADIO_OTA_PACKAGE_NAME)
		result.Data["md5"] = release.Md5
		result.Data["size"] = release.Size
		result.Extra["error"] = nil
	}
	return nil
}

func ProvideQueryData(dal *release.Dal, parsedParams *ParsedParams, result *models.QueryResult) error {
	switch parsedParams.Type {
	case ota_constant.TYPE_DSDS, ota_constant.TYPE_DSDS_RFIC:
		data, err := generateQueryData(dal, parsedParams.CpMap[parsedParams.Type], parsedParams.Type)
		if err != nil {
			return err
		}
		result.Data[ota_constant.TYPE_DSDS] = data
		fallthrough
	case ota_constant.TYPE_SINGLE, ota_constant.TYPE_SINGLE_RFIC:
		data, err := generateQueryData(dal, parsedParams.CpMap[parsedParams.Type], parsedParams.Type)
		if err != nil {
			return err
		}
		result.Data[ota_constant.TYPE_SINGLE] = data
	}
	return nil
}

func generateQueryData(dal *release.Dal, cp_info *CpInfo, _type string) (map[string]map[string][]string, error) {
	cp_list, err := getCpVersionList(dal, cp_info)
	if err != nil {
		return nil, err
	}
	if cp_list != nil && len(cp_list) > 0 {
		data := make(map[string]map[string][]string)
		mode := cp_info.Mode
		sim := cp_info.Sim
		for _, cp_version := range cp_list {
			image_list, err := getImageList(dal, mode, sim, cp_version, _type)
			if err != nil {
				return nil, err
			}
			data[cp_version] = image_list
		}
		return data, nil
	}
	return nil, nil
}

func getImageList(dal *release.Dal, mode string, sim string, cp_version string, _type string) (map[string][]string, error) {
	data := make(map[string][]string)

	arbi_list, err := getArbiList(dal, mode, sim, cp_version, _type)
	if err != nil {
		return nil, err
	}

	grbi_list, err := getGrbiList(dal, mode, sim, cp_version, _type)
	if err != nil {
		return nil, err
	}

	rfic_list, err := getRficList(dal, mode, sim, cp_version, _type)
	if err != nil {
		return nil, err
	}

	switch _type {
	case ota_constant.TYPE_SINGLE:
		if arbi_list != nil && len(arbi_list) > 0 {
			data[ota_constant.ID_ARBI] = arbi_list
		}
		if grbi_list != nil && len(grbi_list) > 0 {
			data[ota_constant.ID_GRBI] = grbi_list
		}
	case ota_constant.TYPE_DSDS:
		if arbi_list != nil && len(arbi_list) > 0 {
			data[ota_constant.ID_ARB2] = arbi_list
		}
		if grbi_list != nil && len(grbi_list) > 0 {
			data[ota_constant.ID_GRB2] = grbi_list
		}
	case ota_constant.TYPE_SINGLE_RFIC:
		if arbi_list != nil && len(arbi_list) > 0 {
			data[ota_constant.ID_ARBI] = arbi_list
		}
		if grbi_list != nil && len(grbi_list) > 0 {
			data[ota_constant.ID_GRBI] = grbi_list
		}
		if rfic_list != nil && len(rfic_list) > 0 {
			data[ota_constant.ID_RFIC] = rfic_list
		}
	case ota_constant.TYPE_DSDS_RFIC:
		if arbi_list != nil && len(arbi_list) > 0 {
			data[ota_constant.ID_ARB2] = arbi_list
		}
		if grbi_list != nil && len(grbi_list) > 0 {
			data[ota_constant.ID_GRB2] = grbi_list
		}
		if rfic_list != nil && len(rfic_list) > 0 {
			data[ota_constant.ID_RFI2] = rfic_list
		}
	default:
		return nil, fmt.Errorf("Unknow type : %s", _type)
	}

	return data, nil
}

func getArbiList(dal *release.Dal, mode string, sim string, cp_version string, _type string) ([]string, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE mode='%s' AND sim='%s' AND version='%s' AND flag=%d",
		cp_constant.TABLE_CP, mode, sim, cp_version, cp_constant.AVAILABLE_FLAG)
	cp, err := release.FindCpRelease(dal, query)
	if err != nil {
		return nil, err
	}

	arbi_list := make([]string, 0, 10)
	query = fmt.Sprintf("SELECT * FROM %s where cp_id=%d AND flag=%d", cp_constant.TABLE_ARBI, cp.Id, cp_constant.AVAILABLE_FLAG)
	arbis, err := release.FindArbiList(dal, query)
	if err != nil {
		return nil, err
	}

	if len(arbis) == 0 {
		return nil, nil
	}

	for _, arbi := range arbis {
		rel_path := arbi.RelPath
		if !strings.Contains(rel_path, "/RFIC/") {
			arbi_list = append(arbi_list, arbi.RelPath)
		}
	}
	return arbi_list, nil
}

func getGrbiList(dal *release.Dal, mode string, sim string, cp_version string, _type string) ([]string, error) {
	max_version := cp_policy.QuantitateVersion(cp_version)
	min_version := (max_version / 1000000) * 1000000
	query := fmt.Sprintf("SELECT * FROM %s WHERE mode='%s' AND sim='%s' AND (version_scalar<=%d AND version_scalar>=%d ) AND flag=%d ORDER BY version_scalar DESC",
		cp_constant.TABLE_CP, mode, sim, max_version, min_version, cp_constant.AVAILABLE_FLAG)
	//	fmt.Println("cp query : " + query)
	cps, err := release.FindCpReleaseList(dal, query)
	if err != nil {
		return nil, err
	}

	for _, cp := range cps {
		grbi_list := make([]string, 0, 10)
		query = fmt.Sprintf("SELECT * FROM %s where cp_id=%d AND flag=%d", cp_constant.TABLE_GRBI, cp.Id, cp_constant.AVAILABLE_FLAG)
		grbis, err := release.FindGrbiList(dal, query)
		if err != nil {
			return nil, err
		}

		if len(grbis) == 0 {
			continue
		}

		for _, grbi := range grbis {
			grbi_list = append(grbi_list, grbi.RelPath)
		}
		return grbi_list, nil
	}

	return nil, nil
}

func getRficList(dal *release.Dal, mode string, sim string, cp_version string, _type string) ([]string, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE mode='%s' AND sim='%s' AND version='%s' AND flag=%d",
		cp_constant.TABLE_CP, mode, sim, cp_version, cp_constant.AVAILABLE_FLAG)
	cp, err := release.FindCpRelease(dal, query)
	if err != nil {
		return nil, err
	}

	arbi_list := make([]string, 0, 10)
	query = fmt.Sprintf("SELECT * FROM %s where cp_id=%d AND flag=%d", cp_constant.TABLE_ARBI, cp.Id, cp_constant.AVAILABLE_FLAG)
	arbis, err := release.FindArbiList(dal, query)
	if err != nil {
		return nil, err
	}

	if len(arbis) == 0 {
		return nil, nil
	}

	for _, arbi := range arbis {
		rel_path := arbi.RelPath
		if strings.Contains(rel_path, "/RFIC/") {
			arbi_list = append(arbi_list, arbi.RelPath)
		}
	}
	return arbi_list, nil
}

func getCpVersionList(dal *release.Dal, cp_info *CpInfo) ([]string, error) {
	cp_list := make([]string, 0, 10)
	mode := cp_info.Mode
	sim := cp_info.Sim
	version := cp_info.Version
	version_scalar := cp_policy.QuantitateVersion(version)

	query := fmt.Sprintf("SELECT * FROM %s WHERE mode='%s' AND sim='%s' AND flag=%d AND version_scalar > %d ORDER BY version_scalar DESC",
		cp_constant.TABLE_CP, mode, sim, cp_constant.AVAILABLE_FLAG, version_scalar)
	list, err := doGetCpVersionList(dal, query)
	if err != nil {
		return nil, err
	}
	cp_list = append(cp_list, list...)

	query = fmt.Sprintf("SELECT * FROM %s WHERE mode='%s' AND sim='%s' AND flag=%d AND version_scalar < %d ORDER BY version_scalar DESC LIMIT 5",
		cp_constant.TABLE_CP, mode, sim, cp_constant.AVAILABLE_FLAG, version_scalar)
	list, err = doGetCpVersionList(dal, query)
	if err != nil {
		return nil, err
	}
	cp_list = append(cp_list, list...)
	if len(cp_list) == 0 {
		return nil, nil
	}

	return cp_list, nil
}

func doGetCpVersionList(dal *release.Dal, query string) ([]string, error) {
	version_list := make([]string, 0, 10)
	cps, err := release.FindCpReleaseList(dal, query)
	if err != nil {
		return nil, err
	}
	for _, cp := range cps {
		version_list = append(version_list, cp.Version)
	}

	if len(version_list) > 0 {
		return version_list, nil
	} else {
		return nil, nil
	}
}
