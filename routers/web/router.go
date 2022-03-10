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
		}
	}
}
