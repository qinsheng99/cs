package main

import (
	"fmt"

	"cve-sa-backend/iniconf"
	manageRouters "cve-sa-backend/routers/manage"
	webRouters "cve-sa-backend/routers/web"

	"github.com/gin-gonic/gin"
)

func main() {
	// init config
	inErr := iniconf.InitConf()
	if inErr != nil {
		fmt.Println("inErr: ", inErr)
		return
	}
	// init logs
	logErr := iniconf.InitLogger()
	if logErr != nil {
		fmt.Println("logErr: ", logErr)
		return
	}
	// init mysql
	gormErr := iniconf.InitGormMysql()
	if gormErr != nil {
		fmt.Println("gormErr: ", gormErr)
		return
	}

	// init mysql
	esErr := iniconf.InitEs()
	if esErr != nil {
		fmt.Println("esErr: ", esErr)
		return
	}

	r := gin.Default()

	webRouters.WebRouters(r)
	manageRouters.ManageRouters(r)

	Conf, err := iniconf.Cfg.GetSection("basic")
	if err != nil {
		fmt.Println("Fail to load section 'server': ", err)
		return
	}

	port := Conf.Key("PORT").MustString("8080")
	runErr := r.Run(":" + port)
	if runErr != nil {
		iniconf.SLog.Error("server run failed,", runErr)
		fmt.Println("runErr: ", runErr)
		return
	}
}
