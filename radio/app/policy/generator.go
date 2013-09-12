package policy

import (
	"fmt"
	//	"github.com/jsli/gtbox/file"
	"github.com/jsli/ota/radio/app/models"
	"strings"
	ota_constant "github.com/jsli/ota/radio/app/constant"
)

func GenerateOtaPackage(dal *models.Dal, parsedParams *ParsedParams, versions string) error {
	var mode, version string
//	image_list := make([]string, 0, 4)
	version_list := strings.Split(versions, "-")
	single_info, ok := parsedParams.CpMap[ota_constant.TYPE_SINGLE]
	if ok {
		arbi := single_info.ImageMap[ota_constant.ID_ARBI]
		mode = arbi.Mode
		version = version_list[0]
	}

	dsds_info, ok := parsedParams.CpMap[ota_constant.TYPE_DSDS]
	if ok {
		arbi := dsds_info.ImageMap[ota_constant.ID_ARB2]
		mode = fmt.Sprintf("%s-%s", mode, arbi.Mode)
		version = fmt.Sprintf("%s-%s", version, version_list[1])
	}

	//1.create new dtim (xxx.rb + radio.dtim)
	//2.create radio.img (dtim + cp)

	//3.copy template files and generate update_pkg.zip (template + radio.img)

	//4.generate update.zip (updatetool + update_pkg.zip)

	//5.insert db
	release := &models.RadioOtaRelease{}
	release.Mode = mode
	release.Version = version
	release.Md5 = "md5md5md5md5md5md5md5md5md5md5md5md5md5"
	release.Size = 1024
	release.Flag = 1

	id, err := release.Save(dal)
	if id < 0 || err != nil {
		return err
	}

	return nil
}
