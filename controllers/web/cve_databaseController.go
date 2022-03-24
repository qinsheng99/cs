package controllers

import (
	"net/http"

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

func GetByCveIdAndPackageName(c *gin.Context) {
	cveId := c.DefaultQuery("cveId", "")
	packageName := c.DefaultQuery("packageName", "")

	if cveId == "" {
		c.JSON(http.StatusBadRequest, "Required String parameter 'cveId' is not present")
		return
	}

	if packageName == "" {
		c.JSON(http.StatusBadRequest, "Required String parameter 'packagName' is not present")
		return
	}

	result, err := web.GetByCveIdAndPackageName(cveId, packageName)
	if err != nil {
		tools.Failure(c)
		return
	}

	tools.Success(c, result)
	return
}

func GetCVEProductPackageListByCveId(c *gin.Context) {
	cveId := c.DefaultQuery("cveId", "")

	if cveId == "" {
		c.JSON(http.StatusBadRequest, "Required String parameter 'cveId' is not present")
		return
	}

	result, err := web.GetCVEProductPackageListByCveId(cveId)
	if err != nil {
		tools.Failure(c)
		return
	}

	tools.Success(c, result)
	return
}

func GetCVEProductPackageList(c *gin.Context) {
	cveId := c.DefaultQuery("cveId", "")
	packageName := c.DefaultQuery("packageName", "")

	if cveId == "" {
		c.JSON(http.StatusBadRequest, "Required String parameter 'cveId' is not present")
		return
	}

	if packageName == "" {
		c.JSON(http.StatusBadRequest, "Required String parameter 'packagName' is not present")
		return
	}

	result, err := web.GetCVEProductPackageList(cveId, packageName)
	if err != nil {
		tools.Failure(c)
		return
	}

	tools.Success(c, result)
	return
}
