package log

import (
	"gosunxifel/config"

	"github.com/donnie4w/go-logger/logger"
)

/*
const (
	ALL LEVEL = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
	OFF
)
*/
const (
	Log_console    bool         = true
	Log_Dir        string       = "logs"
	Log_Filename   string       = "app.log"
	Log_MaxBackup  int32        = 1000
	Log_SizeBackup int64        = 10240
	Log_Level      logger.LEVEL = logger.DEBUG
)

func init() {
	//指定是否控制台打印，默认为true
	logger.SetConsole(Log_console)
	//指定日志文件备份方式为文件大小的方式
	//第一个参数为日志文件存放目录
	//第二个参数为日志文件命名
	//第三个参数为备份文件最大数量
	//第四个参数为备份文件大小
	//第五个参数为文件大小的单位
	logger.SetRollingFile(config.GetRootdir()+"/logs", Log_Filename, Log_MaxBackup, Log_SizeBackup, logger.KB)
	//	logger.SetRollingFile()
	//指定日志文件备份方式为日期的方式
	//第一个参数为日志文件存放目录
	//第二个参数为日志文件命名
	//	logger.SetRollingDaily("./", "test.log")

	//指定日志级别  ALL，DEBUG，INFO，WARN，ERROR，FATAL，OFF 级别由低到高
	//一般习惯是测试阶段为debug，生成环境为info以上
	level := config.Cfg.Section("log").Key("level").MustInt()

	logger.SetLevel(logger.LEVEL(level))

	//		logger.Debug("Debug>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
	//	logger.Info("Info>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
	//	logger.Warn("Warn>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
	//	logger.Error("Error>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
	//	logger.Fatal("Fatal>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
}

func Debug(v ...interface{}) {
	logger.Debug(v...)
}
func Info(v ...interface{}) {
	logger.Info(v...)
}
func Warn(v ...interface{}) {
	logger.Warn(v...)
}
func Error(v ...interface{}) {
	logger.Error(v...)
}
func Fatal(v ...interface{}) {
	logger.Fatal(v...)
}
