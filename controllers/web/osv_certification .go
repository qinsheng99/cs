package web

import (
	"net/http"
	"strconv"

	"cve-sa-backend/handles"
	"cve-sa-backend/iniconf"
	cveSa "cve-sa-backend/utils/entity/cve_sa"
	"cve-sa-backend/utils/tools"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Osv struct {
}

func (o *Osv) FindAllOsv(c *gin.Context) {
	var req cveSa.RequestOsv

	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		iniconf.SLog.Error(err)
		tools.QueryFailure(c)
		return
	}

	result, err := handles.OsvHandle.FindAllOsv(req)
	if err != nil {
		iniconf.SLog.Error("FindAllOsv :", err)
		tools.Failure(c)
		return
	}
	tools.Success(c, result)
}

func (o *Osv) GetOsvName(c *gin.Context) {
	datas, err := handles.OsvHandle.GetOsvName()
	if err != nil {
		iniconf.SLog.Error(err)
		tools.Failure(c)
		return
	}

	tools.Success(c, datas)
}

func (o *Osv) GetType(c *gin.Context) {
	datas, err := handles.OsvHandle.GetType()
	if err != nil {
		iniconf.SLog.Error(err)
		tools.Failure(c)
		return
	}

	tools.Success(c, datas)
}

func (o *Osv) GetOne(c *gin.Context) {
	Id, ok := c.GetQuery("id")
	if !ok {
		c.JSON(http.StatusBadRequest, "Required String parameter 'cveId' is not present")
		return
	}

	id, err := strconv.ParseInt(Id, 10, 64)
	if err != nil {
		tools.Failure(c)
		return
	}

	datas, err := handles.OsvHandle.GetOne(id)
	if err != nil {
		tools.Failure(c)
		return
	}

	tools.Success(c, datas)
	return
}
