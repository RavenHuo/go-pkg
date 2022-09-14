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

type Logger struct{}

func (d Logger) Error(ctx context.Context, msg string) {
	logrus.WithContext(ctx).Error(msg)
}
func (d Logger) Errorf(ctx context.Context, format string, args ...interface{}) {
	logrus.WithContext(ctx).Errorf(format, args...)
}
func (d Logger) Info(ctx context.Context, msg string) {
	logrus.WithContext(ctx).Info(msg)
}
func (d Logger) Infof(ctx context.Context, format string, arg ...interface{}) {
	logrus.WithContext(ctx).Infof(format, arg...)
}
func (d Logger) Warn(ctx context.Context, msg string) {
	logrus.WithContext(ctx).Warn(msg)
}
func (d Logger) Warnf(ctx context.Context, format string, arg ...interface{}) {
	logrus.WithContext(ctx).Warnf(format, arg...)
}

func (d Logger) Debug(ctx context.Context, msg string) {
	logrus.WithContext(ctx).Debug(msg)
}
func (d Logger) Debugf(ctx context.Context, format string, arg ...interface{}) {
	logrus.WithContext(ctx).Debugf(format, arg...)
}
