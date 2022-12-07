/**
 * @Author raven
 * @Description
 * @Date 2022/8/29
 **/
package log

import (
	"context"
)

var logger ILogger

type ILogger interface {
	Debug(ctx context.Context, msg string)
	Debugf(ctx context.Context, format string, v ...interface{})

	Info(ctx context.Context, msg string)
	Infof(ctx context.Context, format string, v ...interface{})

	Warn(ctx context.Context, msg string)
	Warnf(ctx context.Context, format string, v ...interface{})

	Error(ctx context.Context, msg string)
	Errorf(ctx context.Context, format string, v ...interface{})
}

func getLogger() ILogger {
	return logger
}

func init() {
	logger = &contextLogger{
		log: getLogrus(),
	}
}
