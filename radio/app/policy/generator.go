package policy

import (
	"encoding/json"
	"fmt"
	"github.com/jsli/gtbox/archive"
	"github.com/jsli/gtbox/file"
	"github.com/jsli/gtbox/pathutil"
	"github.com/jsli/gtbox/sys"
	ota_constant "github.com/jsli/ota/radio/app/constant"
	"github.com/jsli/ota/radio/app/models"
	"math/rand"
	"strings"
	"time"
)

const (
	DEBUG = false
)

func GenerateOtaPackage(dal *models.Dal, task *models.ReleaseCreationTask, root_path string) (*models.RadioOtaRelease, error) {
	//0.parse request from json string which stored in DB.
	update_request, err := ParseRequest(task.UpdateRequest)

	//1.create all paths
	zip_path := fmt.Sprintf("%s%s", root_path, ota_constant.ZIP_DIR_NAME)
	pathutil.MkDir(zip_path)
	template_path := fmt.Sprintf("%s%s", ota_constant.TEMPLATE_ROOT, "HELAN")
	update_pkg_path := fmt.Sprintf("%s%s", root_path, ota_constant.UPDATE_PKG_NAME)
	radio_ota_path := fmt.Sprintf("%s%s", root_path, ota_constant.RADIO_OTA_PACKAGE_NAME)

	if DEBUG {
		fmt.Println(zip_path)
		fmt.Println(template_path)
		fmt.Println(update_pkg_path)
		fmt.Println(radio_ota_path)
	}

	//2.copy template files
	err = file.CopyDir(template_path, zip_path)

	/*------------temp : copy radio.img ---------------*/
	_, err = file.CopyFile("/home/manson/desktop/radio/radio.img", zip_path+"radio.img")

	//3.create Radio.dtim, Radio.img
	image_list := GenerateImageList(update_request)

	//	4. archive all files
	err = archive.ArchiveZipFile(zip_path, update_pkg_path)

	//5.generate update.zip (updatetool + update_pkg.zip)
	params := make([]string, 5)
	params[0] = ota_constant.OTA_CMD_PARAM_PLATFORM_PREFIX + ota_constant.MODEL_TO_PLATFORM[update_request.Device.Model]
	params[1] = ota_constant.OTA_CMD_PARAM_PRODUCT_PREFIX + update_request.Device.Model
	params[2] = ota_constant.OTA_CMD_PARAM_OEM_PREFIX
	params[3] = ota_constant.OTA_CMD_PARAM_OUTPUT_PREFIX + radio_ota_path
	params[4] = ota_constant.OTA_CMD_PARAM_INPUT_PREFIX + update_pkg_path
	generateOtaPackage(ota_constant.OTA_MAKE_CMD, params)

	//6.insert db
	release := &models.RadioOtaRelease{}
	release.Flag = ota_constant.FLAG_AVAILABLE
	release.FingerPrint = GenerateOtaPackageFingerPrint(image_list)
	release.Md5, err = file.Md5SumFile(radio_ota_path)
	if err != nil {
		return nil, err
	}
	release.Size, err = file.GetFileSize(radio_ota_path)
	if err != nil {
		return nil, err
	}

	id, err := release.Save(dal)
	if id < 0 || err != nil {
		return nil, err
	}

	//	7.copy final file to public directory, checksum
	final_dir := fmt.Sprintf("%s%s/", ota_constant.RADIO_OTA_RELEASE_ROOT, release.FingerPrint)
	pathutil.MkDir(final_dir)
	final_path := fmt.Sprintf("%s%s", final_dir, ota_constant.RADIO_OTA_PACKAGE_NAME)
	checksum_path := fmt.Sprintf("%s%s", final_dir, ota_constant.CHECKSUM_TXT_NAME)
	file.CopyFile(radio_ota_path, final_path)
	RecordMd5(final_path, checksum_path)

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
	fp := file.Md5SumString(joined_str)
	return fp
}

func GenerateRandFileName() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%d", r.Int63())
}

func generateRadioImage(cmd string, image_list []string) error {
	fmt.Printf("exec : %s, params = %s", cmd, image_list)
	res, output, err := sys.ExecCmd(cmd, image_list)
	if DEBUG {
		fmt.Println(output)
	}
	if !res || err != nil {
		return fmt.Errorf("%s failed: %s\n\tdetail message: %s\n", cmd, err, output)
	}
	return nil
}

func generateOtaPackage(cmd string, params []string) error {
	res, output, err := sys.ExecCmd(cmd, params)
	if DEBUG {
		fmt.Println(output)
	}
	if !res || err != nil {
		return fmt.Errorf("%s failed: %s\n\tdetail message: %s\n", cmd, err, output)
	}
	return nil
}

func GenerateTestUpdateRequest() (string, *models.UpdateRequest) {
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
	//	fmt.Println(js_str)
	return js_str, update_request
}
