package utils

import (
	"github.com/jsli/ota/radio/app/constant"
	"regexp"
	"strconv"
	"strings"
)

//files and their offset must match!
var FILE_KEY_LIST = [4]string{"single_cp", "single_dsp", "dsds_cp", "dsds_dsp"}
var FILE_OFFSET_LIST = [4]int64{0, 8388608, 10485760, 18874368}

var MODEL_2_CPPREFIX = map[string]string{"pxa986ff_def": "KL", "pxa988ff_def": "EM", "pxa1088ff_def": "HL_WB", "pxa1t88ff_def": "HL_TD"}

//HLTD_CP_2.48.000:TTD_WK_HLTD_MSA_2.48.000
var CPVPattern, _ = regexp.Compile(`\d+\.\d+\.\d+`)

func IsAvailableCPV(cpv string) (bool, string) {
	cpv_list := CPVPattern.FindAllString(cpv, -1)
	if cpv_list == nil || len(cpv_list) != 2 {
		return false, ""
	} else {
		return true, cpv_list[0]
	}
}

func IsAvailableType(t string) bool {
	for _, _t := range constant.TYPE_LIST {
		if t == _t {
			return true
		}
	}
	return false
}

func IsAvailableModel(model string) bool {
	for _, m := range constant.MODEL_LIST {
		if model == m {
			return true
		}
	}
	return false
}

func ConvertString2Int64(str string) int64 {
	v_int64, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return -1
	}
	return v_int64
}

func ConvertVersion2Int64(version string) int64 {
	v_list := SplitVersion(version)
	v_int_str := strings.Join(v_list, "")
	v_int64, err := strconv.ParseInt(v_int_str, 10, 64)
	if err != nil {
		return -1
	}
	return v_int64
}

func SplitVersion(path string) []string {
	return strings.Split(path, ".")
}
