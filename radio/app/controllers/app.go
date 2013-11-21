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

func (c App) Maintenance() revel.Result {
	result, found := revel.Config.Bool("maintenance")
	if found && result {
		result := models.NewMaintainResult()
		return c.RenderJson(result)
	}
	return nil
}

func (c App) RenderError(err_if models.ErrorInterface, status int, error_code int, error_message string) revel.Result {
	c.Response.Status = status
	err_if.SetErrorCode(error_code)
	err_if.SetErrorMessage(error_message)
	return c.RenderJson(err_if)
}

func (c App) Render400(err_if models.ErrorInterface, err error) revel.Result {
	return c.RenderError(err_if, http.StatusBadRequest,
		ota_constant.ERROR_CODE_DROPPED, fmt.Sprintf("%s", err))
}

func (c App) Render500(err_if models.ErrorInterface, err error) revel.Result {
	return c.RenderError(err_if, http.StatusInternalServerError,
		ota_constant.ERROR_CODE_DROPPED, fmt.Sprintf("%s", err))
}

func (c App) Render404(err_if models.ErrorInterface, err error) revel.Result {
	return c.RenderError(err_if, http.StatusNotFound,
		ota_constant.ERROR_CODE_DROPPED, fmt.Sprintf("%s", err))
}

func (c App) Render404WithCode(err_if models.ErrorInterface, code int, err string) revel.Result {
	return c.RenderError(err_if, http.StatusNotFound,
		code, err)
}
