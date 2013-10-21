package policy

import (
	cp_constant "github.com/jsli/cp_release/constant"
	"github.com/jsli/ota/radio/app/models"
)

func SortCps(update_request *models.UpdateRequest) []models.CpRequest {
	sorted := make([]models.CpRequest, 2)
	request_cps := update_request.Cps
	for _, cp := range request_cps {
		switch cp.Mode {
		case cp_constant.MODE_HLWB, cp_constant.MODE_HLTD, cp_constant.MODE_LWG:
			sorted[0] = cp
		case cp_constant.MODE_HLWB_DSDS, cp_constant.MODE_HLTD_DSDS, cp_constant.MODE_LTG:
			sorted[1] = cp
		}
	}
	return sorted
}
