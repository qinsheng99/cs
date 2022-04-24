package web

import (
	"cve-sa-backend/handles"
	"cve-sa-backend/iniconf"
	cveSa "cve-sa-backend/utils/entity/cve_sa"
	"cve-sa-backend/utils/tools"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type DriverCompatibility struct {
}

func (d *DriverCompatibility) GetOsListForDriver(c *gin.Context) {
	lang := c.DefaultQuery("lang", "zh")

	datas, err := handles.DriverHandle.GetOsList(lang)
	if err != nil {
		iniconf.SLog.Error("getOSList :", err)
		tools.Failure(c)
		return
	}
	tools.Success(c, datas)
}

func (d *DriverCompatibility) GetArchitectureListForDriver(c *gin.Context) {
	lang := c.DefaultQuery("lang", "zh")

	datas, err := handles.DriverHandle.GetArchitectureList(lang)
	if err != nil {
		iniconf.SLog.Error("getArchitectureList :", err)
		tools.Failure(c)
		return
	}
	tools.Success(c, datas)
}

func (d *DriverCompatibility) FindAllDriverCompatibility(c *gin.Context) {
	var req cveSa.OeCompSearchRequest

	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		iniconf.SLog.Error(err)
		tools.QueryFailure(c)
		return
	}

	datas, err := handles.DriverHandle.FindAllDriverCompatibility(req)
	if err != nil {
		iniconf.SLog.Error("findAllDriverCompatibility", err)
		tools.Failure(c)
		return
	}
	tools.Success(c, datas)
}
