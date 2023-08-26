/**
 * @Author raven
 * @Description
 * @Date 2022/7/7
 **/
package log

import (
	"context"
	"github.com/sirupsen/logrus"
)

type contextLogger struct {
	logRus *logrus.Logger
}

func (d *contextLogger) Error(ctx context.Context, msg string) {
	d.logRus.WithContext(ctx).Error(msg)
}
func (d *contextLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	d.logRus.WithContext(ctx).Errorf(format, args...)
}
func (d *contextLogger) Info(ctx context.Context, msg string) {
	d.logRus.WithContext(ctx).Info(msg)
}
func (d *contextLogger) Infof(ctx context.Context, format string, arg ...interface{}) {
	d.logRus.WithContext(ctx).Infof(format, arg...)
}
func (d *contextLogger) Warn(ctx context.Context, msg string) {
	d.logRus.WithContext(ctx).Warn(msg)
}
func (d *contextLogger) Warnf(ctx context.Context, format string, arg ...interface{}) {
	d.logRus.WithContext(ctx).Warnf(format, arg...)
}

func (d *contextLogger) Debug(ctx context.Context, msg string) {
	d.logRus.WithContext(ctx).Debug(msg)
}
func (d *contextLogger) Debugf(ctx context.Context, format string, arg ...interface{}) {
	d.logRus.WithContext(ctx).Debugf(format, arg...)
}

func defaultLogger() ILogger {
	return &contextLogger{
		logRus: getLogrus(),
	}
}
