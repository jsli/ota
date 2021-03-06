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
	CP_RELEASE_SYNC_ROOT   = CP_RELEASE_ROOT + "CP_SYNC/"
	CP_RELEASE_ROOT_FINAL  = CP_SERVER_MIRROR_ROOT
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
	ID_RFIC = "RFBI"
	ID_ARB2 = "ARB2"
	ID_GRB2 = "GRB2"
	ID_RFI2 = "RFB2"

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

	RESULT_KEY_IMAGE_NAME    = "image_name"
	RESULT_KEY_IMAGE_ARRAY   = "image_array"
	RESULT_KEY_MODE_NAME     = "mode_name"
	RESULT_KEY_VERSION_ARRAY = "version_array"
	RESULT_KEY_VERSION_NO    = "version_no"

	TAG_1088   = "1088"
	TAG_1T88   = "1T88"
	TAG_1920   = "1920"
	MODEL_1088 = "PXA1088"
	MODEL_1920 = "PXA1920"

	TAG_JB_4_2      = "4.2"
	TAG_JB_4_3      = "4.3"
	PLATFORM_JB_4_2 = "jb-4.2"
	PLATFORM_JB_4_3 = "jb-4.3"

	BOARD_FF  = "FF"
	BOARD_DKB = "DKB"
)

const (
	REQUEST_PARAM_APIVERSION = "api_version"
	API_VERSION_1_0          = "1.0"
	API_VERSION_2_0          = "2.0"
	CURRENT_API_VERSION      = API_VERSION_2_0
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

	ERROR_CODE_NOERR                 = 0
	ERROR_CODE_DROPPED               = -1
	ERROR_CODE_DISABLE               = 2
	ERROR_CODE_NOT_CREATED           = 4
	ERROR_CODE_CREATING              = 8
	ERROR_CODE_CREATE_FAILED         = 32
	ERROR_CODE_CREATE_REQUEST_FAILED = 64
	ERROR_CODE_MAINTAIN              = 128
	ERROR_CODE_INVALIDATED_DTIM      = 256
	ERROR_CODE_NO_AVAILABLE_UPDATE   = 512
	ERROR_CODE_INVALIDATED_REQUEST   = 1024

	ERROR_MSG_NO_AVAILABLE_CP    = "Cannot find available CP release."
	ERROR_MSG_NO_AVAILABLE_IMAGE = "Cannot find available image [%s]."
	ERROR_MSG_NO_ILLEGAL_REQUEST = "Illegal update request: [%s]."

	RETRY_COUNT = 5
)

const (
	TIME_FMT = "2006-01-02 15:04:05"
)

const (
	QUERY_MODE_STRICT = true
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
