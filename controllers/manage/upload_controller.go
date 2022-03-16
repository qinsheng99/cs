package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"cve-sa-backend/handles/manage"
	"cve-sa-backend/iniconf"
	"cve-sa-backend/utils/tools"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusOK, "Uploading failed. Please select a file.")
		return
	}
	fmt.Println(file)

}

func DeleteCVE(c *gin.Context) {
	cveId := c.PostForm("deleteCVEID")
	packageName := c.PostForm("packageName")

	if cveId == "" {
		c.String(http.StatusOK, "Please enter CVE number.")
		return
	}

	result, err := manage.DeleteCVE(cveId, packageName)
	if err != nil {
		iniconf.Log.Error("delete cve error", zap.String("error", err.Error()))
		tools.Failure(c)
		return
	}
	tools.Success(c, result)
	return
}

func SyncCve(c *gin.Context) {
	cveFileName := c.PostForm("cveNo")

	if cveFileName == "" {
		c.String(http.StatusOK, "CVE is null. Please input CVE.")
		return
	}

	if !strings.HasSuffix(cveFileName, ".xml") && !strings.HasSuffix(cveFileName, ".XML") {
		cveFileName += ".xml"
	}

}