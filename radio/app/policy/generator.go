package policy

import (
	"encoding/json"
	"fmt"
	cp_constant "github.com/jsli/cp_release/constant"
	"github.com/jsli/gtbox/archive"
	"github.com/jsli/gtbox/file"
	"github.com/jsli/gtbox/ota"
	"github.com/jsli/gtbox/pathutil"
	ota_constant "github.com/jsli/ota/radio/app/constant"
	"github.com/jsli/ota/radio/app/models"
	"strings"
)

func GenerateOtaPackage(dal *models.Dal, dtim_info *DtimInfo, update_request *models.UpdateRequest, image_list []string, json_str string, root_path string) (*models.RadioOtaRelease, error) {
	//1.create new dtim (xxx.rb + radio.dtim)
	//	CreateRadioImage(dtim_info.DtimPath, image_list)
	//2.create radio.img (dtim + cp)
	zip_path := fmt.Sprintf("%s%s", root_path, ota_constant.ZIP_DIR_NAME)
	pathutil.MkDir(zip_path)
//	fmt.Println(zip_path)

	//3.copy template files and generate update_pkg.zip (template + radio.img)
	template_path := fmt.Sprintf("%s%s", ota_constant.TEMPLATE_ROOT, "HELAN")
//	fmt.Println(template_path)

	err := file.CopyDir(template_path, zip_path)
	fmt.Println(err)

	/*------------temp : copy radio.img ---------------*/
	_, err = file.CopyFile("/home/manson/desktop/radio/radio.img", zip_path+"radio.img")
//	fmt.Println(err)

	update_pkg_path := fmt.Sprintf("%s%s", root_path, ota_constant.UPDATE_PKG_NAME)
//	fmt.Println(update_pkg_path)
	err = archive.ArchiveZipFile(zip_path, update_pkg_path)
//	fmt.Println(err)

	//4.generate update.zip (updatetool + update_pkg.zip)
	radio_ota_path := fmt.Sprintf("%s%s", root_path, ota_constant.RADIO_OTA_PACKAGE_NAME)
	cmd_params := make([]string, 5)
	cmd_params[0] = ota_constant.OTA_CMD_PARAM_PLATFORM_PREFIX + ota_constant.MODEL_TO_PLATFORM[update_request.Device.Model]
	cmd_params[1] = ota_constant.OTA_CMD_PARAM_PRODUCT_PREFIX + update_request.Device.Model
	cmd_params[2] = ota_constant.OTA_CMD_PARAM_OEM_PREFIX
	cmd_params[3] = ota_constant.OTA_CMD_PARAM_OUTPUT_PREFIX + radio_ota_path
	cmd_params[4] = ota_constant.OTA_CMD_PARAM_INPUT_PREFIX + update_pkg_path
	ota.GenerateOtaPackage(ota_constant.OTA_PKG_MAKE_CMD, cmd_params)

	//5.insert db
	release := &models.RadioOtaRelease{}
	release.FingerPrint = GenerateOtaPackageFingerPrint(image_list)
	release.Md5 = "md5md5md5md5md5md5md5md5md5md5md5md5md5"
	release.Size = 1024
	release.Flag = ota_constant.AVAILABLE_FLAG
	release.Detail = json_str

//	id, err := release.Save(dal)
//	if id < 0 || err != nil {
//		return nil, err
//	}
	return release, nil
}

func GenerateImageList(update_request *models.UpdateRequest) []string {
	request_cps := update_request.Cps
	image_list := make([]string, 0, 10)
	for _, image := range request_cps {
		image_map := image.Images
		if arbel, ok := image_map[ota_constant.KEY_ARBEL]; ok {
			image_list = append(image_list, arbel)
		}
		if msa, ok := image_map[ota_constant.KEY_MSA]; ok {
			image_list = append(image_list, msa)
		}
		if rfic, ok := image_map[ota_constant.KEY_RFIC]; ok {
			image_list = append(image_list, rfic)
		}
	}
	return image_list
}

func GenerateOtaPackageFingerPrint(image_list []string) string {
	joined_str := strings.Join(image_list, "\n")
	//	fmt.Println(joined_str)

	fp := file.Md5SumString(joined_str)
	//	fmt.Println(fp)
	return fp
}

func CreateRadioImage(dtim_path string, image_list []string) (string, error) {
	cmd := ota_constant.RESIGN_DTIM_CMD
	param_list := make([]string, 0, 10)
	param_list = append(param_list, dtim_path)

	for index, path := range image_list {
		mode := strings.Split(path, "/")[0]
		prefix := cp_constant.MODE_TO_ROOT_PATH[mode]
		image_list[index] = fmt.Sprintf("%s%s", prefix, path)
	}

	param_list = append(param_list, image_list[:4]...)
	err := ota.GenerateRadioImage(cmd, param_list)
	return "", err
}

func GenerateTestUpdateRequest() string {
	update_request := &models.UpdateRequest{}

	device_info := models.DeviceInfo{}
	device_info.Model = "pxa1t88ff_def"
	device_info.MacAddr = "08:11:96:8a:a4:38"
	update_request.Device = device_info

	cps := make([]models.CpRequest, 0, 2)

	hltd := models.CpRequest{}
	hltd.Mode = "HLTD"
	hltd.Version = "2.10.000"
	hltd_images := make(map[string]string)
	hltd_images["ARBEL"] = "LWG/HL_CP_2.40.000/HL_CP/Seagull/HL_LWG_DKB.bin"
	hltd_images["MSA"] = "LWG/HL_CP_2.40.000/HL_MSA_2.40.000/HL_LWG_M09_B0_SKL_Flash.bin"
	hltd_images["RFIC"] = "LWG/HL_CP_2.40.000/RFIC/1920_FF/Skylark_LWG.bin"
	hltd.Images = hltd_images
	cps = append(cps, hltd)

	hltd_dsds := models.CpRequest{}
	hltd_dsds.Mode = "HLTD_DSDS"
	hltd_dsds.Version = "3.10.000"
	hltd_dsds_images := make(map[string]string)
	hltd_dsds_images["ARBEL"] = "LTG/HL_CP_3.40.000/HL_CP/Seagull/HL_LTG_DL.bin"
	hltd_dsds_images["MSA"] = "LTG/HL_CP_3.40.000/HL_MSA_3.40.000/HL_DL_M09_Y0_AI_SKL_Flash.bin"
	hltd_dsds_images["RFIC"] = "LTG/HL_CP_3.40.000/RFIC/1920_FF/Skylark_LTG.bin"
	hltd_dsds.Images = hltd_dsds_images
	cps = append(cps, hltd_dsds)

	update_request.Cps = cps

	js_byte, err := json.Marshal(update_request)
	if err != nil {
		panic(err)
	}
	js_str := string(js_byte)
	fmt.Println(js_str)
	return js_str
}
