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
	"github.com/robfig/revel"
	"math/rand"
	"strings"
	"time"
)

func GenerateOtaPackage(dal *models.Dal, task *models.ReleaseCreationTask, root_path string) (*models.RadioOtaRelease, error) {
	//0.parse request from json string which stored in DB.
	update_request, err := ParseRequest(task.UpdateRequest)
	if err != nil {
		return nil, err
	}

	//1.create all paths
	zip_path := fmt.Sprintf("%s%s", root_path, ota_constant.ZIP_DIR_NAME)
	pathutil.MkDir(zip_path)
	template_path := ota_constant.MODEL_TO_TEMPLATE[update_request.Device.Model]
	update_pkg_path := fmt.Sprintf("%s%s", root_path, ota_constant.UPDATE_PKG_NAME)
	radio_ota_path := fmt.Sprintf("%s%s", root_path, ota_constant.RADIO_OTA_PACKAGE_NAME)
	radio_dtim_path := fmt.Sprintf("%s%s", root_path, ota_constant.RADIO_DTIM_NAME)
	radio_image_path := fmt.Sprintf("%s%s", root_path, ota_constant.RADIO_IMAGE_NAME)

	//2.copy template files
	err = file.CopyDir(template_path, zip_path)
	if err != nil {
		return nil, err
	}

	//3.create Radio.dtim, Radio.img
	err = file.WriteBytes2File(task.Data, radio_dtim_path)
	if err != nil {
		return nil, err
	}

	image_list_final := make([]string, 0, 5)
	image_list := GenerateImageList(update_request)
	for _, image_rel_path := range image_list {
		dest_path := root_path + image_rel_path
		_, err := file.CopyFile(ota_constant.CP_SERVER_MIRROR_ROOT+image_rel_path, dest_path)
		if err != nil {
			return nil, err
		}
		image_list_final = append(image_list_final, dest_path)
	}

	err = gzipCpImage(image_list_final)
	if err != nil {
		revel.INFO.Printf("zip image error : %s\n", err)
		return nil, err
	}
	for _, image_path := range image_list_final {
		file.CopyFile(image_path+".gz", image_path)
	}

	err = generateRadioImage(radio_dtim_path, radio_image_path, image_list_final)
	if err != nil {
		return nil, err
	}
	_, err = file.CopyFile(radio_image_path, fmt.Sprintf("%s%s", zip_path, ota_constant.RADIO_IMAGE_NAME))

	//	4. archive all files
	err = archive.ArchiveZipFile(zip_path, update_pkg_path)
	if err != nil {
		return nil, err
	}

	//5.generate update.zip (updatetool + update_pkg.zip)
	params := make([]string, 5)
	params[0] = ota_constant.OTA_CMD_PARAM_PLATFORM_PREFIX + update_request.Device.Platform
	params[1] = ota_constant.OTA_CMD_PARAM_PRODUCT_PREFIX + update_request.Device.Model
	params[2] = ota_constant.OTA_CMD_PARAM_OEM_PREFIX
	params[3] = ota_constant.OTA_CMD_PARAM_OUTPUT_PREFIX + radio_ota_path
	params[4] = ota_constant.OTA_CMD_PARAM_INPUT_PREFIX + update_pkg_path
	err = generateOtaPackage(params)
	if err != nil {
		return nil, err
	}

	//6.insert db
	release := &models.RadioOtaRelease{}
	release.Flag = ota_constant.FLAG_AVAILABLE
	release.FingerPrint = task.FingerPrint
	release.ReleaseNote = "empty"
	release.CreatedTs = time.Now().Unix()
	release.ModifiedTs = release.CreatedTs
	release.Model = update_request.Device.Model
	release.Platform = update_request.Device.Platform
	release.Md5, err = file.Md5SumFile(radio_ota_path)
	if err != nil {
		return nil, err
	}
	release.Size, err = file.GetFileSize(radio_ota_path)
	if err != nil {
		return nil, err
	}

	release.Delete(dal)
	file.DeleteDir(fmt.Sprintf("%s%s", ota_constant.RADIO_OTA_RELEASE_ROOT, release.FingerPrint))
	id, err := release.Save(dal)
	if id < 0 || err != nil {
		return nil, err
	}
	release.Id = id

	//	7.copy final file to public directory, checksum
	final_dir := fmt.Sprintf("%s%s/", ota_constant.RADIO_OTA_RELEASE_ROOT, release.FingerPrint)
	pathutil.MkDir(final_dir)
	final_path := fmt.Sprintf("%s%s", final_dir, ota_constant.RADIO_OTA_PACKAGE_NAME)
	checksum_path := fmt.Sprintf("%s%s", final_dir, ota_constant.CHECKSUM_TXT_NAME)
	_, err = file.CopyFile(radio_ota_path, final_path)
	if err != nil {
		return nil, err
	}
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

func generateOtaPackage(params []string) error {
	res, output, err := sys.ExecCmd(ota_constant.OTA_MAKE_CMD, params)
	revel.INFO.Println("generate ota package: \n", output)
	if !res || err != nil {
		return fmt.Errorf("%s failed: %s\n\tdetail message: %s\n", ota_constant.OTA_MAKE_CMD, err, output)
	}
	return nil
}

func gzipCpImage(path_list []string) error {
	for _, path := range path_list {
		params := make([]string, 0, 5)
		params = append(params, ota_constant.GZIP_CMD_PARAMS...)
		params = append(params, path)
		res, output, err := sys.ExecCmd(ota_constant.GZIP_CMD_NAME, params)
		if !res || err != nil {
			return fmt.Errorf("%s failed: %s\n\tdetail message: %s\n", ota_constant.GZIP_CMD_PARAMS, err, output)
		}
	}

	return nil
}

func generateRadioImage(radio_dtim_path string, radio_image_path string, image_list []string) error {
	params := make([]string, 0, 5)
	params = append(params, radio_dtim_path)
	params = append(params, radio_image_path)
	params = append(params, image_list...)
	res, output, err := sys.ExecCmd(ota_constant.RESIGN_DTIM_CMD, params)
	if !res || err != nil {
		return fmt.Errorf("%s failed: %s\n\tdetail message: %s\n", ota_constant.RESIGN_DTIM_CMD, err, output)
	}
	return nil
}

func GenerateTestUpdateRequest() (string, *models.UpdateRequest) {
	update_request := &models.UpdateRequest{}

	device_info := models.DeviceInfo{}
	device_info.Model = "PXA1920_FF_V10"
	device_info.Platform = "4.3"
	update_request.Device = device_info

	cps := make([]models.CpRequest, 0, 2)

	hltd := models.CpRequest{}
	hltd.Mode = "LWG"
	hltd.Version = "2.41.000"
	hltd_images := make(map[string]string)
	hltd_images["ARBEL"] = "LTE/LWG/HL_CP_2.41.000/HL_CP/Seagull/HL_LWG_DKB.bin"
	hltd_images["MSA"] = "LTE/LWG/HL_CP_2.41.000/HL_MSA_2.41.000/HL_LWG_M09_B0_SKL_Flash.bin"
	hltd_images["RFIC"] = "LTE/LWG/HL_CP_2.41.000/RFIC/1920_FF/Skylark_LWG.bin"
	hltd.Images = hltd_images
	cps = append(cps, hltd)

	//	hltd_dsds := models.CpRequest{}
	//	hltd_dsds.Mode = "LTG"
	//	hltd_dsds.Version = "3.41.000"
	//	hltd_dsds_images := make(map[string]string)
	//	hltd_dsds_images["ARBEL"] = "LTE/LTG/HL_CP_3.41.000/HL_CP/Seagull/HL_LTG_DL.bin"
	//	hltd_dsds_images["MSA"] = "LTE/LTG/HL_CP_3.41.000/HL_MSA_3.41.000/HL_DL_M09_Y0_AI_SKL_Flash.bin"
	//	hltd_dsds_images["RFIC"] = "LTE/LTG/HL_CP_3.41.000/RFIC/1920_FF/Skylark_LTG.bin"
	//	hltd_dsds.Images = hltd_dsds_images
	//	cps = append(cps, hltd_dsds)

	update_request.Cps = cps

	js_byte, err := json.Marshal(update_request)
	if err != nil {
		panic(err)
	}
	js_str := string(js_byte)
	fmt.Println(js_str)
	return js_str, update_request
}
