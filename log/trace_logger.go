/**
 * @Author raven
 * @Description
 * @Date 2022/12/7
 **/
package log

import (
	"context"
	"github.com/RavenHuo/go-kit/trace"
	"github.com/sirupsen/logrus"
)

type TraceLogger struct {
	logRus *logrus.Logger
}

func (d *TraceLogger) Error(ctx context.Context, msg string) {
	d.logRus.WithField(trace.TraceIdField, trace.GetTraceId(ctx)).WithContext(ctx).Error(msg)
}
func (d *TraceLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	d.logRus.WithField(trace.TraceIdField, trace.GetTraceId(ctx)).WithContext(ctx).Errorf(format, args...)
}
func (d *TraceLogger) Info(ctx context.Context, msg string) {
	d.logRus.WithField(trace.TraceIdField, trace.GetTraceId(ctx)).WithContext(ctx).Info(msg)
}
func (d *TraceLogger) Infof(ctx context.Context, format string, arg ...interface{}) {
	d.logRus.WithField(trace.TraceIdField, trace.GetTraceId(ctx)).WithContext(ctx).Infof(format, arg...)
}
func (d *TraceLogger) Warn(ctx context.Context, msg string) {
	d.logRus.WithField(trace.TraceIdField, trace.GetTraceId(ctx)).WithContext(ctx).Warn(msg)
}
func (d *TraceLogger) Warnf(ctx context.Context, format string, arg ...interface{}) {
	d.logRus.WithField(trace.TraceIdField, trace.GetTraceId(ctx)).WithContext(ctx).Warnf(format, arg...)
}

func (d *TraceLogger) Debug(ctx context.Context, msg string) {
	d.logRus.WithField(trace.TraceIdField, trace.GetTraceId(ctx)).WithContext(ctx).Debug(msg)
}
func (d *TraceLogger) Debugf(ctx context.Context, format string, arg ...interface{}) {
	d.logRus.WithField(trace.TraceIdField, trace.GetTraceId(ctx)).WithContext(ctx).Debugf(format, arg...)
}

func BuildTraceLogger() ILogger {
	return &TraceLogger{
		logRus: getLogrus(),
	}
}
