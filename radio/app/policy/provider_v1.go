package policy

import (
	"fmt"
	"github.com/jsli/cp_release/release"
	ota_constant "github.com/jsli/ota/radio/app/constant"
	"github.com/jsli/ota/radio/app/models"
)

type ContentProviderV1 struct {
	CommonContentProvider
}

type DataV1 struct {
}

func (cpv1 *ContentProviderV1) ProvideQueryData(dal *release.Dal, dtim_info *DtimInfo, result models.DataSetter) error {
	result_data := make(map[string]interface{})
	current := models.CurrentCps{}
	available := make(map[string]map[string]map[string][]string)
	for _, cp_info := range dtim_info.CpMap {
		images := getCurrentCpImages(cp_info)
		if len(images) >= 2 {
			ci := models.CpAndImages{}
			ci.Version = cp_info.Version
			ci.Images = images
			current[cp_info.Mode] = ci
		}

		data, err := getCpAndImages(dal, cp_info, dtim_info.HasRFIC)
		if err != nil {
			return err
		}

		available[cp_info.Mode] = data
	}

	if len(available) == 0 || len(current) == 0 {
		return fmt.Errorf(ota_constant.ERROR_MSG_NO_AVAILABLE_CP)
	}

	result_data["available"] = available
	result_data["current"] = current
	result.SetData(result_data)

	return nil
}
