package controllers

import (
	"fmt"
	"github.com/jsli/cp_release/release"
	ota_constant "github.com/jsli/ota/radio/app/constant"
	"github.com/jsli/ota/radio/app/models"
	"github.com/jsli/ota/radio/app/policy"
	"github.com/robfig/revel"
	"github.com/robfig/revel/cache"
	"net/http"
	"time"
)

type Radio struct {
	*revel.Controller
}

func (c Radio) Index() revel.Result {
	if checkMaintain() {
		return c.Redirect("/radio/maintain")
	}
	_, j := policy.GenerateTestUpdateRequest()
	return c.RenderJson(j)
}

func (c Radio) Maintain() revel.Result {
	result := models.NewMaintainResult()
	return c.RenderJson(result)
}

func checkMaintain() bool {
	result, found := revel.Config.Bool("maintain")
	if found && result {
		return true
	}
	return false
}

func (c Radio) OtaCreate() revel.Result {
	if checkMaintain() {
		return c.Redirect("/radio/maintain")
	}

	revel.INFO.Println("OtaCreate request: ", c.Request)
	result := models.NewRadioOtaReleaseResult()

	validator := &policy.RadioValidator{}
	dtim_info, err := validator.ValidateAndParseRadioDtim(c.Params)
	if err != nil {
		revel.ERROR.Println("http.StatusBadRequest: ", err)
		c.Response.Status = http.StatusBadRequest
		result.Extra.ErrorCode = ota_constant.ERROR_CODE_DROPPED
		result.Extra.ErrorMessage = fmt.Sprintf("%s", err)
		return c.RenderJson(result)
	}

	update_request, request_json, err := validator.ValidateUpdateRequest(c.Params)
	if err != nil {
		revel.ERROR.Println("http.StatusBadRequest: ", err)
		c.Response.Status = http.StatusBadRequest
		result.Extra.ErrorCode = ota_constant.ERROR_CODE_DROPPED
		result.Extra.ErrorMessage = fmt.Sprintf("%s", err)
		return c.RenderJson(result)
	}

	err = validator.CompareRequestAndDtim(update_request, dtim_info)
	if err != nil {
		revel.ERROR.Println("http.StatusBadRequest: ", err)
		c.Response.Status = http.StatusBadRequest
		result.Extra.ErrorCode = ota_constant.ERROR_CODE_DROPPED
		result.Extra.ErrorMessage = fmt.Sprintf("%s", err)
		return c.RenderJson(result)
	}

	dal, err := models.NewDal()
	if err != nil {
		revel.ERROR.Println("http.StatusInternalServerError: ", err)
		c.Response.Status = http.StatusInternalServerError
		result.Extra.ErrorCode = ota_constant.ERROR_CODE_DROPPED
		return c.RenderJson(result)
	}
	defer dal.Close()

	update_request.Cps = policy.SortCps(update_request)
	sorted_image_list := policy.GenerateImageList(update_request)
	fp := policy.GenerateOtaPackageFingerPrint(sorted_image_list)
	fp = fmt.Sprintf("%s.%s.%s", update_request.Device.Model, update_request.Device.Platform, fp)

	if err := cache.Get(fp, &result); err == nil {
		return c.RenderJson(result)
	}

	radio, err := policy.ProvideRadioRelease(dal, dtim_info, result, fp)
	if err != nil {
		revel.ERROR.Println("http.StatusInternalServerError: ", err)
		c.Response.Status = http.StatusInternalServerError
		result.Extra.ErrorCode = ota_constant.ERROR_CODE_DROPPED
		return c.RenderJson(result)
	}

	if radio == nil {
		task, err := models.FindReleaseCreationTaskByFp(dal, fp)
		if err != nil {
			revel.ERROR.Println("FindReleaseCreationTaskByFp failed: ", err)
		}
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
				revel.ERROR.Println("http.StatusInternalServerError: ", err)
				result.Extra.ErrorCode = ota_constant.ERROR_CODE_CREATE_REQUEST_FAILED
				result.Extra.ErrorMessage = "Create creation task Failed, try later"
			} else {
				revel.INFO.Println("OtaCreate request, create task: ", task.UpdateRequest)
				result.Extra.ErrorCode = ota_constant.ERROR_CODE_NOT_CREATED
				result.Extra.ErrorMessage = "Update package will be created later"
			}
		} else {
			switch task.Flag {
			case ota_constant.FLAG_DISABLE:
				fallthrough
			case ota_constant.FLAG_DROPPED:
				result.Extra.ErrorCode = ota_constant.ERROR_CODE_DROPPED
				result.Extra.ErrorMessage = "Bad creation task, drop it"
			case ota_constant.FLAG_AVAILABLE:
				fallthrough
			case ota_constant.FLAG_INIT:
				result.Extra.ErrorCode = ota_constant.ERROR_CODE_NOT_CREATED
				result.Extra.ErrorMessage = "Update package will be created later"
			case ota_constant.FLAG_CREATING:
				result.Extra.ErrorCode = ota_constant.ERROR_CODE_CREATING
				result.Extra.ErrorMessage = "Update package is creating"
			case ota_constant.FLAG_CREATED:
				result.Extra.ErrorCode = ota_constant.ERROR_CODE_NOERR
				result.Extra.ErrorMessage = nil
			case ota_constant.FLAG_CREATE_FAILED:
				result.Extra.ErrorCode = ota_constant.ERROR_CODE_CREATE_REQUEST_FAILED
				result.Extra.ErrorMessage = "Update package created failed"
			default:
				result.Extra.ErrorCode = ota_constant.ERROR_CODE_NOERR
				result.Extra.ErrorMessage = nil
			}
		}

		c.Response.Status = http.StatusNotFound
		return c.RenderJson(result)
	} else {
		revel.INFO.Println("OtaCreate, find release: ", radio)
		result.Data.Url = fmt.Sprintf("http://%s/static/%s/%s", c.Request.Host, radio.FingerPrint, ota_constant.RADIO_OTA_PACKAGE_NAME)
		result.Data.Md5 = radio.Md5
		result.Data.Size = radio.Size
		result.Data.CreatedTime = policy.FormatTime(radio.CreatedTs)

		cache.Set(fp, result, 60*time.Second)
	}

	return c.RenderJson(result)
}

func (c Radio) Query() revel.Result {
	if checkMaintain() {
		return c.Redirect("/radio/maintain")
	}

	revel.INFO.Println("Query request: ", c.Request)
	validator := &policy.RadioValidator{}
	dtim_info, err := validator.ValidateAndParseRadioDtim(c.Params)
	if err != nil {
		revel.ERROR.Println("http.StatusBadRequest: ", err)
		c.Response.Status = http.StatusBadRequest
		return c.RenderJson(nil)
	}

	dal, err := release.NewDal()
	if err != nil {
		revel.ERROR.Println("http.StatusInternalServerError: ", err)
		c.Response.Status = http.StatusInternalServerError
		return c.RenderJson(nil)
	}
	defer dal.Close()

	var result_cached models.QueryResult
	if err := cache.Get(dtim_info.MD5Dtim, &result_cached); err == nil {
		return c.RenderJson(result_cached)
	}

	result := models.NewQueryResult()
	err = policy.ProvideQueryData(dal, dtim_info, result)
	if err != nil {
		revel.ERROR.Println("http.StatusInternalServerError: ", err)
		c.Response.Status = http.StatusInternalServerError
		result.Extra.ErrorMessage = err
		return c.RenderJson(result)
	} else {
		if result.Data.Available == nil || len(result.Data.Available) == 0 {
			c.Response.Status = http.StatusNotFound
			return c.RenderJson(nil)
		}
	}

	revel.INFO.Println("Query request, result: ", result)
	cache.Set(dtim_info.MD5Dtim, result, 60*time.Second)
	return c.RenderJson(result)
}
