package policy

import (
	"fmt"
	cp_constant "github.com/jsli/cp_release/constant"
	"github.com/jsli/cp_release/release"
	ota_constant "github.com/jsli/ota/radio/app/constant"
	"github.com/jsli/ota/radio/app/models"
)

type ContentProviderV2 struct {
	CommonContentProvider
}

func (cpv2 *ContentProviderV2) ProvideQueryData(dal *release.Dal, dtim_info *DtimInfo, result models.DataSetter) error {
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

	mode_count := len(available)
	available_array := make([]models.ModeNode, mode_count)
	for mode, versions := range available {
		mode_node := models.ModeNode{}
		mode_node.Mode = mode
		version_count := len(versions)
		cp_array := make([]models.CpNode, 0, version_count)
		for version, images := range versions {
			cp := models.CpNode{}
			image_count := len(images)
			image_array := make([]models.ImageNode, image_count)
			for id, image_list := range images {
				image_node := models.ImageNode{}
				image_node.ImageName = id
				image_node.Images = image_list

				switch id {
				case ota_constant.KEY_ARBEL:
					image_array[0] = image_node
				case ota_constant.KEY_MSA:
					image_array[1] = image_node
				case ota_constant.KEY_RFIC:
					image_array[2] = image_node
				}
			}
			cp.VersionNo = version
			cp.ImageArray = image_array
			cp_array = append(cp_array, cp)
		}
		mode_node.CpArray = SortCpArray(cp_array)

		switch mode {
		case cp_constant.MODE_HLWB, cp_constant.MODE_HLTD, cp_constant.MODE_LWG:
			available_array[0] = mode_node
		case cp_constant.MODE_HLWB_DSDS, cp_constant.MODE_HLTD_DSDS, cp_constant.MODE_LTG:
			available_array[1] = mode_node
		}
	}

	result_data["available"] = available_array
	result_data["current"] = current
	result.SetData(result_data)

	return nil
}
