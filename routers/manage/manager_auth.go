package routers

import (
	"net/http"

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
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username != uploadUser.username || password != uploadUser.password {
		c.Abort()
		c.String(http.StatusOK, "Username or password error")
		return
	}
	c.Next()
}
