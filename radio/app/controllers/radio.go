package controllers

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jsli/cp_release/release"
	"github.com/jsli/ota/radio/app/models"
	"github.com/jsli/ota/radio/app/policy"
	"github.com/robfig/revel"
	"net/http"
)

type Radio struct {
	*revel.Controller
}

func (c Radio) Index() revel.Result {
	return c.Render()
}

func (c Radio) OtaCreate() revel.Result {
	validator := &policy.RadioValidator{}
	parsedParams, err := validator.PostValidate(c.Params)

	if err != nil {
		fmt.Println(err)
		c.Response.Status = http.StatusBadRequest
		return c.RenderJson(nil)
	}

	dal, err := models.NewDal()
	if err != nil {
		c.Response.Status = http.StatusInternalServerError
		return c.RenderJson(nil)
	}
	defer dal.Close()

	var versions string
	c.Params.Bind(&versions, "version")
	fmt.Println(versions)

	result := models.NewRadioOtaReleaseResult()
	err = policy.ProvideRadioRelease(dal, parsedParams, result, versions)
	if err != nil {
		c.Response.Status = http.StatusInternalServerError
		result.Extra["error"] = err
		return c.RenderJson(result)
	} else {
		if len(result.Data) == 0 {
			err = policy.GenerateOtaPackage(dal, parsedParams, versions)
			if err != nil {
				c.Response.Status = http.StatusInternalServerError
				result.Extra["error"] = err
				return c.RenderJson(result)
			}
			err = policy.ProvideRadioRelease(dal, parsedParams, result, versions)
			if err != nil {
				c.Response.Status = http.StatusInternalServerError
				result.Extra["error"] = err
				return c.RenderJson(result)
			} else {
				result.Data["url"] = fmt.Sprintf("http://%s/static/%s", c.Request.Host, result.Data["url"])
				return c.RenderJson(result)
			}
		} else {
			result.Data["url"] = fmt.Sprintf("http://%s/static/%s", c.Request.Host, result.Data["url"])
			return c.RenderJson(result)
		}
	}

	return c.RenderJson(result)
}

func (c Radio) Query() revel.Result {
	validator := &policy.RadioValidator{}
	parsedParams, err := validator.PostValidate(c.Params)
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
	err = policy.ProvideQueryData(dal, parsedParams, result)
	if err != nil {
		c.Response.Status = http.StatusInternalServerError
		result.Extra["error"] = err
		return c.RenderJson(result)
	} else {
		if result.Data == nil || len(result.Data) == 0 {
			c.Response.Status = http.StatusNotFound
			return c.RenderJson(nil)
		} else {
			result.Extra["error"] = nil
		}
	}

	return c.RenderJson(result)
}
