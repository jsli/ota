package constant

import (
	cp_constant "github.com/jsli/cp_release/constant"
	"regexp"
)

const (
	OTA_ROOT         = "/home/manson/OTA/"
	TMP_FILE_ROOT    = OTA_ROOT + "tmp/"
	DTIM_UPLOAD_ROOT = TMP_FILE_ROOT + "dtim_upload/"

	RADIO_OTA_PACKAGE_NAME = "update.zip"
)

var (
	//ARBI|LTG|SINGLE|HLLTE/HLLTE_CP_2.29.000/Seagull/HL_LTG.bin
	//GRBI|LTG|SINGLE|HLLTE/HLLTE_CP_2.29.000/TTD_WK_NL_MSA_2.29.000/HL_DL_M09_Y0_AI_SKL_Flash.bin
	TestDataHLTD = [][]string{
		{"ARBI", "TG", "SINGLE", "HLTD/HLTD_CP_2.42.000/xxx.bin"},
		{"GRBI", "TG", "SINGLE", "HLTD/HLTD_CP_2.42.000/xxxFlash.bin"},
	}
	TestDataHLWB = [][]string{
		{"ARBI", "WG", "SINGLE", "HLWB/HLWB_CP_1.55.000/xxx.bin"},
		{"GRBI", "WG", "SINGLE", "HLWB/HLWB_CP_1.55.000/xxxFlash.bin"},
	}
	TestDataHLTD_DSDS = [][]string{
		{"ARBI", "TG", "SINGLE", "HLTD/HLTD_CP_2.42.000/xxx.bin"},
		{"GRBI", "TG", "SINGLE", "HLTD/HLTD_CP_2.42.000/xxxFlash.bin"},
		{"ARB2", "TG", "DSDS", "HLTD_DSDS/HLTD_DSDS_CP_3.28.000/xxx.bin"},
		{"GRB2", "TG", "DSDS", "HLTD_DSDS/HLTD_DSDS_CP_3.28.000/xxxFlash.bin"},
	}
	TestDataHLWB_DSDS = [][]string{
		{"ARBI", "WG", "SINGLE", "HLWB/HLWB_CP_1.55.000/xxx.bin"},
		{"GRBI", "WG", "SINGLE", "HLWB/HLWB_CP_1.55.000/xxxFlash.bin"},
		{"ARB2", "WG", "DSDS", "HLWB_DSDS/HLWB_CP_2.58.917/xxx.bin"},
		{"GRB2", "WG", "DSDS", "HLWB_DSDS/HLWB_CP_2.58.917/xxxFlash.bin"},
	}
	TestDataHLTDR = [][]string{
		{"ARBI", "TG", "SINGLE", "HLTD/HLTD_CP_2.42.000/xxx.bin"},
		{"GRBI", "TG", "SINGLE", "HLTD/HLTD_CP_2.42.000/xxxFlash.bin"},
		{"RFIC", "TG", "SINGLE", "HLTD/HLTD_CP_2.42.000/xxxrfic.bin"},
	}
)

const (
	RADIO_DTIM_SIZE = 4096
	RADIO_DTIM_NAME = "radio.dtim"

	TYPE_SINGLE      = cp_constant.SIM_SINGLE
	TYPE_DSDS        = cp_constant.SIM_DSDS
	TYPE_SINGLE_RFIC = "SINGLE_RFIC"
	TYPE_DSDS_RFIC   = "DSDS_RFIC"

	ID_ARBI = "ARBI"
	ID_GRBI = "GRBI"
	ID_RFIC = "RFIC"
	ID_ARB2 = "ARB2"
	ID_GRB2 = "GRB2"
	ID_RFI2 = "RFI2"
)

const (
	TABLE_RADIO_OTA_RELEASE = "radio_ota_release"
)

var (
	ID_TO_TYPE = map[string]string{
		ID_ARBI: TYPE_SINGLE,
		ID_GRBI: TYPE_SINGLE,
		ID_ARB2: TYPE_DSDS,
		ID_GRB2: TYPE_DSDS,
	}
)

const (
	OTA_RELEASE_STATUS_DISABLED  = 0
	OTA_RELEASE_STATUS_AVAILABLE = 1
	OTA_RELEASE_STATUS_DROPPED   = 2
)

const (
	VERSION_DIVIDER = "-"
)

const (
	FULL_BUILD_DIR   = "/home/manson/temp/test/CM/"
	UPLOAD_ROOT_DIR  = "/home/manson/temp/test/CP/"
	LOG_ROOG_DIR     = "/home/manson/temp/test/logs/"
	RELEASE_DIR      = "/home/manson/temp/test/release/"
	TEMP_DIR         = "/home/manson/temp/test/tmp/"
	RADIO_IMAGE_DIR  = "/home/manson/temp/test/radio_image/"
	UPLOAD_LOG_NAME  = "upload.log"
	RADIO_IMAGE_NAME = "radio.img"
	OTA_PKG_NAME     = "update.zip"
	UPDATE_PKG_NAME  = "update_pkg.zip"
	MD5_FILE_NAME    = "checksum.txt"
	RADIO_IMAGE_SIZE = 20971520
	ZIP_SUFFIX       = ".zip"
)

const (
	DEFAULT_DIR_ACCESS  = 0755
	DEFAULT_FILE_ACCESS = 0644
	IMAGE_FILE_ACCESS   = 0644
	DEFAULT_BUFFER_SIZE = 4096
)

const (
	OTA_CMD_PARAM_PLATFORM_PREFIX = "--platform="
	OTA_CMD_PARAM_PRODUCT_PREFIX  = "--product="
	OTA_CMD_PARAM_OEM_PREFIX      = "--oem="
	OTA_CMD_PARAM_OUTPUT_PREFIX   = "--output="
	OTA_CMD_PARAM_INPUT_PREFIX    = "--zipfile="
)

const (
	SINGLE_CP  = "single_cp.bin"
	SINGLE_DSP = "single_dsp.bin"
	DSDS_CP    = "dsds_cp.bin"
	DSDS_DSP   = "dsds_dsp.bin"
)

const (
	OTA_PKG_MAKE_CMD = "/home/manson/server/ota/new/radio/updatetool/updatemk"
)

var OFFSET_MAP = map[string]int64{
	SINGLE_CP: 0, SINGLE_DSP: 8388608,
	DSDS_CP: 10485760, DSDS_DSP: 18874368,
}

var CP_VERSION_REX = regexp.MustCompile(`^\d+\.\d+\.\d+$`)
var MODEL_2_CPPREFIX = map[string]string{"pxa986ff_def": "KL", "pxa988ff_def": "EM", "pxa1088ff_def": "HL_WB", "pxa1t88ff_def": "HL_TD"}
var COPY_FILE_LIST []string = []string{"SYSTEM/build.prop", "RECOVERY/RAMDISK/etc/recovery.fstab"}
var COPY_DIR_LIST []string = []string{"OTA/", "META/"}
var UPLOAD_FILE_LIST = []string{"cp", "dsp"}

var TYPE_LIST = []string{"single", "dsds"}
var MODEL_LIST = []string{"pxa1088ff_def", "pxa1t88ff_def"}
