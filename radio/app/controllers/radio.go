package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jsli/cp_release/release"
	ota_constant "github.com/jsli/ota/radio/app/constant"
	"github.com/jsli/ota/radio/app/models"
	"github.com/jsli/ota/radio/app/policy"
	"github.com/robfig/revel"
)

func init() {
	//revel.InterceptMethod(Radio.LogVisitorByIP, revel.BEFORE)
	revel.InterceptMethod((*Radio).Prepare, revel.BEFORE)
}

type Radio struct {
	App
	ApiVersion string
	Provider   policy.ContentProvider
}

func (c Radio) Index() revel.Result {
	_, j := policy.GenerateTestUpdateRequest()
	return c.RenderJson(j)
}

func (c Radio) LogVisitorByIP() revel.Result {
	rdal, err := models.NewRedisDal()
	if err != nil {
		return nil
	}
	defer rdal.Close()
	ip := policy.FilterIp(c.Request.RemoteAddr)
	rdal.Incr(ip)
	return nil
}

func (c *Radio) Prepare() revel.Result {
	if c.Request.Method == "GET" {
		return nil
	}

	var api_version string = ""
	c.Params.Bind(&api_version, ota_constant.REQUEST_PARAM_APIVERSION)
	c.ApiVersion = api_version

	switch api_version {
	case ota_constant.API_VERSION_1_0:
		c.Provider = new(policy.ContentProviderV1)
	case ota_constant.API_VERSION_2_0:
		c.Provider = new(policy.ContentProviderV2)
	default:
		c.ApiVersion = ota_constant.API_VERSION_1_0
		c.Provider = new(policy.ContentProviderV1)
	}
	return nil
}

func (c Radio) OtaCreate() revel.Result {
	result := models.NewResult()
	validator := &policy.RadioValidator{}
	dtim_info, err := validator.ValidateAndParseRadioDtim(c.Params)
	if err != nil {
		return c.Render400(result, err)
	}

	update_request, err := validator.ValidateUpdateRequest(c.Params)
	if err != nil {
		return c.Render400WithCode(result, ota_constant.ERROR_CODE_INVALIDATED_REQUEST,
			fmt.Sprintf(ota_constant.ERROR_MSG_NO_ILLEGAL_REQUEST, update_request))
	}

	err = validator.CompareRequestAndDtim(update_request, dtim_info)
	if err != nil {
		return c.Render400(result, err)
	}

	update_request.Cps = policy.SortCps(update_request)
	request_json_byte, err := json.Marshal(update_request)
	if err != nil {
		return c.Render400(result, err)
	}
	request_json := string(request_json_byte)

	dal, err := models.NewDal()
	if err != nil {
		return c.Render500(result, err)
	}
	defer dal.Close()

	sorted_image_list := policy.GenerateImageList(update_request)
	fp := policy.GenerateOtaPackageFingerPrint(sorted_image_list)
	fp = fmt.Sprintf("%s.%s.%s", update_request.Device.Model, update_request.Device.Platform, fp)

	radio, err := c.Provider.ProvideRadioRelease(dal, dtim_info, fp)
	if err != nil {
		return c.Render500(result, err)
	}

	if radio == nil {
		task, err := models.FindReleaseCreationTaskByFp(dal, fp)
		if err != nil {
		}
		var (
			err_code int
			err_msg  string
		)
		if task == nil {
			task := &models.ReleaseCreationTask{}
			task.Flag = ota_constant.FLAG_INIT
			task.UpdateRequest = request_json
			task.Data = dtim_info.BinaryData
			task.Model = update_request.Device.Model
			task.Platform = update_request.Device.Platform
			task.FingerPrint = fp
			task.CreatedTs = time.Now().Unix()
			task.ModifiedTs = task.CreatedTs

			id, err := task.Save(dal)
			if id < 0 || err != nil {
				err_code = ota_constant.ERROR_CODE_CREATE_REQUEST_FAILED
				err_msg = "Create creation task Failed, try later"
			} else {
				err_code = ota_constant.ERROR_CODE_NOT_CREATED
				err_msg = "Update package will be created later"
			}
		} else {
			switch task.Flag {
			case ota_constant.FLAG_DISABLE:
				fallthrough
			case ota_constant.FLAG_DROPPED:
				err_code = ota_constant.ERROR_CODE_DROPPED
				err_msg = "Bad creation task, drop it"
			case ota_constant.FLAG_AVAILABLE:
				fallthrough
			case ota_constant.FLAG_INIT:
				err_code = ota_constant.ERROR_CODE_NOT_CREATED
				err_msg = "Update package will be created later"
			case ota_constant.FLAG_CREATING:
				err_code = ota_constant.ERROR_CODE_CREATING
				err_msg = "Update package is creating"
			case ota_constant.FLAG_CREATED:
				err_code = ota_constant.ERROR_CODE_NOERR
				err_msg = ""
			case ota_constant.FLAG_CREATE_FAILED:
				err_code = ota_constant.ERROR_CODE_CREATE_REQUEST_FAILED
				err_msg = "Update package created failed"
			default:
				err_code = ota_constant.ERROR_CODE_NOERR
				err_msg = ""
			}
		}
		return c.Render404WithCode(result, err_code, err_msg)
	} else {
		c.LogVisitorByIP()
		data := models.ReleaseResultData{}
		data.Url = fmt.Sprintf("http://%s/static/%s/%s", c.Request.Host, radio.FingerPrint, ota_constant.RADIO_OTA_PACKAGE_NAME)
		data.Md5 = radio.Md5
		data.Size = radio.Size
		data.CreatedTime = policy.FormatTime(radio.CreatedTs)
		result.SetData(data)
	}

	return c.RenderJson(result)
}

func (c Radio) Query() revel.Result {
	result := models.NewResult()
	validator := &policy.RadioValidator{}
	dtim_info, err := validator.ValidateAndParseRadioDtim(c.Params)
	if err != nil {
		return c.Render400WithCode(result, ota_constant.ERROR_CODE_INVALIDATED_DTIM, fmt.Sprintf("%s", err))
	}

	dal, err := release.NewDal()
	if err != nil {
		return c.Render500(result, err)
	}
	defer dal.Close()

	err = c.Provider.ProvideQueryData(dal, dtim_info, result)
	if err != nil {
		return c.Render404WithCode(result, ota_constant.ERROR_CODE_NO_AVAILABLE_UPDATE, fmt.Sprintf("%s", err))
	}

	return c.RenderJson(result)
}
