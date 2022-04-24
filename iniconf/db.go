package iniconf

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySqlConf struct {
	DbHost    string
	DbPort    int64
	DbUser    string
	DbPwd     string
	DbName    string
	DbMaxConn int
	DbMaxidle int
	DbType    string
}

//var Mysql *sql.DB

const CONNMAXLIFTIME = 900

func NewMysql(myc *MySqlConf) {
	logConf, err := Cfg.GetSection("db")
	if err != nil {
		SLog.Error("Fail to load section 'server': ", err)
		return
	}
	dbHost := logConf.Key("DB_HOST").MustString("DB_URL")
	myc.DbHost = os.Getenv(dbHost)
	//myc.DbHost = dbHost
	dbPort := logConf.Key("DB_PORT").MustInt64(3306)
	//myc.DbPort = os.Getenv(dbPort)
	myc.DbPort = dbPort
	dbUser := logConf.Key("DB_USER").MustString("DB_USER")
	myc.DbUser = os.Getenv(dbUser)
	//myc.DbUser = dbUser
	dbPwd := logConf.Key("DB_PWD").MustString("DB_PWD")
	myc.DbPwd = os.Getenv(dbPwd)
	//myc.DbPwd = dbPwd
	dbName := logConf.Key("DB_NAME").MustString("DB_NAME")
	//myc.DbName = os.Getenv(dbName)
	myc.DbName = dbName
	dbMaxidle := logConf.Key("DB_MAXIDLE").MustInt(10)
	//myc.DbMaxidle = os.Getenv(dbMaxidle)
	myc.DbMaxidle = dbMaxidle
	dbMaxconn := logConf.Key("DB_MAXCONN").MustInt(100)
	//myc.DbMaxConn = os.Getenv(dbMaxconn)
	myc.DbMaxConn = dbMaxconn
	dbType := logConf.Key("DB_TYPE").MustString("DB_TYPE")
	//myc.DbName = os.Getenv(dbName)
	myc.DbType = dbType
}

//func InitMysql() error {
//	myc := MySqlConf{}
//	NewMysql(&myc)
//	connErr := errors.New("")
//	mysqlConf := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", myc.DbUser, myc.DbPwd, myc.DbHost, myc.DbPort, myc.DbName)
//	Mysql, connErr = sql.Open(myc.DbType, mysqlConf)
//	if connErr != nil {
//		Logs.Error("InitMysql, connErr: ", connErr)
//		return connErr
//	}
//	Mysql.SetConnMaxLifetime(CONNMAXLIFTIME)
//	Mysql.SetMaxOpenConns(myc.DbMaxConn)
//	Mysql.SetMaxIdleConns(myc.DbMaxidle)
//	//验证连接
//	if err := Mysql.Ping(); err != nil {
//		Logs.Error("open database fail")
//		return err
//	}
//	Logs.Info("connnect success")
//	return nil
//}

var DB *gorm.DB

func InitGormMysql() error {
	myc := MySqlConf{}
	NewMysql(&myc)
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local", myc.DbUser, myc.DbPwd, myc.DbHost, myc.DbPort, myc.DbName)
	db, err := gorm.Open(gormmysql.New(gormmysql.Config{
		DSN:                       dsn,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{})
	if err != nil {
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetConnMaxLifetime(CONNMAXLIFTIME)
	sqlDB.SetMaxOpenConns(myc.DbMaxConn)
	sqlDB.SetMaxIdleConns(myc.DbMaxidle)

	DB = db

	return nil
}

func GetDb() *gorm.DB {
	return DB
}
