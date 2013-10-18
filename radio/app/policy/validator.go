package policy

import (
	"fmt"
	cp_constant "github.com/jsli/cp_release/constant"
	cp_policy "github.com/jsli/cp_release/policy"
	"github.com/jsli/gtbox/file"
	"github.com/jsli/gtbox/ota"
	ota_constant "github.com/jsli/ota/radio/app/constant"
	"github.com/robfig/revel"
	"os"
	"strings"
)

const ()

/* FMT:
 *<id>|<network>| <sim>  |   <path>
 *ARBI|   LTG   | SINGLE |HLLTE/HLLTE_CP_3.29.000/Seagull/HL_LTG.bin
 *GRBI|   LTG   | SINGLE |HLLTE/HLLTE_CP_3.29.000/TTD_WK_NL_MSA_3.29.000/HL_DL_M09_Y0_AI_SKL_Flash.bin
 */
type CpImage struct {
	Id      string
	Network string
	Sim     string
	Path    string
	Prefix  string
	Version string
	Mode    string
}

func (ci *CpImage) LoadSelf(attrs []string) {
	ci.Id = attrs[0]
	ci.Network = attrs[1]
	ci.Sim = attrs[2]
	ci.Path = attrs[3]

	slice := strings.Split(attrs[3], "/")
	ci.Mode = slice[0]
	ci.Version = cp_policy.ExtractVersion(slice[1])
	ci.Prefix = strings.TrimSuffix(slice[1], ci.Version)
}

func (ci *CpImage) Validate() (err error) {
	err = ValidateSim(ci.Sim)
	err = ValidateNetwork(ci.Network)
	return err
}

type CpInfo struct {
	Mode     string
	Version  string
	Network  string
	Sim      string
	Prefix   string
	ImageMap map[string]*CpImage
}

type ParsedParams struct {
	DtimPath string
	Type     string
	HasRFIC  bool
	CpMap    map[string]*CpInfo
}

type Validator interface {
	Validate(params *revel.Params) (ParsedParams, error)
}

type RadioValidator struct {
}

func (v *RadioValidator) PostValidate(params *revel.Params) (*ParsedParams, error) {
	fh_arr, ok := params.Files[ota_constant.RADIO_DTIM_NAME]
	if !ok || len(fh_arr) <= 0 {
		return nil, fmt.Errorf("Post request lost file : %s", ota_constant.RADIO_DTIM_NAME)
	}

	input, err := fh_arr[0].Open()
	if err != nil {
		return nil, err
	}
	defer input.Close()
	rand_file_path := fmt.Sprintf("%s%s", ota_constant.DTIM_UPLOAD_ROOT, randFileName())
	err = file.WriteReader2File(input, rand_file_path)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(rand_file_path)
	if err != nil {
		return nil, err
	}
	if fi.Size() != ota_constant.RADIO_DTIM_SIZE {
		return nil, fmt.Errorf("%s uploaded size is too big: %d. Need %d only.",
			ota_constant.RADIO_DTIM_NAME, fi.Size(), ota_constant.RADIO_DTIM_SIZE)
	}

	//	images, err := ota.ParseDtim(rand_file_path)
	images, err := ota_constant.TestDataLTER, nil
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

	parsedParams := &ParsedParams{}
	parsedParams.DtimPath = rand_file_path

	count := len(images)
	switch count {
	case 2:
		parsedParams.HasRFIC = false
		parsedParams.Type = ota_constant.TYPE_SINGLE
	case 3:
		parsedParams.HasRFIC = true
		parsedParams.Type = ota_constant.TYPE_SINGLE_RFIC
	case 4:
		parsedParams.HasRFIC = false
		parsedParams.Type = ota_constant.TYPE_DSDS
	case 6:
		parsedParams.HasRFIC = true
		parsedParams.Type = ota_constant.TYPE_DSDS_RFIC
	default:
		return nil, fmt.Errorf("Illegal cp information from %s, image count must be 2 or 4, NOT %d", rand_file_path, count)
	}

	cp_image_list := make([]*CpImage, count)
	for index, image := range images {
		if len(image) != 4 {
			return nil, fmt.Errorf("Illegal image information from %s, image's attr count must be 4, NOT %d", rand_file_path, len(image))
		}

		cp_image := &CpImage{}
		cp_image.LoadSelf(image)

		err := cp_image.Validate()
		if err != nil {
			return nil, err
		}
		cp_image_list[index] = cp_image
	}

	parsedParams.CpMap = make(map[string]*CpInfo)
	switch parsedParams.Type {
	case ota_constant.TYPE_DSDS_RFIC:
		cp_info := &CpInfo{}
		cp_info.ImageMap = make(map[string]*CpImage)
		cp_info.Mode = cp_image_list[3].Mode
		cp_info.Network = cp_image_list[3].Network
		cp_info.Sim = cp_image_list[3].Sim
		cp_info.Version = cp_image_list[3].Version
		cp_info.Prefix = cp_image_list[3].Prefix
		cp_info.ImageMap[ota_constant.ID_ARB2] = cp_image_list[3]
		cp_info.ImageMap[ota_constant.ID_GRB2] = cp_image_list[4]
		cp_info.ImageMap[ota_constant.ID_RFI2] = cp_image_list[5]
		parsedParams.CpMap[ota_constant.TYPE_DSDS_RFIC] = cp_info
		fallthrough
	case ota_constant.TYPE_SINGLE_RFIC:
		cp_info := &CpInfo{}
		cp_info.ImageMap = make(map[string]*CpImage)
		cp_info.Mode = cp_image_list[0].Mode
		cp_info.Network = cp_image_list[0].Network
		cp_info.Sim = cp_image_list[0].Sim
		cp_info.Version = cp_image_list[0].Version
		cp_info.Prefix = cp_image_list[0].Prefix
		cp_info.ImageMap[ota_constant.ID_ARBI] = cp_image_list[0]
		cp_info.ImageMap[ota_constant.ID_GRBI] = cp_image_list[1]
		cp_info.ImageMap[ota_constant.ID_RFIC] = cp_image_list[2]
		parsedParams.CpMap[ota_constant.TYPE_SINGLE_RFIC] = cp_info
	case ota_constant.TYPE_DSDS:
		cp_info := &CpInfo{}
		cp_info.ImageMap = make(map[string]*CpImage)
		cp_info.Mode = cp_image_list[2].Mode
		cp_info.Network = cp_image_list[2].Network
		cp_info.Sim = cp_image_list[2].Sim
		cp_info.Version = cp_image_list[2].Version
		cp_info.Prefix = cp_image_list[2].Prefix
		cp_info.ImageMap[ota_constant.ID_ARB2] = cp_image_list[2]
		cp_info.ImageMap[ota_constant.ID_GRB2] = cp_image_list[3]
		parsedParams.CpMap[ota_constant.TYPE_DSDS] = cp_info
		fallthrough
	case ota_constant.TYPE_SINGLE:
		cp_info := &CpInfo{}
		cp_info.ImageMap = make(map[string]*CpImage)
		cp_info.Mode = cp_image_list[0].Mode
		cp_info.Network = cp_image_list[0].Network
		cp_info.Sim = cp_image_list[0].Sim
		cp_info.Version = cp_image_list[0].Version
		cp_info.Prefix = cp_image_list[0].Prefix
		cp_info.ImageMap[ota_constant.ID_ARBI] = cp_image_list[0]
		cp_info.ImageMap[ota_constant.ID_GRBI] = cp_image_list[1]
		parsedParams.CpMap[ota_constant.TYPE_SINGLE] = cp_info
	}
	return parsedParams, nil
}

func ValidateNetwork(network string) error {
	for _, _network := range cp_constant.NETWORK_LIST {
		if strings.ToUpper(network) == _network {
			return nil
		}
	}
	return fmt.Errorf("Illegal NETWORK: %s", network)
}

func ValidateSim(sim string) error {
	for _, _sim := range cp_constant.SIM_LIST {
		if strings.ToUpper(sim) == _sim {
			return nil
		}
	}
	return fmt.Errorf("Illegal SIM: %s", sim)
}

func randFileName() string {
	rand := ota.GenerateRandFileName()
	return fmt.Sprintf("%s.%s", ota_constant.RADIO_DTIM_NAME, rand)
}
