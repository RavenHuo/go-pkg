/**
 * @Author raven
 * @Description
 * @Date 2022/12/6
 **/
package log

import (
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

var log *logrus.Logger
var logOnce sync.Once

func getLogrus() *logrus.Logger {
	logOnce.Do(
		func() {
			log = buildLogrus()
		})
	return log
}

func buildLogrus() *logrus.Logger {

	//实例化
	lr := logrus.New()

	//设置输出
	lr.Out = os.Stdout
	lr.AddHook(newFileHook())
	//设置日志级别
	lr.SetLevel(logrus.DebugLevel)
	lr.SetReportCaller(true)
	//设置日志格式
	lr.SetFormatter(&logrus.TextFormatter{
		ForceQuote:      true, //键值对加引号
		TimestampFormat: "2006-01-02 15:04:05",
	})
	return lr
}
