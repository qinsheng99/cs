package controllers

import (
	"cve-sa-backend/handles/web"
	"cve-sa-backend/iniconf"
	cveSa "cve-sa-backend/utils/entity/cve_sa"
	"cve-sa-backend/utils/tools"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func FindAllCVEDatabase(c *gin.Context) {
	var req cveSa.RequestData

	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		iniconf.SLog.Error(err)
		tools.QueryFailure(c)
		return
	}
	datas, err := web.FindAllCVEDatabase(req)
	if err != nil {
		iniconf.SLog.Error("findAllSecurityNotice :", err)
		tools.Failure(c)
		return
	}

	tools.Success(c, datas)
}
