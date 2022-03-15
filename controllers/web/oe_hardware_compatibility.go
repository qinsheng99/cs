package controllers

import (
	"cve-sa-backend/handles/web"
	"cve-sa-backend/iniconf"
	cveSa "cve-sa-backend/utils/entity/cve_sa"
	"cve-sa-backend/utils/tools"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func GetOsListForHardware(c *gin.Context) {
	lang := c.DefaultQuery("lang", "zh")

	data, err := web.GetOsForHardware(lang)
	if err != nil {
		iniconf.SLog.Error("getOSList :", err)
		tools.Failure(c)
		return
	}

	tools.Success(c, data)
}

func GetArchitectureListForHardware(c *gin.Context) {
	lang := c.DefaultQuery("lang", "zh")

	data, err := web.GetArchitectureListForHardware(lang)
	if err != nil {
		iniconf.SLog.Error("getArchitectureList :", err)
		tools.Failure(c)
		return
	}

	tools.Success(c, data)
}

func FindAllHardwareCompatibility(c *gin.Context) {
	var req cveSa.OeCompSearchRequest

	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		iniconf.SLog.Error("findAllHardwareCompatibility :", err)
		tools.QueryFailure(c)
		return
	}

	datas, err := web.FindAllHardwareCompatibility(req)
	if err != nil {
		iniconf.SLog.Error(err)
		tools.Failure(c)
		return
	}

	tools.Success(c, datas)
}

func GetHardwareCompatibilityById(c *gin.Context) {
	var Id struct {
		Id int64 `form:"id"`
	}

	if err := c.ShouldBindWith(&Id, binding.Query); err != nil {
		iniconf.SLog.Error(err)
		tools.QueryFailure(c)
		return
	}

	data, err := web.GetHardwareCompatibilityById(Id.Id)
	if err != nil {
		iniconf.SLog.Error("getHardwareCompatibilityById :", err)
		tools.Failure(c)
		return
	}

	tools.Success(c, data)
}

func GetOeHardwareAdapterListByHardwareId(c *gin.Context) {
	var Id struct {
		Id int64 `form:"hardwareId"`
	}

	if err := c.ShouldBindWith(&Id, binding.Query); err != nil {
		iniconf.SLog.Error(err)
		tools.QueryFailure(c)
		return
	}

	datas, err := web.ByhardwareId(Id.Id)
	if err != nil {
		iniconf.SLog.Error("getOEHardwareAdapterListByHardwareId :", err)
		tools.Failure(c)
		return
	}

	tools.Success(c, datas)
}

func GetCpu(c *gin.Context) {
	lang := c.DefaultQuery("lang", "zh")

	datas, err := web.GetCpuList(lang)
	if err != nil {
		iniconf.SLog.Error(err)
		tools.Failure(c)
		return
	}

	tools.Success(c, datas)
}
