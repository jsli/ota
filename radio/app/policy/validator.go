package policy

import (
	"fmt"
	cp_constant "github.com/jsli/cp_release/constant"
	cp_policy "github.com/jsli/cp_release/policy"
	ota_constant "github.com/jsli/ota/radio/app/constant"
	"github.com/jsli/ota/radio/app/models"
	"github.com/robfig/revel"
	"net/url"
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

type DtimInfo struct {
	HasRFIC    bool
	CpMap      map[string]*CpInfo
	BinaryData []byte
}

type Validator interface {
	ValidateAndParseRadioDtim(params *revel.Params) (DtimInfo, error)
	ValidateUpdateRequest(params *revel.Params) (*models.UpdateRequest, string, error)
}

type RadioValidator struct {
}

func (v *RadioValidator) ValidateUpdateRequest(params *revel.Params) (*models.UpdateRequest, string, error) {
	var request_str string = ""
	params.Bind(&request_str, "request")
	if request_str == "" {
		return nil, request_str, fmt.Errorf("Illegal param [request]!")
	}

	request_str, err := url.QueryUnescape(request_str)
	if err != nil {
		return nil, "", fmt.Errorf("Illegal format [request] : %s", err)
	}
	request, err := ParseRequest(request_str)
	return request, request_str, err
}

func (v *RadioValidator) ValidateAndParseRadioDtim(params *revel.Params) (*DtimInfo, error) {
	fh_arr, ok := params.Files[ota_constant.RADIO_DTIM_NAME]
	if !ok || len(fh_arr) <= 0 {
		return nil, fmt.Errorf("Post request lost file : %s", ota_constant.RADIO_DTIM_NAME)
	}

	input, err := fh_arr[0].Open()
	if err != nil {
		return nil, err
	}
	defer input.Close()

	return ParseDtim(input)
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
