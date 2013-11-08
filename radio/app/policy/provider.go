package policy

import (
	"fmt"
	cp_constant "github.com/jsli/cp_release/constant"
	cp_policy "github.com/jsli/cp_release/policy"
	"github.com/jsli/cp_release/release"
	"github.com/jsli/gtbox/file"
	ota_constant "github.com/jsli/ota/radio/app/constant"
	"github.com/jsli/ota/radio/app/models"
	"strings"
)

func ProvideRadioRelease(dal *models.Dal, dtim_info *DtimInfo, result *models.RadioOtaReleaseResult, fp string) (*models.RadioOtaRelease, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE fingerprint='%s' AND flag=%d LIMIT 1",
		ota_constant.TABLE_RADIO_OTA_RELEASE, fp, cp_constant.AVAILABLE_FLAG)
	//	fmt.Println(query)
	release, err := models.FindRadioOtaRelease(dal, query)
	if err != nil {
		return nil, err
	}

	if release != nil {
		return release, nil
	}
	return nil, nil
}

func ProvideQueryData(dal *release.Dal, dtim_info *DtimInfo, result *models.QueryResult) error {
	//	result.Data = make(map[string]interface{})
	current := make(map[string]interface{})
	available := make(map[string]interface{})
	for _, cp_info := range dtim_info.CpMap {
		image_list := getCurrentInfo(cp_info)
		if len(image_list) > 0 {
			current_detail := make(map[string]interface{})
			current_detail[ota_constant.KEY_VERSION] = cp_info.Version
			current_detail[ota_constant.KEY_IMAGES] = image_list
			current[cp_info.Mode] = current_detail
		}

		data, err := getCpAndImages(dal, cp_info, dtim_info.HasRFIC)
		if err != nil {
			return err
		}
		filterByParams(data, dtim_info)
		filterByRuleFile(data, cp_info)
		available[cp_info.Mode] = data
	}
	result.Data.Available = available
	result.Data.Current = current

	return nil
}

func getCurrentInfo(cp_info *CpInfo) map[string]string {
	data := make(map[string]string)
	for key, value := range cp_info.ImageMap {
		switch key {
		case ota_constant.ID_ARBI, ota_constant.ID_ARB2:
			data[ota_constant.KEY_ARBEL] = value.Path
		case ota_constant.ID_GRBI, ota_constant.ID_GRB2:
			data[ota_constant.KEY_MSA] = value.Path
		case ota_constant.ID_RFIC, ota_constant.ID_RFI2:
			data[ota_constant.KEY_RFIC] = value.Path
		}
	}
	return data
}

func filterByParams(data map[string]map[string][]string, dtim_info *DtimInfo) {
}

func filterByRuleFile(data map[string]map[string][]string, cp_info *CpInfo) {
	filter_map := make(map[string][]string)
	for _, key := range ota_constant.KEY_LIST {
		filter_map[key] = getFilterList(cp_info.Mode, key)
	}

	for key, _ := range data {
		original_data := data[key]
		for key, value := range original_data {
			filter := filter_map[key]
			filtered := make([]string, 0, 10)
			for _, path := range value {
				if check(path, filter) {
					filtered = append(filtered, path)
				} else {
				}
			}
			original_data[key] = filtered
		}
	}
}

func getFilterList(mode string, key string) []string {
	path := fmt.Sprintf("%s%s_%s", cp_constant.FILTER_ROOT, mode, key)
	content, err := file.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	content = strings.TrimSpace(content)
	return strings.Split(content, "\n")
}

func check(str string, filter []string) bool {
	for _, value := range filter {
		if strings.Contains(str, value) {
			return false
		}
	}
	return true
}

/*
 * return:
 * map[string]map[string][]string : map[CP_VERSION] map[IMAGE_KEY] images list
 * @IMAGE_ID refer to
 *	KEY_ARBEL = "ARBEL"
 *	KEY_MSA = "MSA"
 *	KEY_RFIC = "RFIC"
 */
func getCpAndImages(dal *release.Dal, cp_info *CpInfo, hasRFIC bool) (map[string]map[string][]string, error) {
	cp_list, err := getCpList(dal, cp_info)
	if err != nil {
		return nil, err
	}

	if cp_list != nil {
		data := make(map[string]map[string][]string)
		for _, cp := range cp_list {
			image_list, err := getImagesByCp(dal, cp, hasRFIC)
			if err != nil {
				//return nil, err
				continue
			}
			data[cp.Version] = image_list
		}
		return data, nil
	}
	return nil, nil
}

/*
 * return:
 * map[string][]string : map[CP_VERSION] images list
 */
func getImagesByCp(dal *release.Dal, cp *release.CpRelease, hasRFIC bool) (map[string][]string, error) {
	data := make(map[string][]string)

	arbi_list, err := getArbiList(dal, cp)
	if err != nil {
		return nil, err
	}
	if arbi_list != nil {
		data[ota_constant.KEY_ARBEL] = arbi_list
	}

	grbi_list, err := getGrbiList(dal, cp)
	if err != nil {
		return nil, err
	}
	if grbi_list != nil {
		data[ota_constant.KEY_MSA] = grbi_list
	}

	if hasRFIC {
		rfic_list, err := getRficList(dal, cp)
		if err != nil {
			return nil, err
		}
		if rfic_list != nil {
			data[ota_constant.KEY_RFIC] = rfic_list
		}
	}

	return data, nil
}

func getArbiList(dal *release.Dal, cp *release.CpRelease) ([]string, error) {
	arbi_list := make([]string, 0, 10)
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
	return arbi_list, nil
}

func getGrbiList(dal *release.Dal, cp *release.CpRelease) ([]string, error) {
	max_version := cp_policy.QuantitateVersion(cp.Version)
	min_version := (max_version / 1000000) * 1000000
	query := fmt.Sprintf("SELECT * FROM %s WHERE mode='%s' AND sim='%s' AND (version_scalar<=%d AND version_scalar>=%d ) AND flag=%d ORDER BY version_scalar DESC",
		cp_constant.TABLE_CP, cp.Mode, cp.Sim, max_version, min_version, cp_constant.AVAILABLE_FLAG)
	cps, err := release.FindCpReleaseList(dal, query)
	if err != nil {
		return nil, err
	}

	grbi_list := make([]string, 0, 10)
	for _, _cp := range cps {
		query = fmt.Sprintf("SELECT * FROM %s where cp_id=%d AND flag=%d", cp_constant.TABLE_GRBI, _cp.Id, cp_constant.AVAILABLE_FLAG)
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
		break
	}

	if len(grbi_list) == 0 {
		return nil, fmt.Errorf("Cannot find the right Image.")
	}

	return grbi_list, nil
}

func getRficList(dal *release.Dal, cp *release.CpRelease) ([]string, error) {
	rfic_list := make([]string, 0, 10)
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
	return rfic_list, nil
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
	fmt.Println(query)
	list, err := doGetCpList(dal, query)
	if err != nil {
		return nil, err
	}
	cp_list = append(cp_list, list...)

	query = fmt.Sprintf("SELECT * FROM %s WHERE mode='%s' AND prefix='%s' AND sim='%s' AND flag=%d AND version_scalar < %d ORDER BY version_scalar DESC LIMIT 5",
		cp_constant.TABLE_CP, mode, prefix, sim, cp_constant.AVAILABLE_FLAG, version_scalar)
	fmt.Println(query)
	list, err = doGetCpList(dal, query)
	if err != nil {
		return nil, err
	}
	cp_list = append(cp_list, list...)

	if len(cp_list) == 0 {
		return nil, nil
	}
	return cp_list, nil
}

func doGetCpList(dal *release.Dal, query string) ([]*release.CpRelease, error) {
	cps, err := release.FindCpReleaseList(dal, query)
	if err != nil {
		return nil, err
	}

	if len(cps) > 0 {
		return cps, nil
	} else {
		return nil, nil
	}
}
