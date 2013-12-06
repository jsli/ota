package controllers

import (
	"fmt"
	ota_constant "github.com/jsli/ota/radio/app/constant"
	"github.com/jsli/ota/radio/app/models"
	"github.com/robfig/revel"
	"net/http"
)

func init() {
	revel.InterceptMethod(App.Maintenance, revel.BEFORE)
}

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Apis() revel.Result {
	result := models.NewApiResult()
	radio_api := models.API{}
	radio_api["query"] = "/radio/query"
	radio_api["ota_create"] = "/radio/ota/create"
	result.Radio = radio_api
	return c.RenderJson(result)
}

func (c App) Maintenance() revel.Result {
	result, found := revel.Config.Bool("maintenance")
	if found && result {
		result := models.NewMaintainResult()
		return c.RenderJson(result)
	}
	return nil
}

func (c App) RenderError(err_setter models.ErrorSetter, status int, error_code int, error_message string) revel.Result {
	c.Response.Status = status
	err_setter.SetErrorCode(error_code)
	err_setter.SetErrorMessage(error_message)
	return c.RenderJson(err_setter)
}

func (c App) Render400(err_setter models.ErrorSetter, err error) revel.Result {
	return c.RenderError(err_setter, http.StatusBadRequest,
		ota_constant.ERROR_CODE_DROPPED, fmt.Sprintf("%s", err))
}

func (c App) Render400WithCode(err_setter models.ErrorSetter, code int, err string) revel.Result {
	return c.RenderStatusAndCode(err_setter, http.StatusBadRequest, code, err)
}

func (c App) Render500(err_setter models.ErrorSetter, err error) revel.Result {
	return c.RenderError(err_setter, http.StatusInternalServerError,
		ota_constant.ERROR_CODE_DROPPED, fmt.Sprintf("%s", err))
}

func (c App) Render404(err_setter models.ErrorSetter, err error) revel.Result {
	return c.RenderError(err_setter, http.StatusNotFound,
		ota_constant.ERROR_CODE_DROPPED, fmt.Sprintf("%s", err))
}

func (c App) Render404WithCode(err_setter models.ErrorSetter, code int, err string) revel.Result {
	return c.RenderStatusAndCode(err_setter, http.StatusNotFound, code, err)
}

func (c App) RenderStatusAndCode(err_setter models.ErrorSetter, status int, code int, err string) revel.Result {
	return c.RenderError(err_setter, status, code, err)
}
