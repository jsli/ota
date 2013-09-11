package constant

import (
	"regexp"
)

const (
	FULL_BUILD_DIR   = "/home/manson/temp/test/CM/"
	UPLOAD_ROOT_DIR  = "/home/manson/temp/test/CP/"
	LOG_ROOG_DIR     = "/home/manson/temp/test/logs/"
	RELEASE_DIR      = "/home/manson/temp/test/release/"
	TEMP_DIR         = "/home/manson/temp/test/tmp/"
	UPLOAD_LOG_NAME  = "upload.log"
	RADIO_IMAGE_NAME = "radio.img"
	RADIO_PKG_NAME   = "update.zip"
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

var CP_VERSION_REX = regexp.MustCompile(`^\d+\.\d+\.\d+$`)

var COPY_FILE_LIST []string = []string{"SYSTEM/build.prop", "RECOVERY/RAMDISK/etc/recovery.fstab"}
var COPY_DIR_LIST []string = []string{"OTA/", "META/"}
