package controllers

import (
	"cve-sa-backend/handles/web"
	"cve-sa-backend/iniconf"
	cveSa "cve-sa-backend/utils/entity/cve_sa"
	"cve-sa-backend/utils/tools"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func FindAllSecurity(c *gin.Context) {
	var req cveSa.RequestData

	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		iniconf.SLog.Error(err)
		tools.QueryFailure(c)
		return
	}
	datas, err := web.FindAllSecurity(req)
	if err != nil {
		iniconf.SLog.Error("findAllSecurityNotice :", err)
		tools.Failure(c)
		return
	}

	tools.Success(c, datas)
}
func GetSecurityNoticePackageByPackageName(c *gin.Context) {
	var req cveSa.RequestData

	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		iniconf.SLog.Error(err)
		tools.QueryFailure(c)
		return
	}
	datas, err := web.GetSecurityNoticePackageByPackageName(req.PackageName)
	if err != nil {
		iniconf.SLog.Error("getSecurityNoticePackageByPackageName :", err)
		tools.Failure(c)
		return
	}
	tools.Success(c, datas)
}

func NoticeByCVEID(c *gin.Context) {
	cveId := c.DefaultQuery("cveId", "")
	if cveId == "" {
		tools.Success(c, nil)
		return
	}
	datas, err := web.NoticeByCVEID(cveId)
	if err != nil {
		iniconf.SLog.Error("getCVEDatabaseByCVEID :", err)
		tools.Failure(c)
		return
	}
	tools.Success(c, datas)
}

func ByCveIdAndAffectedComponent(c *gin.Context) {
	cveId := c.DefaultQuery("cveId", "")
	affectedComponent := c.DefaultQuery("affectedComponent", "")
	datas, err := web.ByCveIdAndAffectedComponent(cveId, affectedComponent)
	if err != nil {
		iniconf.SLog.Error("byCveIdAndAffectedComponent :", err)
		tools.Failure(c)
		return
	}
	tools.Success(c, datas)
}

func NoticeBySecurityNoticeNo(c *gin.Context) {
	securityNoticeNo := c.DefaultQuery("securityNoticeNo", "")
	if securityNoticeNo == "" {
		tools.Success(c, nil)
		return
	}
	data, err := web.NoticeBySecurityNoticeNo(securityNoticeNo)
	if err != nil {
		iniconf.SLog.Error("getSecurityNoticeBySecurityNoticeNo :", err)
		tools.Failure(c)
		return
	}
	tools.Success(c, data)
}