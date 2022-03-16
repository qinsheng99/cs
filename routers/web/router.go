package routers

import (
	controllers "cve-sa-backend/controllers/web"

	"github.com/gin-gonic/gin"
)

func WebRouters(r *gin.Engine) {

	website := r.Group("/cve-security-notice-server")
	{
		cve := website.Group("/cvedatabase")
		{
			cve.GET("findAll", func(context *gin.Context) {

			})
		}
		securityNotice := website.Group("/securitynotice")
		{
			securityNotice.POST("/findAll", controllers.FindAllSecurity)
			securityNotice.POST("/getPackageLink", controllers.GetSecurityNoticePackageByPackageName)
			securityNotice.GET("/getByCveId", controllers.NoticeByCVEID)
			securityNotice.GET("/byCveIdAndAffectedComponent", controllers.ByCveIdAndAffectedComponent)
			securityNotice.GET("/getBySecurityNoticeNo", controllers.NoticeBySecurityNoticeNo)
		}

		oeDriverCompatibility := website.Group("/drivercomp")
		{
			oeDriverCompatibility.POST("/findAll", controllers.FindAllDriverCompatibility)
			oeDriverCompatibility.GET("/getOS", controllers.GetOsListForDriver)
			oeDriverCompatibility.GET("/getArchitecture", controllers.GetArchitectureListForDriver)
		}

		oeHardwareCompatibility := website.Group("/hardwarecomp")
		{
			oeHardwareCompatibility.POST("/findAll", controllers.FindAllHardwareCompatibility)
			oeHardwareCompatibility.GET("/getOS", controllers.GetOsListForHardware)
			oeHardwareCompatibility.GET("/getArchitecture", controllers.GetArchitectureListForHardware)
			oeHardwareCompatibility.GET("/getOne", controllers.GetHardwareCompatibilityById)
			oeHardwareCompatibility.GET("/getAdapterList", controllers.GetOeHardwareAdapterListByHardwareId)
			oeHardwareCompatibility.GET("/getCpu", controllers.GetCpu)
		}
	}
}
