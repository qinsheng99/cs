package routers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var uploadUser struct {
	username string
	password string
}

func init() {
	uploadUser.username = "admin"
	uploadUser.password = "admin@1234"
}

func managerAuth(c *gin.Context) {
	username := strings.TrimSpace(c.PostForm("username"))
	password := strings.TrimSpace(c.PostForm("password"))
	if username != uploadUser.username || password != uploadUser.password {
		c.Abort()
		c.JSON(http.StatusOK, "Username or password error")
		return
	}
	c.Next()
}
