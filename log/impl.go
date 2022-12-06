/**
 * @Author raven
 * @Description
 * @Date 2022/12/6
 **/
package log

import "context"

func Debug(ctx context.Context, msg string) {
	getContextLogger().Debug(ctx, msg)
}
func Debugf(ctx context.Context, format string, v ...interface{}) {
	getContextLogger().Debugf(ctx, format, v...)
}

func Info(ctx context.Context, msg string) {
	getContextLogger().Info(ctx, msg)
}
func Infof(ctx context.Context, format string, v ...interface{}) {
	getContextLogger().Infof(ctx, format, v...)
}

func Warn(ctx context.Context, msg string) {
	getContextLogger().Warn(ctx, msg)
}
func Warnf(ctx context.Context, format string, v ...interface{}) {
	getContextLogger().Warnf(ctx, format, v...)
}

func Error(ctx context.Context, msg string) {
	getContextLogger().Error(ctx, msg)
}
func Errorf(ctx context.Context, format string, v ...interface{}) {
	getContextLogger().Errorf(ctx, format, v...)
}
