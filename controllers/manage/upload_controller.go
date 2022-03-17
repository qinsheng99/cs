package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"cve-sa-backend/handles/manage"
	"cve-sa-backend/iniconf"
	"cve-sa-backend/utils"
	"cve-sa-backend/utils/tools"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, "Uploading failed. Please select a file.")
		return
	}
	fmt.Println(file)

}

func DeleteCVE(c *gin.Context) {
	cveId := c.PostForm("deleteCVEID")
	packageName := c.PostForm("packageName")

	if cveId == "" {
		c.JSON(http.StatusOK, "Please enter CVE number.")
		return
	}

	result, err := manage.DeleteCVE(cveId, packageName)
	if err != nil {
		iniconf.Log.Error("delete cve error", zap.String("error", err.Error()))
		c.JSON(http.StatusOK, "deleteCVE failed. An exception occurred.")
		return
	}
	c.JSON(http.StatusOK, result)
	return
}

func SyncCve(c *gin.Context) {
	cveFileName := c.PostForm("cveNo")

	if cveFileName == "" {
		c.JSON(http.StatusOK, "CVE is null. Please input CVE.")
		return
	}

	if !strings.HasSuffix(cveFileName, ".xml") && !strings.HasSuffix(cveFileName, ".XML") {
		cveFileName += ".xml"
	}

	err := manage.SyncCve(cveFileName)
	fmt.Println(err)

}

func SyncHardwareCompatibility(c *gin.Context) {
	result, err := manage.SyncHardwareCompatibility()
	if err != nil {
		iniconf.SLog.Error("syncHardwareCompatibility failed :", err)
		c.JSON(http.StatusOK, "syncHardwareCompatibility failed. An exception occurred."+result+" "+err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func SyncDriverCompatibility(c *gin.Context) {
	result, err := manage.SyncDriverCompatibility()
	if err != nil {
		iniconf.SLog.Error("syncHardwareCompatibility failed :", err)
		c.JSON(http.StatusOK, "syncHardwareCompatibility failed. An exception occurred."+result+" "+err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func TransferOldData(c *gin.Context) {
	cve := c.PostForm("saNo")
	if cve == "" {
		c.JSON(http.StatusOK, "SA is null. Please input SA.")
		return
	}
	result, err := manage.TransferData(cve)
	if err != nil {
		iniconf.SLog.Error("transferOldData failed,", err)
		c.JSON(http.StatusOK, "transferOldData failed. An exception occurred.")
		return
	}
	c.JSON(http.StatusOK, result)
}

func SyncSA(c *gin.Context) {
	saFileName := c.PostForm("saNo")
	if saFileName == "" {
		c.JSON(http.StatusOK, "SA is null. Please input SA.")
		return
	}
	saFileName = strings.TrimSpace(saFileName)
	if !strings.HasSuffix(saFileName, ".xml") && !strings.HasSuffix(saFileName, ".XML") {
		saFileName += ".xml"
	}

	//result, err := manage.TransferData(saFileName)
	//if err != nil {
	//	iniconf.SLog.Error("SyncSA failed,", err)
	//	c.JSON(http.StatusOK, "SyncSA failed. An exception occurred."+result+err.Error())
	//	return
	//}
	//c.JSON(http.StatusOK, result)

	security, err := manage.ParserSA(saFileName)
	if err != nil {
		return
	}
	securityData := security.RCveSecurityNotice.CveSecurityNotice
	securityData.Updateime = utils.StrToTime(security.RCveSecurityNotice.Updateime)
	tools.Success(c, securityData)
}
