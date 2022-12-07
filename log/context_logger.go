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
	log *logrus.Logger
}

func (d *contextLogger) Error(ctx context.Context, msg string) {
	d.log.WithContext(ctx).Error(msg)
}
func (d *contextLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	d.log.WithContext(ctx).Errorf(format, args...)
}
func (d *contextLogger) Info(ctx context.Context, msg string) {
	d.log.WithContext(ctx).Info(msg)
}
func (d *contextLogger) Infof(ctx context.Context, format string, arg ...interface{}) {
	d.log.WithContext(ctx).Infof(format, arg...)
}
func (d *contextLogger) Warn(ctx context.Context, msg string) {
	d.log.WithContext(ctx).Warn(msg)
}
func (d *contextLogger) Warnf(ctx context.Context, format string, arg ...interface{}) {
	d.log.WithContext(ctx).Warnf(format, arg...)
}

func (d *contextLogger) Debug(ctx context.Context, msg string) {
	d.log.WithContext(ctx).Debug(msg)
}
func (d *contextLogger) Debugf(ctx context.Context, format string, arg ...interface{}) {
	d.log.WithContext(ctx).Debugf(format, arg...)
}
