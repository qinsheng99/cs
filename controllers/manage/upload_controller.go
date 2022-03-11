package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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

}
