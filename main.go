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
	//myErr := iniconf.InitMysql()
	//if myErr != nil {
	//	fmt.Println("myErr: ", myErr)
	//	return
	//}

	gormErr := iniconf.InitGormMysql()
	if gormErr != nil {
		fmt.Println("gormErr: ", gormErr)
		return
	}
	//defer iniconf.Mysql.Close()
	//defer iniconf.LogFile.Close()
	r := gin.Default()
	webRouters.WebRouters(r)
	manageRouters.ManageRouters(r)

	runErr := r.Run(":8081")
	if runErr != nil {
		fmt.Println("runErr: ", runErr)
		return
	}
}
