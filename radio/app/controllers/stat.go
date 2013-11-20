package controllers

import (
	"github.com/jsli/ota/radio/app/models"
	"github.com/jsli/ota/radio/app/policy"
	"github.com/robfig/revel"
	"strconv"
)

type Stat struct {
	*revel.Controller
}

func (c Stat) Index() revel.Result {
	return c.Redirect("/stat/all")
}

func (c Stat) Self() revel.Result {
	rdal, err := models.NewRedisDal()
	if err != nil {
		return c.RenderJson(err)
	}
	ip := policy.FilterIp(c.Request.RemoteAddr)
	counter, err := rdal.Get(ip)
	if err != nil {
		return c.RenderJson(err)
	}

	result := models.StatSelfResult{}
	result.IP = ip
	result.Counter, _ = strconv.Atoi(counter)

	return c.RenderJson(result)
}

func (c Stat) All() revel.Result {
	rdal, err := models.NewRedisDal()
	if err != nil {
		return c.RenderJson(err)
	}

	result := models.StatAllResult{}
	result.Detail = make(map[string]int)
	keys, err := rdal.Keys("*")
	if err != nil {
		return c.RenderJson(err)
	}
	for _, key := range keys {
		count, _ := rdal.Get(key)
		count_i, _ := strconv.Atoi(count)
		result.Detail[key] = count_i
		result.TotalCount += count_i
	}
	return c.RenderJson(result)
}
