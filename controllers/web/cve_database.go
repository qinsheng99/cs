package web

import (
	"net/http"

	"cve-sa-backend/handles"
	"cve-sa-backend/iniconf"
	cveSa "cve-sa-backend/utils/entity/cve_sa"
	"cve-sa-backend/utils/tools"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type CveDatabaseCon struct {
}

func (cv *CveDatabaseCon) FindAllCVEDatabase(c *gin.Context) {
	var req cveSa.RequestData

	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		iniconf.SLog.Error(err)
		tools.QueryFailure(c)
		return
	}
	datas, err := handles.CveDatabaseHandle.FindAllCVEDatabase(req)
	if err != nil {
		iniconf.SLog.Error("FindAllCVEDatabase :", err)
		tools.Failure(c)
		return
	}

	tools.Success(c, datas)
}

func (cv *CveDatabaseCon) GetByCveIdAndPackageName(c *gin.Context) {
	cveId, ok := c.GetQuery("cveId")
	if !ok {
		c.JSON(http.StatusBadRequest, "Required String parameter 'cveId' is not present")
		return
	}
	packageName, ok := c.GetQuery("packageName")
	if !ok {
		c.JSON(http.StatusBadRequest, "Required String parameter 'cveId' is not present")
		return
	}

	result, err := handles.CveDatabaseHandle.GetByCveIdAndPackageName(cveId, packageName)
	if err != nil {
		tools.Failure(c)
		return
	}

	tools.Success(c, result)
	return
}

func (cv *CveDatabaseCon) GetCVEProductPackageListByCveId(c *gin.Context) {
	cveId, ok := c.GetQuery("cveId")
	if !ok {
		c.JSON(http.StatusBadRequest, "Required String parameter 'cveId' is not present")
		return
	}

	result, err := handles.CveDatabaseHandle.GetCVEProductPackageListByCveId(cveId)
	if err != nil {
		tools.Failure(c)
		return
	}

	tools.Success(c, result)
	return
}

func (cv *CveDatabaseCon) GetCVEProductPackageList(c *gin.Context) {
	cveId, ok := c.GetQuery("cveId")
	if !ok {
		c.JSON(http.StatusBadRequest, "Required String parameter 'cveId' is not present")
		return
	}
	packageName, ok := c.GetQuery("packageName")
	if !ok {
		c.JSON(http.StatusBadRequest, "Required String parameter 'cveId' is not present")
		return
	}

	result, err := handles.CveDatabaseHandle.GetCVEProductPackageList(cveId, packageName)
	if err != nil {
		tools.Failure(c)
		return
	}

	tools.Success(c, result)
	return
}
