package policy

import (
	"fmt"

	cp_constant "github.com/jsli/cp_release/constant"
	cp_policy "github.com/jsli/cp_release/policy"
	"github.com/jsli/cp_release/release"
	ota_constant "github.com/jsli/ota/radio/app/constant"
	"github.com/jsli/ota/radio/app/models"
	"github.com/robfig/revel"
)

type ContentProvider interface {
	ProvideRadioRelease(dal *models.Dal, dtim_info *DtimInfo, fp string) (*models.RadioOtaRelease, error)
	ProvideQueryData(dal *release.Dal, dtim_info *DtimInfo, result models.DataSetter) error
}

type CommonContentProvider struct {
}

func (ccp *CommonContentProvider) ProvideRadioRelease(dal *models.Dal, dtim_info *DtimInfo, fp string) (*models.RadioOtaRelease, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE fingerprint='%s' AND flag=%d LIMIT 1",
		ota_constant.TABLE_RADIO_OTA_RELEASE, fp, ota_constant.FLAG_AVAILABLE)
	revel.INFO.Println("query radio ota release: ", query)
	release, err := models.FindRadioOtaRelease(dal, query)
	if err != nil {
		return nil, err
	}

	if release != nil {
		return release, nil
	}
	return nil, nil
}

func getCurrentCpImages(cp_info *CpInfo) map[string]string {
	data := make(map[string]string)
	for key, value := range cp_info.ImageMap {
		data[key] = value.Path
	}
	return data
}

/*
 *    [version]     [image_id]  [image_list]
 * map[string]   map[string]    []string
 */
func getCpAndImages(dal *release.Dal, cp_info *CpInfo, hasRFIC bool) (map[string]map[string][]string, error) {
	cp_list, err := getCpList(dal, cp_info)
	if err != nil {
		return nil, err
	}

	data := make(map[string]map[string][]string)
	for _, cp := range cp_list {
		images_list, err := getImagesByCp(dal, cp, cp_info, hasRFIC)
		if err != nil {
			continue
		}
		data[cp.Version] = images_list
	}

	return data, nil
}

func getImagesByCp(dal *release.Dal, cp *release.CpRelease, cp_info *CpInfo, hasRFIC bool) (map[string][]string, error) {
	data := make(map[string][]string)
	arbi_list, err := getArbiList(dal, cp, cp_info.ImageMap[ota_constant.KEY_ARBEL].Path)
	if err != nil {
		return nil, err
	}
	if arbi_list == nil || len(arbi_list) == 0 {
		return nil, fmt.Errorf(ota_constant.ERROR_MSG_NO_AVAILABLE_IMAGE, ota_constant.KEY_ARBEL)
	}
	data[ota_constant.KEY_ARBEL] = arbi_list

	grbi_list, err := getGrbiList(dal, cp, cp_info.ImageMap[ota_constant.KEY_MSA].Path)
	if err != nil {
		return nil, err
	}
	if grbi_list == nil || len(grbi_list) == 0 {
		return nil, fmt.Errorf(ota_constant.ERROR_MSG_NO_AVAILABLE_IMAGE, ota_constant.KEY_MSA)
	}
	data[ota_constant.KEY_MSA] = grbi_list

	if hasRFIC {
		rfic_list, err := getRficList(dal, cp, cp_info.ImageMap[ota_constant.KEY_RFIC].Path)
		if err != nil {
			return nil, err
		}
		if rfic_list == nil || len(rfic_list) == 0 {
			return nil, fmt.Errorf(ota_constant.ERROR_MSG_NO_AVAILABLE_IMAGE, ota_constant.KEY_RFIC)
		}
		data[ota_constant.KEY_RFIC] = rfic_list
	}

	return data, nil
}

func getArbiList(dal *release.Dal, cp *release.CpRelease, original_arbi string) ([]string, error) {
	arbi_list := make([]string, 0, 5)

	//0. replace version
	arbi_primary, err := ReplaceVersionInPath(original_arbi, cp.Version)
	if err == nil && arbi_primary != "" {
		rel_arbi, err := release.FindArbiByPath(dal, arbi_primary)
		if err == nil && rel_arbi != nil {
			//			fmt.Println("primary arbi ", arbi_primary, " exist")
			arbi_list = append(arbi_list, arbi_primary)
			return arbi_list, nil
		} else {
			//			fmt.Println("primary arbi ", arbi_primary, " not exist")
		}
	}

	//1. search in db
	if !ota_constant.QUERY_MODE_STRICT {
		query := fmt.Sprintf("SELECT * FROM %s where cp_id=%d AND flag=%d", cp_constant.TABLE_ARBI, cp.Id, cp_constant.AVAILABLE_FLAG)
		arbis, err := release.FindArbiList(dal, query)
		if err != nil {
			return nil, err
		}

		if len(arbis) == 0 {
			return nil, fmt.Errorf("Cannot find the right Image.")
		}

		for _, arbi := range arbis {
			arbi_list = append(arbi_list, arbi.RelPath)
		}
	}

	return arbi_list, nil
}

func getGrbiList(dal *release.Dal, cp *release.CpRelease, original_grbi string) ([]string, error) {
	//0. search in primary CP
	grbi_list, err := doGetGrbiList(dal, cp, original_grbi)
	if err == nil && grbi_list != nil && len(grbi_list) > 0 {
		return grbi_list, nil
	}

	max_version := cp.VersionScalar
	min_version := (max_version/1000000)*1000000 + 484848

	query := fmt.Sprintf("SELECT * FROM %s WHERE mode='%s' AND sim='%s' AND prefix='%s' AND version_scalar<%d AND version_scalar >=%d AND flag=%d ORDER BY version_scalar DESC",
		cp_constant.TABLE_CP, cp.Mode, cp.Sim, cp.Prefix, max_version, min_version, cp_constant.AVAILABLE_FLAG)
	cps, err := release.FindCpReleaseList(dal, query)
	if err != nil {
		return nil, err
	}
	for _, _cp := range cps {
		grbi_list, err := doGetGrbiList(dal, _cp, original_grbi)
		if err == nil && grbi_list != nil && len(grbi_list) > 0 {
			return grbi_list, nil
		}
	}

	/*********************************************************************
	query := fmt.Sprintf("SELECT * FROM %s where cp_id=%d AND flag=%d", cp_constant.TABLE_GRBI, cp.Id, cp_constant.AVAILABLE_FLAG)
	grbis, err := release.FindGrbiList(dal, query)
	if err == nil && grbis != nil && len(grbis) > 0 {
		for _, grbi := range grbis {
			grbi_list = append(grbi_list, grbi.RelPath)
		}
		if len(grbi_list) > 0 {
			return grbi_list, nil
		}
	}

	//1. search in lower CP by replacing version and comparing basename
	query = fmt.Sprintf("SELECT * FROM %s WHERE mode='%s' AND sim='%s' AND prefix='%s' AND version_scalar<%d  AND flag=%d ORDER BY version_scalar DESC",
		cp_constant.TABLE_CP, cp.Mode, cp.Sim, cp.Prefix, cp.VersionScalar, cp_constant.AVAILABLE_FLAG)
	cps, err := release.FindCpReleaseList(dal, query)
	if err != nil {
		return nil, err
	}
	for _, _cp := range cps {
		grbi_list, err := doGetGrbiList(dal, _cp, original_grbi)
		if err == nil && grbi_list != nil && len(grbi_list) > 0 {
			return grbi_list, nil
		}
	}

	//2. search any available msa in lower CP
	for _, _cp := range cps {
		query = fmt.Sprintf("SELECT * FROM %s where cp_id=%d AND flag=%d", cp_constant.TABLE_GRBI, _cp.Id, cp_constant.AVAILABLE_FLAG)
		grbis, err := release.FindGrbiList(dal, query)
		if err == nil && grbis != nil && len(grbis) > 0 {
			for _, grbi := range grbis {
				grbi_list = append(grbi_list, grbi.RelPath)
			}
			if len(grbi_list) > 0 {
				return grbi_list, nil
			}
		}
	}
	*********************************************************************/

	return nil, nil
}

func doGetGrbiList(dal *release.Dal, cp *release.CpRelease, original_grbi string) ([]string, error) {
	grbi_list := make([]string, 0, 5)

	//find any MSA in current CP
	query := fmt.Sprintf("SELECT * FROM %s where cp_id=%d AND flag=%d", cp_constant.TABLE_GRBI, cp.Id, cp_constant.AVAILABLE_FLAG)
	grbis, err := release.FindGrbiList(dal, query)
	if err == nil && grbis != nil && len(grbis) > 0 {
		for _, grbi := range grbis {
			grbi_list = append(grbi_list, grbi.RelPath)
		}

		if len(grbi_list) > 0 {
			return grbi_list, nil
		}
	}

	/*********************************************************************
	//0. replace version
	grbi_primary, err := ReplaceVersionInPath(original_grbi, cp.Version)
	if err == nil && grbi_primary != "" {
		rel_grbi, err := release.FindGrbiByPath(dal, grbi_primary)
		if err == nil && rel_grbi != nil {
			//			fmt.Println(" primary grbi ", grbi_primary, " exist")
			grbi_list = append(grbi_list, grbi_primary)
			return grbi_list, nil
		} else {
			//			fmt.Println(" primary grbi ", grbi_primary, " not exist")
		}
	}

	//1. searching by basename
	query := fmt.Sprintf("SELECT * FROM %s where cp_id=%d AND flag=%d", cp_constant.TABLE_GRBI, cp.Id, cp_constant.AVAILABLE_FLAG)
	grbis, err := release.FindGrbiList(dal, query)
	if err == nil && grbis != nil && len(grbis) > 0 {
		original_name := pathutil.BaseName(original_grbi)
		for _, grbi := range grbis {
			base_name := pathutil.BaseName(grbi.RelPath)
			if base_name == original_name {
				grbi_list = append(grbi_list, grbi.RelPath)
			}
		}
		if len(grbi_list) > 0 {
			return grbi_list, nil
		}
	}
	*********************************************************************/

	return nil, fmt.Errorf("Cannot find grbi")
}

func getRficList(dal *release.Dal, cp *release.CpRelease, original_rfic string) ([]string, error) {
	rfic_list, err := doGetRficList(dal, cp, original_rfic)
	if err == nil && rfic_list != nil && len(rfic_list) > 0 {
		return rfic_list, nil
	}

	max_version := cp.VersionScalar
	min_version := (max_version/1000000)*1000000 + 484848

	query := fmt.Sprintf("SELECT * FROM %s WHERE mode='%s' AND sim='%s' AND prefix='%s' AND version_scalar<%d AND version_scalar >=%d AND flag=%d ORDER BY version_scalar DESC",
		cp_constant.TABLE_CP, cp.Mode, cp.Sim, cp.Prefix, max_version, min_version, cp_constant.AVAILABLE_FLAG)
	cps, err := release.FindCpReleaseList(dal, query)
	if err != nil {
		return nil, err
	}
	for _, _cp := range cps {
		rfic_list, err := doGetRficList(dal, _cp, original_rfic)
		if err == nil && rfic_list != nil && len(rfic_list) > 0 {
			return rfic_list, nil
		}
	}

	/*******************************************************************
	//0. replace version
	rfic_primary, err := ReplaceVersionInPath(original_rfic, cp.Version)
	if err == nil && rfic_primary != "" {
		rel_rfic, err := release.FindRficByPath(dal, rfic_primary)
		if err == nil && rel_rfic != nil {
			//			fmt.Println("primary rfic ", rfic_primary, " exist")
			rfic_list = append(rfic_list, rfic_primary)
			return rfic_list, nil
		} else {
			//			fmt.Println("primary rfic ", rfic_primary, " not exist")
		}
	}

	//1. search in db
	if !ota_constant.QUERY_MODE_STRICT {
		query := fmt.Sprintf("SELECT * FROM %s where cp_id=%d AND flag=%d", cp_constant.TABLE_RFIC, cp.Id, cp_constant.AVAILABLE_FLAG)
		rfics, err := release.FindRficList(dal, query)
		if err != nil {
			return nil, err
		}

		if len(rfics) == 0 {
			return nil, fmt.Errorf("Cannot find the right Image.")
		}

		for _, rfic := range rfics {
			rfic_list = append(rfic_list, rfic.RelPath)
		}
	}
	****************************************************************************/

	return nil, nil
}

func doGetRficList(dal *release.Dal, cp *release.CpRelease, original_rfic string) ([]string, error) {
	rfic_list := make([]string, 0, 5)

	//find any rfic in current CP
	query := fmt.Sprintf("SELECT * FROM %s where cp_id=%d AND flag=%d", cp_constant.TABLE_RFIC, cp.Id, cp_constant.AVAILABLE_FLAG)
	rfics, err := release.FindRficList(dal, query)
	if err == nil && rfics != nil && len(rfics) > 0 {
		for _, rfic := range rfics {
			rfic_list = append(rfic_list, rfic.RelPath)
		}

		if len(rfic_list) > 0 {
			return rfic_list, nil
		}
	}

	return nil, fmt.Errorf("Cannot find rfic")
}

func getCpList(dal *release.Dal, cp_info *CpInfo) ([]*release.CpRelease, error) {
	cp_list := make([]*release.CpRelease, 0, 10)
	mode := cp_info.Mode
	sim := cp_info.Sim
	version := cp_info.Version
	prefix := cp_info.Prefix
	version_scalar := cp_policy.QuantitateVersion(version)

	query := fmt.Sprintf("SELECT * FROM %s WHERE mode='%s' AND prefix='%s' AND sim='%s' AND flag=%d AND version_scalar > %d ORDER BY version_scalar DESC",
		cp_constant.TABLE_CP, mode, prefix, sim, cp_constant.AVAILABLE_FLAG, version_scalar)
	revel.INFO.Println("query higher cp: ", query)
	//	fmt.Println(query)
	list, err := doGetCpList(dal, query)
	if err != nil {
		return nil, err
	}
	cp_list = append(cp_list, list...)

	query = fmt.Sprintf("SELECT * FROM %s WHERE mode='%s' AND prefix='%s' AND sim='%s' AND flag=%d AND version_scalar < %d ORDER BY version_scalar DESC LIMIT 5",
		cp_constant.TABLE_CP, mode, prefix, sim, cp_constant.AVAILABLE_FLAG, version_scalar)
	revel.INFO.Println("query lower 5 cp: ", query)
	//	fmt.Println(query)
	list, err = doGetCpList(dal, query)
	if err != nil {
		return nil, err
	}
	cp_list = append(cp_list, list...)

	return cp_list, nil
}

func doGetCpList(dal *release.Dal, query string) ([]*release.CpRelease, error) {
	cps, err := release.FindCpReleaseList(dal, query)
	if err != nil {
		return nil, err
	}

	return cps, nil
}
