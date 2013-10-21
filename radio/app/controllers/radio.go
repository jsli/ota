package controllers

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jsli/cp_release/release"
	"github.com/jsli/gtbox/file"
	"github.com/jsli/gtbox/ota"
	"github.com/jsli/gtbox/pathutil"
	ota_constant "github.com/jsli/ota/radio/app/constant"
	"github.com/jsli/ota/radio/app/models"
	"github.com/jsli/ota/radio/app/policy"
	"github.com/robfig/revel"
	"net/http"
)

type Radio struct {
	*revel.Controller
}

func (c Radio) Index() revel.Result {
	policy.GenerateTestUpdateRequest()
	return c.Render()
}

func (c Radio) OtaCreate() revel.Result {
	root := fmt.Sprintf("%s%s/", ota_constant.TMP_FILE_ROOT, ota.GenerateRandFileName())
	pathutil.MkDir(root)
	//	defer file.DeleteDir(root)

	validator := &policy.RadioValidator{}
	dtim_info, err := validator.ValidateAndParseRadioDtim(c.Params, root, false)
	if err != nil {
		c.Response.Status = http.StatusBadRequest
		return c.RenderJson(nil)
	}

	update_request, request_json, err := validator.ValidateUpdateRequest(c.Params)
	if err != nil {
		c.Response.Status = http.StatusBadRequest
		return c.RenderJson(nil)
	}

	dal, err := models.NewDal()
	if err != nil {
		c.Response.Status = http.StatusInternalServerError
		return c.RenderJson(nil)
	}
	defer dal.Close()

	update_request.Cps = policy.SortCps(update_request)
	sorted_image_list := policy.GenerateImageList(update_request)
	fp := policy.GenerateOtaPackageFingerPrint(sorted_image_list)

	result := models.NewRadioOtaReleaseResult()
	radio, err := policy.ProvideRadioRelease(dal, dtim_info, result, fp)
	if err != nil {
		c.Response.Status = http.StatusInternalServerError
		result.Extra[ota_constant.KEY_ERROR] = err
		return c.RenderJson(result)
	}

	if radio == nil {
		radio, err = policy.GenerateOtaPackage(dal, dtim_info, update_request, sorted_image_list, request_json, root)
		if err != nil {
			c.Response.Status = http.StatusInternalServerError
			result.Extra[ota_constant.KEY_ERROR] = err
			return c.RenderJson(result)
		}
	}

	if radio != nil {
		result.Data[ota_constant.KEY_URL] = fmt.Sprintf("http://%s/static/%s/%s", c.Request.Host, radio.FingerPrint, ota_constant.RADIO_OTA_PACKAGE_NAME)
		result.Data[ota_constant.KEY_MD5] = radio.Md5
		result.Data[ota_constant.KEY_SIZE] = radio.Size
		result.Extra[ota_constant.KEY_ERROR] = nil
	} else {
		c.Response.Status = http.StatusNotFound
		result.Extra[ota_constant.KEY_ERROR] = nil
		return c.RenderJson(nil)
	}

	return c.RenderJson(result)
}

func (c Radio) Query() revel.Result {
	root := fmt.Sprintf("%s%s/", ota_constant.TMP_FILE_ROOT, ota.GenerateRandFileName())
	pathutil.MkDir(root)
	defer file.DeleteDir(root)

	validator := &policy.RadioValidator{}
	dtim_info, err := validator.ValidateAndParseRadioDtim(c.Params, root, true)
	if err != nil {
		fmt.Println(err)
		c.Response.Status = http.StatusBadRequest
		return c.RenderJson(nil)
	}
	dal, err := release.NewDal()
	if err != nil {
		c.Response.Status = http.StatusInternalServerError
		return c.RenderJson(nil)
	}
	defer dal.Close()

	result := models.NewQueryResult()
	err = policy.ProvideQueryData(dal, dtim_info, result)
	if err != nil {
		c.Response.Status = http.StatusInternalServerError
		result.Extra[ota_constant.KEY_ERROR] = err
		return c.RenderJson(result)
	} else {
		if result.Data == nil || len(result.Data) == 0 {
			c.Response.Status = http.StatusNotFound
			return c.RenderJson(nil)
		} else {
			result.Extra[ota_constant.KEY_ERROR] = nil
		}
	}

	return c.RenderJson(result)
}
