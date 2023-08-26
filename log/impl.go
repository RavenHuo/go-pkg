/**
 * @Author raven
 * @Description
 * @Date 2022/12/6
 **/
package log

import "context"

func Debug(ctx context.Context, msg string) {
	logger.Debug(ctx, msg)
}
func Debugf(ctx context.Context, format string, v ...interface{}) {
	logger.Debugf(ctx, format, v...)
}

func Info(ctx context.Context, msg string) {
	logger.Info(ctx, msg)
}
func Infof(ctx context.Context, format string, v ...interface{}) {
	logger.Infof(ctx, format, v...)
}

func Warn(ctx context.Context, msg string) {
	logger.Warn(ctx, msg)
}
func Warnf(ctx context.Context, format string, v ...interface{}) {
	logger.Warnf(ctx, format, v...)
}

func Error(ctx context.Context, msg string) {
	logger.Error(ctx, msg)
}
func Errorf(ctx context.Context, format string, v ...interface{}) {
	logger.Errorf(ctx, format, v...)
}
