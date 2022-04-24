package web

import (
	"cve-sa-backend/handles"
	"cve-sa-backend/iniconf"
	cveSa "cve-sa-backend/utils/entity/cve_sa"
	"cve-sa-backend/utils/tools"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type HardwareCompatibility struct {
}

func (h *HardwareCompatibility) GetOsListForHardware(c *gin.Context) {
	lang := c.DefaultQuery("lang", "zh")

	data, err := handles.HardwareHandle.GetOsForHardware(lang)
	if err != nil {
		iniconf.SLog.Error("getOSList :", err)
		tools.Failure(c)
		return
	}

	tools.Success(c, data)
}

func (h *HardwareCompatibility) GetArchitectureListForHardware(c *gin.Context) {
	lang := c.DefaultQuery("lang", "zh")

	data, err := handles.HardwareHandle.GetArchitectureListForHardware(lang)
	if err != nil {
		iniconf.SLog.Error("getArchitectureList :", err)
		tools.Failure(c)
		return
	}

	tools.Success(c, data)
}

func (h *HardwareCompatibility) FindAllHardwareCompatibility(c *gin.Context) {
	var req cveSa.OeCompSearchRequest

	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		iniconf.SLog.Error("findAllHardwareCompatibility :", err)
		tools.QueryFailure(c)
		return
	}

	datas, err := handles.HardwareHandle.FindAllHardwareCompatibility(req)
	if err != nil {
		iniconf.SLog.Error(err)
		tools.Failure(c)
		return
	}

	tools.Success(c, datas)
}

func (h *HardwareCompatibility) GetHardwareCompatibilityById(c *gin.Context) {
	var Id struct {
		Id int64 `form:"id"`
	}

	if err := c.ShouldBindWith(&Id, binding.Query); err != nil {
		iniconf.SLog.Error(err)
		tools.QueryFailure(c)
		return
	}

	data, err := handles.HardwareHandle.GetHardwareCompatibilityById(Id.Id)
	if err != nil {
		iniconf.SLog.Error("getHardwareCompatibilityById :", err)
		tools.Failure(c)
		return
	}

	tools.Success(c, data)
}

func (h *HardwareCompatibility) GetOeHardwareAdapterListByHardwareId(c *gin.Context) {
	var Id struct {
		Id int64 `form:"hardwareId"`
	}

	if err := c.ShouldBindWith(&Id, binding.Query); err != nil {
		iniconf.SLog.Error(err)
		tools.QueryFailure(c)
		return
	}

	datas, err := handles.HardwareHandle.ByhardwareId(Id.Id)
	if err != nil {
		iniconf.SLog.Error("getOEHardwareAdapterListByHardwareId :", err)
		tools.Failure(c)
		return
	}

	tools.Success(c, datas)
}

func (h *HardwareCompatibility) GetCpu(c *gin.Context) {
	lang := c.DefaultQuery("lang", "zh")

	datas, err := handles.HardwareHandle.GetCpuList(lang)
	if err != nil {
		iniconf.SLog.Error(err)
		tools.Failure(c)
		return
	}

	tools.Success(c, datas)
}
