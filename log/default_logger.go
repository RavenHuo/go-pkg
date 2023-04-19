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
	logRus := logrus.New()

	//设置输出
	logRus.Out = os.Stdout
	logrus.AddHook(newFileHook())
	//设置日志级别
	logRus.SetLevel(logrus.DebugLevel)

	//设置日志格式
	logRus.SetFormatter(&logrus.TextFormatter{
		ForceQuote:      true, //键值对加引号
		TimestampFormat: "2006-01-02 15:04:05",
	})
	return logRus
}
