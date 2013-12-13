package policy

import (
	"fmt"
	cp_constant "github.com/jsli/cp_release/constant"
	cp_policy "github.com/jsli/cp_release/policy"
	"github.com/jsli/cp_release/release"
	ota_constant "github.com/jsli/ota/radio/app/constant"
	"github.com/jsli/ota/radio/app/models"
	"github.com/robfig/revel"
	"net/url"
	"reflect"
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

func (ci CpImage) String() string {
	return fmt.Sprintf("CpImage(Id=%s, Network=%s, Sim=%s, Path=%s, Prefix=%s, Version=%s, Mode=%s)",
		ci.Id, ci.Network, ci.Sim, ci.Path, ci.Prefix, ci.Version, ci.Mode)
}

func (ci *CpImage) LoadSelf(attrs []string) error {
	ci.Id = attrs[0]
	ci.Network = attrs[1]
	ci.Sim = attrs[2]
	ci.Path = attrs[3]

	slice := strings.Split(strings.TrimSpace(ci.Path), "/")
	slice = TrimArrayTail(slice)
	if len(slice) < 3 {
		return fmt.Errorf("Parse image path failed. path = %s", ci.Path)
	}
	ci.Mode = slice[1]
	ci.Version = cp_policy.ExtractVersion(slice[2])
	if ci.Version == "" {
		return fmt.Errorf("Parse version from image path failed. path = %s", ci.Path)
	}
	ci.Prefix = strings.TrimSuffix(slice[2], ci.Version)
	return nil
}

func (ci *CpImage) Validate() (err error) {
	//	if err := ValidateSim(ci.Sim); err != nil {
	//		return err
	//	}
	if err := ValidateNetwork(ci.Network); err != nil {
		return err
	}

	dal, err := release.NewDal()
	if err != nil {
		return fmt.Errorf("Validate dtim error")
	}
	defer dal.Close()

	var image interface{}
	switch ci.Id {
	case ota_constant.ID_ARBI, ota_constant.ID_ARB2:
		image, _ = release.FindArbiByPath(dal, ci.Path)
	case ota_constant.ID_GRBI, ota_constant.ID_GRB2:
		image, _ = release.FindGrbiByPath(dal, ci.Path)
	case ota_constant.ID_RFIC, ota_constant.ID_RFI2:
		image, _ = release.FindRficByPath(dal, ci.Path)
	}

	v := reflect.ValueOf(image)
	if v.IsNil() {
		return fmt.Errorf("Illegal image path in dtim: %s", ci.Path)
	}

	return nil
}

type CpInfo struct {
	Mode     string
	Version  string
	Network  string
	Sim      string
	Prefix   string
	ImageMap map[string]*CpImage
}

type DtimInfo struct {
	HasRFIC    bool
	CpMap      map[string]*CpInfo
	BinaryData []byte
	MD5Dtim    string
}

type Validator interface {
	ValidateAndParseRadioDtim(params *revel.Params) (DtimInfo, error)
	ValidateUpdateRequest(params *revel.Params) (*models.UpdateRequest, string, error)
}

type RadioValidator struct {
}

func (v *RadioValidator) ValidateUpdateRequest(params *revel.Params) (*models.UpdateRequest, error) {
	var request_str string = ""
	params.Bind(&request_str, "request")
	if request_str == "" {
		return nil, fmt.Errorf("Illegal param [request]!")
	}

	request_str, err := url.QueryUnescape(request_str)
	if err != nil {
		return nil, fmt.Errorf("Illegal format [request] : %s", err)
	}
	request, err := ParseRequest(request_str)
	return request, err
}

func (v *RadioValidator) ValidateAndParseRadioDtim(params *revel.Params) (*DtimInfo, error) {
	dtim_byte, err := readDtimByte(params)
	if err != nil {
		return nil, err
	}

	return ParseDtim(dtim_byte)
}

func readDtimByte(params *revel.Params) ([]byte, error) {
	fh_arr, ok := params.Files[ota_constant.RADIO_DTIM_NAME]
	if !ok || len(fh_arr) <= 0 {
		return nil, fmt.Errorf("Post request lost file : %s", ota_constant.RADIO_DTIM_NAME)
	}

	input, err := fh_arr[0].Open()
	if err != nil {
		return nil, err
	}
	defer input.Close()

	buffer := make([]byte, 4096)
	n, err := input.Read(buffer)
	if err != nil {
		return nil, err
	}
	return buffer[:n], nil
}

func (v *RadioValidator) CompareRequestAndDtim(request *models.UpdateRequest, dtim_info *DtimInfo) error {
	request_cps_arr := request.Cps
	dtim_cps_map := dtim_info.CpMap

	request_cp_len := len(request_cps_arr)
	dtim_cp_len := len(dtim_cps_map)
	if request_cp_len == 0 {
		return fmt.Errorf("None CP request")
	}

	if request_cp_len == dtim_cp_len {
		return nil
	}

	match := false
	for mode, dtim_cp := range dtim_cps_map {
		found := false
		for _, request_cp := range request_cps_arr {
			if mode == request_cp.Mode {
				found = true
				match = true
				break
			}
		}

		if !found {
			request := models.CpRequest{}
			request.Mode = mode
			request.Version = dtim_cp.Version
			images := make(map[string]string)

			if arbi, ok := dtim_cp.ImageMap[ota_constant.KEY_ARBEL]; ok {
				images[ota_constant.KEY_ARBEL] = arbi.Path
			} else {
				return fmt.Errorf("dtim and request unmatch!!!")
			}

			if grbi, ok := dtim_cp.ImageMap[ota_constant.KEY_MSA]; ok {
				images[ota_constant.KEY_MSA] = grbi.Path
			} else {
				return fmt.Errorf("dtim and request unmatch!!!")
			}

			if dtim_info.HasRFIC {
				if rfic, ok := dtim_cp.ImageMap[ota_constant.KEY_RFIC]; ok {
					images[ota_constant.KEY_RFIC] = rfic.Path
				} else {
					return fmt.Errorf("dtim and request unmatch!!!")
				}
			}

			request.Images = images
			request_cps_arr = append(request_cps_arr, request)
			//			fmt.Println(request)
		}

		request.Cps = request_cps_arr
	}

	if !match {
		return fmt.Errorf("request and dtim MODE unmatched!!!")
	}
	return nil
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

func ValidateImageId(id string) error {
	for _, _id := range ota_constant.IMAGE_ID_LIST {
		if strings.ToUpper(id) == _id {
			return nil
		}
	}
	return fmt.Errorf("Illegal Image Id: %s", id)
}
