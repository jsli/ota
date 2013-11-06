package controllers

import (
	"fmt"
	"github.com/jsli/cp_release/release"
	ota_constant "github.com/jsli/ota/radio/app/constant"
	"github.com/jsli/ota/radio/app/models"
	"github.com/jsli/ota/radio/app/policy"
	"github.com/robfig/revel"
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
	validator := &policy.RadioValidator{}
	dtim_info, err := validator.ValidateAndParseRadioDtim(c.Params)
	if err != nil {
		revel.ERROR.Println("http.StatusBadRequest: ", err)
		c.Response.Status = http.StatusBadRequest
		return c.RenderJson(nil)
	}

	update_request, request_json, err := validator.ValidateUpdateRequest(c.Params)
	if err != nil {
		revel.ERROR.Println("http.StatusBadRequest: ", err)
		c.Response.Status = http.StatusBadRequest
		return c.RenderJson(nil)
	}

	dal, err := models.NewDal()
	if err != nil {
		revel.ERROR.Println("http.StatusInternalServerError: ", err)
		c.Response.Status = http.StatusInternalServerError
		return c.RenderJson(nil)
	}
	defer dal.Close()

	update_request.Cps = policy.SortCps(update_request)
	sorted_image_list := policy.GenerateImageList(update_request)
	fp := policy.GenerateOtaPackageFingerPrint(sorted_image_list)
	fp = fmt.Sprintf("%s.%s.%s", update_request.Device.Model, update_request.Device.Platform, fp)

	result := models.NewRadioOtaReleaseResult()
	radio, err := policy.ProvideRadioRelease(dal, dtim_info, result, fp)
	if err != nil {
		revel.ERROR.Println("http.StatusInternalServerError: ", err)
		c.Response.Status = http.StatusInternalServerError
		result.Extra.ErrorMessage = err
		return c.RenderJson(result)
	}

	if radio == nil {
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
			result.Extra.ErrorMessage = "Duplicated creation task, wait"
		} else {
			revel.INFO.Println("OtaCreate request, create task: ", task.UpdateRequest)
			result.Extra.ErrorMessage = "Create creation task, try later"
		}

		c.Response.Status = http.StatusNotFound
		return c.RenderJson(result)
	} else {
		revel.INFO.Println("OtaCreate, find release: ", radio)
		result.Data.Url = fmt.Sprintf("http://%s/static/%s/%s", c.Request.Host, radio.FingerPrint, ota_constant.RADIO_OTA_PACKAGE_NAME)
		result.Data.Md5 = radio.Md5
		result.Data.Size = radio.Size
		result.Data.CreatedTime = policy.FormatTime(radio.CreatedTs)
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
	return c.RenderJson(result)
}
