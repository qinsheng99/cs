package routers

import (
	"cve-sa-backend/controllers"
	"cve-sa-backend/controllers/web"

	"github.com/gin-gonic/gin"
)

func WebRouters(r *gin.Engine) {

	website := r.Group("/cve-security-notice-server")
	{
		cve := website.Group("/cvedatabase")
		cveCon := controllers.Con.Web.CveDatabaseCon
		func(cveCon web.CveDatabaseCon) {
			cve.POST("/findAll", cveCon.FindAllCVEDatabase)
			cve.GET("/getByCveIdAndPackageName", cveCon.GetByCveIdAndPackageName)
			cve.GET("/getPackageByCveId", cveCon.GetCVEProductPackageListByCveId)
			cve.GET("/getCVEProductPackageList", cveCon.GetCVEProductPackageList)
		}(cveCon)

		securityNotice := website.Group("/securitynotice")
		securityNoticeCon := controllers.Con.Web.SecurityNoticeCon
		func(securityNoticeCon web.SecurityNoticeCon) {
			securityNotice.POST("/findAll", securityNoticeCon.FindAllSecurity)
			securityNotice.POST("/getPackageLink", securityNoticeCon.GetSecurityNoticePackageByPackageName)
			securityNotice.GET("/getByCveId", securityNoticeCon.NoticeByCVEID)
			securityNotice.GET("/byCveIdAndAffectedComponent", securityNoticeCon.ByCveIdAndAffectedComponent)
			securityNotice.GET("/getBySecurityNoticeNo", securityNoticeCon.NoticeBySecurityNoticeNo)
		}(securityNoticeCon)

		oeDriverCompatibility := website.Group("/drivercomp")
		oeDriverCompatibilityCon := controllers.Con.Web.DriverCompatibility
		func(oeDriverCompatibilityCon web.DriverCompatibility) {
			oeDriverCompatibility.POST("/findAll", oeDriverCompatibilityCon.FindAllDriverCompatibility)
			oeDriverCompatibility.GET("/getOS", oeDriverCompatibilityCon.GetOsListForDriver)
			oeDriverCompatibility.GET("/getArchitecture", oeDriverCompatibilityCon.GetArchitectureListForDriver)
		}(oeDriverCompatibilityCon)

		oeHardwareCompatibility := website.Group("/hardwarecomp")
		oeHardwareCompatibilityCon := controllers.Con.Web.HardwareCompatibility
		func(oeHardwareCompatibilityCon web.HardwareCompatibility) {
			oeHardwareCompatibility.POST("/findAll", oeHardwareCompatibilityCon.FindAllHardwareCompatibility)
			oeHardwareCompatibility.GET("/getOS", oeHardwareCompatibilityCon.GetOsListForHardware)
			oeHardwareCompatibility.GET("/getArchitecture", oeHardwareCompatibilityCon.GetArchitectureListForHardware)
			oeHardwareCompatibility.GET("/getOne", oeHardwareCompatibilityCon.GetHardwareCompatibilityById)
			oeHardwareCompatibility.GET("/getAdapterList", oeHardwareCompatibilityCon.GetOeHardwareAdapterListByHardwareId)
			oeHardwareCompatibility.GET("/getCpu", oeHardwareCompatibilityCon.GetCpu)
		}(oeHardwareCompatibilityCon)

		osvCompatibility := website.Group("/osv")
		osvCompatibilityCon := controllers.Con.Web.Osv
		func(osvCompatibilityCon web.Osv) {
			osvCompatibility.POST("/findAll", osvCompatibilityCon.FindAllOsv)
			osvCompatibility.GET("/getOsName", osvCompatibilityCon.GetOsvName)
			osvCompatibility.GET("/getType", osvCompatibilityCon.GetType)
			osvCompatibility.GET("/getOne", osvCompatibilityCon.GetOne)
		}(osvCompatibilityCon)
	}

	es := r.Group("/es")
	esCon := controllers.Con.Web.EsController
	func(osvCompatibilityCon web.EsController) {
		es.GET("/refresh", esCon.Refresh)
		es.POST("/find", esCon.Find)
	}(esCon)
}
