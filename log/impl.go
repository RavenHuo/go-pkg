/**
 * @Author raven
 * @Description
 * @Date 2022/12/6
 **/
package log

import "context"

func Debug(ctx context.Context, msg string) {
	getLogger().Debug(ctx, msg)
}
func Debugf(ctx context.Context, format string, v ...interface{}) {
	getLogger().Debugf(ctx, format, v...)
}

func Info(ctx context.Context, msg string) {
	getLogger().Info(ctx, msg)
}
func Infof(ctx context.Context, format string, v ...interface{}) {
	getLogger().Infof(ctx, format, v...)
}

func Warn(ctx context.Context, msg string) {
	getLogger().Warn(ctx, msg)
}
func Warnf(ctx context.Context, format string, v ...interface{}) {
	getLogger().Warnf(ctx, format, v...)
}

func Error(ctx context.Context, msg string) {
	getLogger().Error(ctx, msg)
}
func Errorf(ctx context.Context, format string, v ...interface{}) {
	getLogger().Errorf(ctx, format, v...)
}
