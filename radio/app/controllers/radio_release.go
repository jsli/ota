package controllers

import (
	"github.com/jsli/ota/radio/app/models"
	"github.com/robfig/revel"
	"net/http"
)

func (c Radio) ReleaseIndex() revel.Result {
	return c.RenderJson("Radio Release Index")
}

func (c Radio) Release(fp string) revel.Result {
	dal, err := models.NewDal()
	if err != nil {
		c.Response.Status = http.StatusInternalServerError
		return c.RenderJson(nil)
	}
	defer dal.Close()

	release, err := models.FindRadioOtaReleaseByFp(dal, fp)
	if err != nil {
		c.Response.Status = http.StatusInternalServerError
		return c.RenderJson(err)
	}
	if release != nil {
		return c.RenderJson(release)
	}

	return c.RenderJson(nil)
}
