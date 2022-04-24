package manage

import (
	"errors"
	"net/http"
	"strings"

	"cve-sa-backend/handles"
	"cve-sa-backend/iniconf"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UploadController struct {
}

func (u *UploadController) DeleteCVE(c *gin.Context) {
	cveId := c.PostForm("deleteCVEID")
	packageName := c.PostForm("packageName")

	if cveId == "" {
		c.JSON(http.StatusOK, "Please enter CVE number.")
		return
	}
	cveId = strings.TrimSpace(cveId)
	packageName = strings.TrimSpace(packageName)

	result, err := handles.UploadHandle.DeleteCVE(cveId, packageName)
	if err != nil {
		iniconf.Log.Error("delete cve error", zap.String("error", err.Error()))
		c.JSON(http.StatusOK, "deleteCVE failed. An exception occurred.")
		return
	}
	c.JSON(http.StatusOK, result)
	return
}

func (u *UploadController) DeleteSA(c *gin.Context) {
	saNo := c.PostForm("deleteSAID")
	if saNo == "" {
		c.JSON(http.StatusOK, "Please enter CVE number.")
		return
	}

	saNo = strings.TrimSpace(saNo)
	err := handles.UploadHandle.DeleteSA(saNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			iniconf.Log.Warn("SA NO error, manage.DeleteSA errors is:", zap.Error(err))
			c.JSON(http.StatusOK, "SA NO error")
			return
		} else {
			iniconf.Log.Error("manage.DeleteSA errors", zap.Error(err))
			c.JSON(http.StatusOK, "deleteCVE failed. An exception occurred.")
			return
		}
	}
	c.JSON(http.StatusOK, "deleteCVE success.")
	return
}

func (u *UploadController) GetHttpParserBeanListByCve(c *gin.Context) {
	cve := c.PostForm("cveNo")
	packageName := c.PostForm("packageName")
	if cve == "" {
		c.JSON(http.StatusOK, "CVE is null. Please input CVE.")
		return
	}

	cve = strings.TrimSpace(cve)
	packageName = strings.TrimSpace(packageName)
	result, err := handles.UploadHandle.GetHttpParserBeanListByCve(cve, packageName)
	if err != nil {
		iniconf.Log.Error("manage.GetHttpParserBeanListByCve error", zap.Error(err), zap.String("cve", cve), zap.String("packageName", packageName))
		c.JSON(http.StatusOK, "GetHttpParserBeanListByCve failed. An exception occurred.")
		return
	}

	c.JSON(http.StatusOK, result)
	return
}

func (u *UploadController) SyncCve(c *gin.Context) {
	cveFileName := c.PostForm("cveNo")

	if cveFileName == "" {
		c.JSON(http.StatusOK, "CVE is null. Please input CVE.")
		return
	}
	cveFileName = strings.TrimSpace(cveFileName)

	if !strings.HasSuffix(cveFileName, ".xml") && !strings.HasSuffix(cveFileName, ".XML") {
		cveFileName += ".xml"
	}

	result, err := handles.UploadHandle.SyncCve(cveFileName)
	if err != nil {
		iniconf.Log.Error("manage.SyncCve error:", zap.Error(err))
		c.JSON(http.StatusOK, "SyncCve failed. An exception occurred.")
		return
	}
	c.JSON(http.StatusOK, result)
	return
}

func (u *UploadController) SyncHardwareCompatibility(c *gin.Context) {
	result, err := handles.UploadHandle.SyncHardwareCompatibility()
	if err != nil {
		iniconf.SLog.Error("syncHardwareCompatibility failed :", err)
		c.JSON(http.StatusOK, "syncHardwareCompatibility failed. An exception occurred."+result+err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (u *UploadController) SyncDriverCompatibility(c *gin.Context) {
	result, err := handles.UploadHandle.SyncDriverCompatibility()
	if err != nil {
		iniconf.SLog.Error("syncHardwareCompatibility failed :", err)
		c.JSON(http.StatusOK, "syncHardwareCompatibility failed. An exception occurred."+result+err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (u *UploadController) TransferOldData(c *gin.Context) {
	cve := c.PostForm("saNo")
	if cve == "" {
		c.JSON(http.StatusOK, "SA is null. Please input SA.")
		return
	}
	cve = strings.TrimSpace(cve)
	result, err := handles.UploadHandle.TransferData(cve)
	if err != nil {
		iniconf.SLog.Error("transferOldData failed,", err)
		c.JSON(http.StatusOK, "transferOldData failed. An exception occurred.")
		return
	}
	c.JSON(http.StatusOK, result)
}

func (u *UploadController) SyncSA(c *gin.Context) {
	saFileName := c.PostForm("saNo")
	if saFileName == "" {
		c.JSON(http.StatusOK, "SA is null. Please input SA.")
		return
	}
	saFileName = strings.TrimSpace(saFileName)
	if !strings.HasSuffix(saFileName, ".xml") && !strings.HasSuffix(saFileName, ".XML") {
		saFileName += ".xml"
	}

	result, err := handles.UploadHandle.SyncSA(saFileName)
	if err != nil {
		iniconf.SLog.Error("SyncSA failed,", err)
		c.JSON(http.StatusOK, "SyncSA failed. An exception occurred."+result)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (u *UploadController) SyncAll(c *gin.Context) {
	result, err := handles.UploadHandle.SyncAll()

	if err != nil {
		iniconf.SLog.Error("SyncAll failed,", err)
		c.JSON(http.StatusOK, result+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (u *UploadController) SyncOsv(c *gin.Context) {
	result, err := handles.UploadHandle.SyncOsv()
	if err != nil {
		iniconf.SLog.Error("syncOsv failed :", err)
		c.JSON(http.StatusOK, "syncOsv failed. An exception occurred."+result+err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}
