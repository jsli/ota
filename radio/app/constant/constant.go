package constant

import (
	cp_constant "github.com/jsli/cp_release/constant"
)

const (
	OTA_ROOT         = cp_constant.OTA_ROOT
	TMP_FILE_ROOT    = OTA_ROOT + "tmp/"
	DTIM_UPLOAD_ROOT = TMP_FILE_ROOT + "dtim_upload/"
	SCRIPTS_ROOT     = OTA_ROOT + "scripts/"
	FILTER_ROOT      = SCRIPTS_ROOT + "filter/"
	TEMPLATE_ROOT    = OTA_ROOT + "template/"

	ZIP_DIR_NAME           = "zip/"
	UPDATE_PKG_NAME        = "update_pkg.zip"
	RADIO_OTA_PACKAGE_NAME = "update.zip"
	RADIO_DTIM_NAME        = "Radio.dtim"
	RADIO_IMAGE_NAME       = "Radio.img"

	TEMPLATE_HELAN    = "HELAN"
	TEMPLATE_HELANLTE = "HELANLTE"
)

//command path
const (
	OTA_CMD_PARAM_PLATFORM_PREFIX = "--platform="   //[]
	OTA_CMD_PARAM_PRODUCT_PREFIX  = "--product="    //[ro.product.model]
	OTA_CMD_PARAM_OEM_PREFIX      = "--oem=marvell" //hard code here
	OTA_CMD_PARAM_OUTPUT_PREFIX   = "--output="
	OTA_CMD_PARAM_INPUT_PREFIX    = "--zipfile="

	RESIGN_DTIM_CMD  = SCRIPTS_ROOT + "dtim/resigndtim.rb"
	OTA_PKG_MAKE_CMD = "/home/manson/server/ota/new/radio/updatetool/updatemk"
)

const (
	RADIO_DTIM_SIZE = 4096

	TYPE_SINGLE      = cp_constant.SIM_SINGLE //2 image2
	TYPE_DSDS        = cp_constant.SIM_DSDS   //4 images
	TYPE_SINGLE_RFIC = "SINGLE_RFIC"          //3 images
	TYPE_DSDS_RFIC   = "DSDS_RFIC"            //6 images

	ID_ARBI = "ARBI"
	ID_GRBI = "GRBI"
	ID_RFIC = "RFIC"
	ID_ARB2 = "ARB2"
	ID_GRB2 = "GRB2"
	ID_RFI2 = "RFI2"

	KEY_ARBEL            = "ARBEL"
	KEY_MSA              = "MSA"
	KEY_RFIC             = "RFIC"
	KEY_RESULT_CURRENT   = "current"
	KEY_RESULT_AVAILABLE = "available"
	KEY_VERSION          = "version"
	KEY_IMAGES           = "images"

	KEY_URL   = "url"
	KEY_MD5   = "md5"
	KEY_SIZE  = "size"
	KEY_ERROR = "error"

	DROPPED_FLAG   = -1
	AVAILABLE_FLAG = 1
	DISABLE_FLAG   = 2

	PLATFORM_JB_4_2 = "jb-4.2"
	PLATFORM_JB_4_3 = "jb-4.3"

	MODEL_1088 = "pxa1088ff_def"
	MODEL_1T88 = "pxa1t88ff_def"
	MODEL_1L88 = "pxa1l88ff_def"

	BOARD_FF  = "FF"
	BOARD_DKB = "DKB"
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
	MODEL_TO_TEMPLATE = map[string]string{
		MODEL_1088: TEMPLATE_HELAN,
		MODEL_1T88: TEMPLATE_HELAN,
		MODEL_1L88: TEMPLATE_HELANLTE,
	}
	MODEL_TO_PLATFORM = map[string]string{
		MODEL_1088: PLATFORM_JB_4_2,
		MODEL_1T88: PLATFORM_JB_4_2,
		MODEL_1L88: PLATFORM_JB_4_3,
	}

	KEY_LIST   = []string{KEY_ARBEL, KEY_MSA, KEY_RFIC}
	MODEL_LIST = []string{MODEL_1088, MODEL_1T88, MODEL_1L88}
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
	TestDataLTER = [][]string{
		{"ARBI", "LWG", "SINGLE", "LWG/HL_CP_2.30.000/HL_CP/Seagull/HL_LWG_DKB.bin"},
		{"GRBI", "LWG", "SINGLE", "LWG/HL_CP_2.30.000/HL_MSA_2.30.000/HL_LWG_M09_B0_SKL_Flash.bin"},
		{"RFIC", "LWG", "SINGLE", "LWG/HL_CP_2.30.000/RFIC/1920_FF/Skylark_LWG.bin"},
		{"ARB2", "LTG", "SINGLE", "LTG/HL_CP_3.30.000/HL_CP/Seagull/HL_LTG_DL_DKB.bin"},
		{"GRB2", "LTG", "SINGLE", "LTG/HL_CP_3.30.000/HL_MSA_3.30.000/HL_DL_M09_Y0_AI_SKL_Flash.bin"},
		{"RFI2", "LTG", "SINGLE", "LTG/HL_CP_3.30.000/RFIC/1920_FF/Skylark_LTG.bin"},
	}
)
