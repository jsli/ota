package constant

import ()

//path
const (
	OTA_ROOT               = "/home/manson/OTA/"
	TMP_FILE_ROOT          = OTA_ROOT + "tmp/"
	DTIM_UPLOAD_ROOT       = TMP_FILE_ROOT + "dtim_upload/"
	SCRIPTS_ROOT           = OTA_ROOT + "scripts/"
	FILTER_ROOT            = SCRIPTS_ROOT + "filter/"
	TEMPLATE_ROOT          = OTA_ROOT + "template/"
	RELEASE_ROOT           = OTA_ROOT + "release/"
	RADIO_OTA_RELEASE_ROOT = RELEASE_ROOT + "RADIO_OTA/"
	TOOLS_ROOT             = OTA_ROOT + "updatetool/"

	ZIP_DIR_NAME           = "zip/"
	UPDATE_PKG_NAME        = "update_pkg.zip"
	RADIO_OTA_PACKAGE_NAME = "update.zip"
	RADIO_DTIM_NAME        = "Radio.dtim"
	RADIO_IMAGE_NAME       = "Radio.img"
	CHECKSUM_TXT_NAME      = "checksum.txt"

	UPDATE_CMD_NAME  = "updatemk"
	RESIGN_DTIM_NAME = "dtim/resigndtim.rb"

	TEMPLATE_HELAN    = "HELAN"
	TEMPLATE_HELANLTE = "HELANLTE"
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
	KEY_ERROR            = "error"

	PLATFORM_JB_4_2 = "jb-4.2"
	PLATFORM_JB_4_3 = "jb-4.3"

	MODEL_1088 = "pxa1088ff_def"
	MODEL_1T88 = "pxa1t88ff_def"
	MODEL_1L88 = "pxa1l88ff_def"

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
)

var (
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
