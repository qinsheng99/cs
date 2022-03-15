package controllers

import (
	"cve-sa-backend/handles/web"
	"cve-sa-backend/iniconf"
	cveSa "cve-sa-backend/utils/entity/cve_sa"
	"cve-sa-backend/utils/tools"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func GetOsListForDriver(c *gin.Context) {
	lang := c.DefaultQuery("lang", "zh")

	datas, err := web.GetOsList(lang)
	if err != nil {
		iniconf.SLog.Error("getOSList :", err)
		tools.Failure(c)
		return
	}
	tools.Success(c, datas)
}

func GetArchitectureListForDriver(c *gin.Context) {
	lang := c.DefaultQuery("lang", "zh")

	datas, err := web.GetArchitectureList(lang)
	if err != nil {
		iniconf.SLog.Error("getArchitectureList :", err)
		tools.Failure(c)
		return
	}
	tools.Success(c, datas)
}

func FindAllDriverCompatibility(c *gin.Context) {
	var req cveSa.OeCompSearchRequest

	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		iniconf.SLog.Error(err)
		tools.QueryFailure(c)
		return
	}

	datas, err := web.FindAllDriverCompatibility(req)
	if err != nil {
		iniconf.SLog.Error("findAllDriverCompatibility", err)
		tools.Failure(c)
		return
	}
	tools.Success(c, datas)
}
