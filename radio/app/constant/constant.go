package constant

import (
	cp_constant "github.com/jsli/cp_release/constant"
)

//path
const (
	OTA_ROOT               = "/home/manson/OTA/"
	TMP_FILE_ROOT          = OTA_ROOT + "tmp/"
	SCRIPTS_ROOT           = OTA_ROOT + "scripts/"
	FILTER_ROOT            = SCRIPTS_ROOT + "filter/"
	TEMPLATE_ROOT          = OTA_ROOT + "template/"
	RELEASE_ROOT           = OTA_ROOT + "release/"
	CP_RELEASE_ROOT        = RELEASE_ROOT + "CP/"
	CP_SERVER_MIRROR_ROOT  = CP_RELEASE_ROOT + "CP_SERVER_MIRROR/"
	RADIO_OTA_RELEASE_ROOT = RELEASE_ROOT + "radio_ota/"
	TOOLS_ROOT             = OTA_ROOT + "updatetool/"

	ZIP_DIR_NAME           = "zip/"
	UPDATE_PKG_NAME        = "update_pkg.zip"
	RADIO_OTA_PACKAGE_NAME = "update.zip"
	RADIO_DTIM_NAME        = "radio.dtim"
	RADIO_IMAGE_NAME       = "radio.img"
	CHECKSUM_TXT_NAME      = "checksum.txt"

	UPDATE_CMD_NAME  = "updatemk"
	RESIGN_DTIM_NAME = "dtim/resigndtim.rb"
	GZIP_CMD_NAME    = "gzip"

	TEMPLATE_HELAN_ROOT    = TEMPLATE_ROOT + "HELAN/"
	TEMPLATE_HELANLTE_ROOT = TEMPLATE_ROOT + "HELANLTE/"
)

//command
const (
	OTA_CMD_PARAM_PLATFORM_PREFIX = "--platform="   //[]
	OTA_CMD_PARAM_PRODUCT_PREFIX  = "--product="    //[ro.product.model]
	OTA_CMD_PARAM_OEM_PREFIX      = "--oem=marvell" //hard code here
	OTA_CMD_PARAM_OUTPUT_PREFIX   = "--output="
	OTA_CMD_PARAM_INPUT_PREFIX    = "--zipfile="

	RESIGN_DTIM_CMD = SCRIPTS_ROOT + RESIGN_DTIM_NAME
	OTA_MAKE_CMD    = TOOLS_ROOT + UPDATE_CMD_NAME
)

const (
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
	KEY_URL              = "url"
	KEY_MD5              = "md5"
	KEY_SIZE             = "size"
	KEY_CREATED_TIME     = "created_time"
	KEY_ERROR            = "error"

	TAG_1088   = "1088"
	TAG_1T88   = "1T88"
	TAG_1L88   = "1L88"
	MODEL_1088 = "PXA1088"
	MODEL_1920 = "PXA1L88"

	TAG_JB_4_2      = "4.2"
	TAG_JB_4_3      = "4.3"
	PLATFORM_JB_4_2 = "jb-4.2"
	PLATFORM_JB_4_3 = "jb-4.3"

	BOARD_FF  = "FF"
	BOARD_DKB = "DKB"
)

//database
const (
	TABLE_RADIO_OTA_RELEASE     = "radio_ota_release"
	TABLE_RELEASE_CREATION_TASK = "release_creation_task"

	FLAG_DROPPED       = -1
	FLAG_AVAILABLE     = 1
	FLAG_DISABLE       = 2
	FLAG_INIT          = 4
	FLAG_CREATING      = 8
	FLAG_CREATED       = 16
	FLAG_CREATE_FAILED = 32

	ERROR_CODE_NOERR         = 0
	ERROR_CODE_DROPPED       = FLAG_DROPPED
	ERROR_CODE_DISABLE       = FLAG_DISABLE
	ERROR_CODE_NOT_CREATED   = FLAG_INIT
	ERROR_CODE_CREATING      = FLAG_CREATING
	ERROR_CODE_CREATE_FAILED = FLAG_CREATE_FAILED
)

const (
	TIME_FMT = "2006-01-02 15:04:05"
)

var (
	MODEL_TO_TEMPLATE = map[string]string{
		MODEL_1088: TEMPLATE_HELAN_ROOT,
		MODEL_1920: TEMPLATE_HELANLTE_ROOT,
	}

	MODE_TO_ROOT_PATH = cp_constant.MODE_TO_ROOT_PATH

	KEY_LIST        = []string{KEY_ARBEL, KEY_MSA, KEY_RFIC}
	MODEL_LIST      = []string{MODEL_1088, MODEL_1920}
	IMAGE_ID_LIST   = []string{ID_ARBI, ID_GRBI, ID_RFIC, ID_ARB2, ID_GRB2, ID_RFI2}
	GZIP_CMD_PARAMS = []string{"-n", "-9"}
)
