package policy

import (
	"encoding/json"
	"fmt"
	"github.com/jsli/gtbox/file"
	ota_constant "github.com/jsli/ota/radio/app/constant"
	"github.com/jsli/ota/radio/app/models"
	"io"
	"os"
	"strings"
)

func ParseRequest(request_str string) (*models.UpdateRequest, error) {
	update_request := &models.UpdateRequest{}
	err := json.Unmarshal([]byte(request_str), update_request)
	if err != nil {
		return nil, fmt.Errorf("Illegal format [request] : %s", err)
	}

	model := update_request.Device.Model
	update_request.Device.Model = ConvertModel(model)
	platform := update_request.Device.Platform
	update_request.Device.Platform = ConvertAndroidPlatform(platform)

	return update_request, nil
}

func ParseDtim(dtim_byte []byte) (*DtimInfo, error) {
	//	binary_data, images, err := ParseDtimWithReader(reader)
	images, err := ParseDtimWithByte(dtim_byte)
	if err != nil {
		return nil, err
	}

	//ignore empty items
	for index, d := range images {
		if d == nil || len(d) == 0 {
			images = images[:index]
			break
		}
	}

	count := len(images)
	dtim_info := &DtimInfo{}
	dtim_info.BinaryData = dtim_byte
	cp_image_list := make([]*CpImage, count)
	for index, image := range images {
		if len(image) != 4 {
			return nil, fmt.Errorf("Illegal image information, image's attr count must be 4, NOT %d", len(image))
		}

		cp_image := &CpImage{}
		cp_image.LoadSelf(image)

		err := cp_image.Validate()
		if err != nil {
			return nil, err
		}
		cp_image_list[index] = cp_image
	}

	dtim_info.CpMap = make(map[string]*CpInfo)
	switch count {
	case 6:
		cp_info := &CpInfo{}
		cp_info.ImageMap = make(map[string]*CpImage)
		cp_info.Mode = cp_image_list[3].Mode
		cp_info.Network = cp_image_list[3].Network
		cp_info.Sim = cp_image_list[3].Sim
		cp_info.Version = cp_image_list[3].Version
		cp_info.Prefix = cp_image_list[3].Prefix
		cp_info.ImageMap[ota_constant.KEY_ARBEL] = cp_image_list[3]
		cp_info.ImageMap[ota_constant.KEY_MSA] = cp_image_list[4]
		cp_info.ImageMap[ota_constant.KEY_RFIC] = cp_image_list[5]
		dtim_info.CpMap[cp_info.Mode] = cp_info
		fallthrough
	case 3:
		dtim_info.HasRFIC = true
		cp_info := &CpInfo{}
		cp_info.ImageMap = make(map[string]*CpImage)
		cp_info.Mode = cp_image_list[0].Mode
		cp_info.Network = cp_image_list[0].Network
		cp_info.Sim = cp_image_list[0].Sim
		cp_info.Version = cp_image_list[0].Version
		cp_info.Prefix = cp_image_list[0].Prefix
		cp_info.ImageMap[ota_constant.KEY_ARBEL] = cp_image_list[0]
		cp_info.ImageMap[ota_constant.KEY_MSA] = cp_image_list[1]
		cp_info.ImageMap[ota_constant.KEY_RFIC] = cp_image_list[2]
		dtim_info.CpMap[cp_info.Mode] = cp_info
	case 4:
		cp_info := &CpInfo{}
		cp_info.ImageMap = make(map[string]*CpImage)
		cp_info.Mode = cp_image_list[2].Mode
		cp_info.Network = cp_image_list[2].Network
		cp_info.Sim = cp_image_list[2].Sim
		cp_info.Version = cp_image_list[2].Version
		cp_info.Prefix = cp_image_list[2].Prefix
		cp_info.ImageMap[ota_constant.KEY_ARBEL] = cp_image_list[2]
		cp_info.ImageMap[ota_constant.KEY_MSA] = cp_image_list[3]
		dtim_info.CpMap[cp_info.Mode] = cp_info
		fallthrough
	case 2:
		dtim_info.HasRFIC = false
		cp_info := &CpInfo{}
		cp_info.ImageMap = make(map[string]*CpImage)
		cp_info.Mode = cp_image_list[0].Mode
		cp_info.Network = cp_image_list[0].Network
		cp_info.Sim = cp_image_list[0].Sim
		cp_info.Version = cp_image_list[0].Version
		cp_info.Prefix = cp_image_list[0].Prefix
		cp_info.ImageMap[ota_constant.KEY_ARBEL] = cp_image_list[0]
		cp_info.ImageMap[ota_constant.KEY_MSA] = cp_image_list[1]
		dtim_info.CpMap[cp_info.Mode] = cp_info
	default:
		return nil, fmt.Errorf("Illegal cp information, image count must be 2 or 4, NOT %d", count)
	}

	dtim_info.MD5Dtim = Md5Dtim(dtim_byte)
	return dtim_info, nil
}

func ParseDtimWithReader(reader io.Reader) ([]byte, [][]string, error) {
	buffer := make([]byte, 4096)
	n, err := reader.Read(buffer)
	if err != nil {
		return nil, nil, err
	}
	data, err := ParseDtimWithByte(buffer[:n])
	return buffer[:n], data, err
}

/*---offset:    3*1024
 *<id>|<network>| <sim>  |   <path>
 *ARBI|   LTG   | SINGLE |HLLTE/HLLTE_CP_3.29.000/Seagull/HL_LTG.bin
 *GRBI|   LTG   | SINGLE |HLLTE/HLLTE_CP_3.29.000/TTD_WK_NL_MSA_3.29.000/HL_DL_M09_Y0_AI_SKL_Flash.bin
 *.............
 *---end:       4*1024
 */
func ParseDtimWithFile(path string) ([]byte, [][]string, error) {
	dtim, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}

	buffer := make([]byte, 1024)
	_, err = dtim.ReadAt(buffer, 3072)
	if err != nil {
		return nil, nil, err
	}
	defer dtim.Close()

	data, err := ParseDtimWithByte(buffer)
	return buffer, data, err
}

func ParseDtimWithByte(dtim_byte []byte) ([][]string, error) {
	data := make([][]string, 10)

	if len(dtim_byte) > 1024 {
		dtim_byte = dtim_byte[3072:]
	}
	text := string(dtim_byte)
	if len(text) <= 0 {
		return nil, fmt.Errorf("Empty dtim %s", dtim_byte)
	}

	image_list := strings.Split(text, "\n")

	counter := 0
	for index, image := range image_list {
		if len(image) < 4 {
			return nil, fmt.Errorf("Empty dtim %s", dtim_byte)
		}

		if index := strings.Index(image, "@"); index >= 0 {
			image = image[:index]
		}

		if ValidateImageId(image[:4]) == nil {
			data[index] = make([]string, 0, 10)
			counter += counter
			attrs := strings.Split(image, "|")
			for _, attr := range attrs {
				data[index] = append(data[index], attr)
			}
		}
	}

	return data, nil
}

func DtimToBinary(path string) []byte {
	data, err := file.ReadBinaryFile(path)
	if err != nil {
		return nil
	}
	return data
}
