package policy

import (
	"fmt"
	"sort"

	cp_constant "github.com/jsli/cp_release/constant"
	cp_policy "github.com/jsli/cp_release/policy"
	"github.com/jsli/ota/radio/app/models"
)

func SortCps(update_request *models.UpdateRequest) []models.CpRequest {
	sorted := make([]models.CpRequest, 2)
	request_cps := update_request.Cps
	if len(request_cps) == 1 {
		return request_cps
	}
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

type SortInt []int64

func (p SortInt) Len() int               { return len(p) }
func (p SortInt) Less(i int, j int) bool { return p[j] < p[i] }
func (p SortInt) Swap(i int, j int)      { p[i], p[j] = p[j], p[i] }

var DEBUG bool = false

func SortCpArray(cp_array []models.CpNode) []models.CpNode {
	order_map := make(map[int64]models.CpNode)
	version_int_arr := SortInt{}
	for _, node := range cp_array {
		version_int := cp_policy.QuantitateVersion(node.VersionNo)
		version_int_arr = append(version_int_arr, version_int)
		order_map[version_int] = node
	}

	sort.Sort(version_int_arr)
	sorted_arr := make([]models.CpNode, 0, 50)
	for _, version := range version_int_arr {
		sorted_arr = append(sorted_arr, order_map[version])
	}

	if DEBUG {
		for _, v1 := range cp_array {
			fmt.Print(v1.VersionNo)
			fmt.Print("  ")
		}
		fmt.Print("\n")

		for _, v2 := range sorted_arr {
			fmt.Print(v2.VersionNo)
			fmt.Print("  ")
		}
		fmt.Print("\n")
	}

	return sorted_arr
}
