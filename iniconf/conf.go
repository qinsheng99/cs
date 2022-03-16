package iniconf

import (
	"fmt"
	"gopkg.in/ini.v1"
)

var Cfg *ini.File

func InitConf() error {
	var err error
	Cfg, err = ini.Load("./conf/app.ini")
	if err != nil {
		fmt.Println("Fail to Load ‘conf/app.ini’:", err)
		return err
	}
	return nil
}