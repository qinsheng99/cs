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

type SecurityNoticeCon struct {
}

func (s *SecurityNoticeCon) FindAllSecurity(c *gin.Context) {
	var req cveSa.RequestData

	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		iniconf.SLog.Error("request err,", err)
		tools.QueryFailure(c)
		return
	}
	datas, err := handles.SecurityHandle.FindAllSecurity(req)
	if err != nil {
		iniconf.SLog.Error("findAllSecurityNotice :", err)
		tools.Failure(c)
		return
	}

	tools.Success(c, datas)
}
func (s *SecurityNoticeCon) GetSecurityNoticePackageByPackageName(c *gin.Context) {
	var req cveSa.RequestData

	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		iniconf.SLog.Error("request err,", err)
		tools.QueryFailure(c)
		return
	}
	datas, err := handles.SecurityHandle.GetSecurityNoticePackageByPackageName(req.PackageName)
	if err != nil {
		iniconf.SLog.Error("getSecurityNoticePackageByPackageName :", err)
		tools.Failure(c)
		return
	}
	tools.Success(c, datas)
}

func (s *SecurityNoticeCon) NoticeByCVEID(c *gin.Context) {
	cveId, ok := c.GetQuery("cveId")
	if !ok {
		c.JSON(http.StatusBadRequest, "Required String parameter 'cveId' is not present")
		return
	}
	if cveId == "" {
		tools.Success(c, nil)
		return
	}
	datas, err := handles.SecurityHandle.NoticeByCVEID(cveId)
	if err != nil {
		iniconf.SLog.Error("getCVEDatabaseByCVEID :", err)
		tools.Failure(c)
		return
	}
	tools.Success(c, datas)
}

func (s *SecurityNoticeCon) ByCveIdAndAffectedComponent(c *gin.Context) {
	cveId, ok := c.GetQuery("cveId")
	if !ok {
		c.JSON(http.StatusBadRequest, "Required String parameter 'cveId' is not present")
		return
	}
	if cveId == "" {
		tools.Success(c, nil)
		return
	}
	affectedComponent, ok := c.GetQuery("affectedComponent")
	if !ok {
		c.JSON(http.StatusBadRequest, "Required String parameter 'affectedComponent' is not present")
		return
	}
	datas, err := handles.SecurityHandle.ByCveIdAndAffectedComponent(cveId, affectedComponent)
	if err != nil {
		iniconf.SLog.Error("byCveIdAndAffectedComponent :", err)
		tools.Failure(c)
		return
	}
	tools.Success(c, datas)
}

func (s *SecurityNoticeCon) NoticeBySecurityNoticeNo(c *gin.Context) {
	securityNoticeNo, ok := c.GetQuery("securityNoticeNo")
	if !ok {
		c.JSON(http.StatusBadRequest, "Required String parameter 'securityNoticeNo' is not present")
		return
	}
	if securityNoticeNo == "" {
		tools.Success(c, nil)
		return
	}
	data, err := handles.SecurityHandle.NoticeBySecurityNoticeNo(securityNoticeNo)
	if err != nil {
		iniconf.SLog.Error("getSecurityNoticeBySecurityNoticeNo :", err)
		tools.Failure(c)
		return
	}
	tools.Success(c, data)
}
