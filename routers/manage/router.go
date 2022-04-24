package routers

import (
	"net/http"

	"cve-sa-backend/controllers"
	"cve-sa-backend/controllers/manage"

	"github.com/gin-gonic/gin"
)

func ManageRouters(r *gin.Engine) {

	r.LoadHTMLGlob("./webapp/*")
	manager := r.Group("/cve-security-notice-server")
	{
		uploadCon := controllers.Con.Manage.UploadController
		func(uploadCon manage.UploadController) {
			manager.GET("/manager", func(context *gin.Context) {
				context.HTML(http.StatusOK, "manager.html", nil)
			})
			manager.Use(managerAuth)
			manager.POST("/deleteCVE", uploadCon.DeleteCVE)
			manager.POST("/deleteSA", uploadCon.DeleteSA)
			manager.POST("/syncUnCVE", uploadCon.SyncCve)
			manager.POST("/getSyncCVE", uploadCon.GetHttpParserBeanListByCve)

			manager.POST("/syncHardware", uploadCon.SyncHardwareCompatibility)
			manager.POST("/syncDriver", uploadCon.SyncDriverCompatibility)
			manager.POST("/transfer", uploadCon.TransferOldData)
			manager.POST("/syncAll", uploadCon.SyncAll)
			manager.POST("/syncSA", uploadCon.SyncSA)
			manager.POST("/syncOsv", uploadCon.SyncOsv)
		}(uploadCon)
	}
}
