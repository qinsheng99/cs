package iniconf

import (
	"fmt"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"io"
	"os"
	"time"
)

//var Logs *logrus.Logger
//var LogFile *os.File
//
//func InitLog() error {
//	logConf, err := Cfg.GetSection("log")
//	if err != nil {
//		fmt.Println("Fail to load section 'server': ", err)
//		return err
//	}
//	logLevel := logConf.Key("LOG_LEVEL").MustUint(5)
//	Logs = logrus.New()
//	lLog := logrus.TraceLevel
//	if logLevel == 5 {
//		lLog = logrus.DebugLevel
//	} else if logLevel == 4 {
//		lLog = logrus.InfoLevel
//	} else if logLevel == 3 {
//		lLog = logrus.WarnLevel
//	} else if logLevel == 2 {
//		lLog = logrus.ErrorLevel
//	} else if logLevel == 1 {
//		lLog = logrus.FatalLevel
//	} else if logLevel == 0 {
//		lLog = logrus.PanicLevel
//	}
//	Logs.SetLevel(lLog)
//	Logs.Formatter = &logrus.TextFormatter{
//		DisableColors:  true,
//		FullTimestamp:  true,
//		DisableSorting: true,
//		ForceColors:    true,
//		ForceQuote:     true,
//	}
//
//	loggerName := logConf.Key("LOG_NAME").MustString("sa-backend")
//	curTime := time.Now()
//	logFileName := fmt.Sprintf("%s_%04d-%02d-%02d.log",
//		loggerName, curTime.Year(), curTime.Month(), curTime.Day())
//	LogFile, err = os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeExclusive)
//	if err != nil {
//		fmt.Printf("try create logfile[%s] error[%s]\n", logFileName, err.Error())
//		return err
//	}
//	Logs.SetOutput(LogFile)
//	return nil
//}

var Log *zap.Logger
var SLog *zap.SugaredLogger

// InitLogger 初始化Logger
func InitLogger() (err error) {
	path, _ := os.Getwd()
	logConf, err := Cfg.GetSection("log")
	if err != nil {
		fmt.Println("Fail to load section 'server': ", err)
		return err
	}
	loggerName := logConf.Key("LOG_NAME").MustString("sa-backend")
	writeSyncer, err := getWriter(path + "/" + loggerName)
	if err != nil {
		return err
	}
	encoder := getEncoder()
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte("Info"))
	if err != nil {
		return
	}

	w := zapcore.NewMultiWriteSyncer(zapcore.AddSync(writeSyncer), zapcore.AddSync(os.Stdout))
	core := zapcore.NewCore(encoder, w, l)

	Log = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(Log)
	SLog = Log.Sugar()
	return
}

func getWriter(filename string) (io.Writer, error) {
	hook, err := rotatelogs.New(
		filename+".%Y-%m-%d",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour*24),
	)

	if err != nil {

		return nil, err
	}
	return hook, nil
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}
