package controllers

import (
	cveSa "cve-sa-backend/utils/entity/cve_sa"
	"cve-sa-backend/utils/tools"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func FindAllSecurity(c *gin.Context) {
	var req cveSa.RequestData

	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		tools.QueryFailure(c)
		return
	}

	tools.Success(c, req)
}
