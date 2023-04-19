/**
 * @Author raven
 * @Description
 * @Date 2023/4/19
 **/
package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

// FileHook 文件输出钩子
type FileHook struct {
	logger *logrus.Logger
}

func newFileHook() *FileHook {
	now := time.Now()
	logFilePath := ""
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/logs/"
	}
	if err := os.MkdirAll(logFilePath, 0777); err != nil {
		fmt.Println(err.Error())
	}
	logFileName := now.Format("2006-01-02-15") + ".log"
	//日志文件
	fileName := path.Join(logFilePath, logFileName)
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			fmt.Println(err.Error())
		}
	}
	//写入文件
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	fileLogger := logrus.New()
	fileLogger.SetOutput(src)
	return &FileHook{
		logger: fileLogger,
	}
}
func (hook *FileHook) Fire(entry *logrus.Entry) error {
	serialized, err := entry.Logger.Formatter.Format(entry)
	if err != nil {
		return err
	}
	hook.logger.Print(string(serialized))
	return nil
}

func (hook *FileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
