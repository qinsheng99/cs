package web

import (
	"cve-sa-backend/handles"
	"cve-sa-backend/iniconf"
	cveSa "cve-sa-backend/utils/entity/cve_sa"
	"cve-sa-backend/utils/tools"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type EsController struct {
}

func (*EsController) Refresh(c *gin.Context) {
	err := handles.EsHandle.Refresh()
	if err != nil {
		tools.FailureErr(c, err)
	}
	tools.Success(c, "refresh success")
}

func (*EsController) Find(c *gin.Context) {
	var req cveSa.OeCompSearchRequest

	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		iniconf.SLog.Error(err)
		tools.QueryFailure(c)
		return
	}

	datas, err := handles.EsHandle.Find(req)
	if err != nil {
		tools.FailureErr(c, err)
		return
	}
	tools.Success(c, datas)
}
