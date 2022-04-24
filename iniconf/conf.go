package iniconf

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

var Cfg *ini.File

type AkAndSk struct {
	AK    string
	SK    string
	Point string
}

type es struct {
	Host string
	Port string
}

var Obs = new(AkAndSk)
var Es = new(es)

func InitConf() error {
	var err error
	Cfg, err = ini.Load("./conf/app.ini")
	if err != nil {
		fmt.Println("Fail to Load ‘conf/app.ini’:", err)
		return err
	}
	err = initObs()
	if err != nil {
		fmt.Println("Fail to Load obs:", err)
		return err
	}

	err = initEs()
	if err != nil {
		fmt.Println("Fail to Load obs:", err)
		return err
	}
	return nil
}

func initObs() error {
	logConf, err := Cfg.GetSection("obs")
	if err != nil {
		return err
	}
	ak := logConf.Key("Ak").MustString("AK")
	sk := logConf.Key("Sk").MustString("SK")
	endpoint := logConf.Key("END_POINT").String()

	Obs.SK = os.Getenv(sk)
	Obs.AK = os.Getenv(ak)
	Obs.Point = endpoint
	return nil
}

func initEs() error {
	logConf, err := Cfg.GetSection("es")
	if err != nil {
		return err
	}
	Es.Host = logConf.Key("ES_HOST").MustString("localhost")
	Es.Port = logConf.Key("ES_PORT").MustString("9200")
	return nil
}
