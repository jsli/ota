package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/jsli/cp_release/release"
	ota_constant "github.com/jsli/ota/radio/app/constant"
	"github.com/jsli/ota/radio/app/models"
	"github.com/jsli/ota/radio/app/policy"
	"github.com/robfig/revel"
	"github.com/robfig/revel/cache"
	"time"
)

func init() {
	revel.InterceptMethod(Radio.LogVisitorByIP, revel.BEFORE)
}

type Radio struct {
	App
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

func (c Radio) OtaCreate() revel.Result {
	result := models.NewRadioOtaReleaseResult()
	validator := &policy.RadioValidator{}
	dtim_info, err := validator.ValidateAndParseRadioDtim(c.Params)
	if err != nil {
		return c.Render400(result, err)
	}

	update_request, err := validator.ValidateUpdateRequest(c.Params)
	if err != nil {
		return c.Render400(result, err)
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

	if err := cache.Get(fp, &result); err == nil {
		return c.RenderJson(result)
	}

	radio, err := policy.ProvideRadioRelease(dal, dtim_info, result, fp)
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
		result.Data.Url = fmt.Sprintf("http://%s/static/%s/%s", c.Request.Host, radio.FingerPrint, ota_constant.RADIO_OTA_PACKAGE_NAME)
		result.Data.Md5 = radio.Md5
		result.Data.Size = radio.Size
		result.Data.CreatedTime = policy.FormatTime(radio.CreatedTs)

		cache.Set(fp, result, 60*time.Second)
	}

	return c.RenderJson(result)
}

func (c Radio) Query() revel.Result {
	result := models.NewQueryResult()
	validator := &policy.RadioValidator{}
	dtim_info, err := validator.ValidateAndParseRadioDtim(c.Params)
	if err != nil {
		return c.Render400(result, err)
	}

	if err := cache.Get(dtim_info.MD5Dtim, result); err == nil {
		return c.RenderJson(result)
	}

	dal, err := release.NewDal()
	if err != nil {
		return c.Render500(result, err)
	}
	defer dal.Close()

	err = policy.ProvideQueryData(dal, dtim_info, result)
	if err != nil {
		return c.Render500(result, err)
	} else {
		if result.Data.Available == nil || len(result.Data.Available) == 0 {
			return c.Render404(result, nil)
		}
	}

	cache.Set(dtim_info.MD5Dtim, result, 60*time.Second)
	return c.RenderJson(result)
}
