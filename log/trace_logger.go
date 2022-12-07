/**
 * @Author raven
 * @Description
 * @Date 2022/12/7
 **/
package log

import (
	"context"
	"github.com/sirupsen/logrus"
)

const TraceIdField = "trace-id"

type TraceLogger struct {
	log *logrus.Logger
}

func (d *TraceLogger) Error(ctx context.Context, msg string) {
	logrus.WithField(TraceIdField, getTraceId(ctx)).WithContext(ctx).Error(msg)
}
func (d *TraceLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	logrus.WithField(TraceIdField, getTraceId(ctx)).WithContext(ctx).Errorf(format, args...)
}
func (d *TraceLogger) Info(ctx context.Context, msg string) {
	logrus.WithField(TraceIdField, getTraceId(ctx)).WithContext(ctx).Info(msg)
}
func (d *TraceLogger) Infof(ctx context.Context, format string, arg ...interface{}) {
	logrus.WithField(TraceIdField, getTraceId(ctx)).WithContext(ctx).Infof(format, arg...)
}
func (d *TraceLogger) Warn(ctx context.Context, msg string) {
	logrus.WithField(TraceIdField, getTraceId(ctx)).WithContext(ctx).Warn(msg)
}
func (d *TraceLogger) Warnf(ctx context.Context, format string, arg ...interface{}) {
	logrus.WithField(TraceIdField, getTraceId(ctx)).WithContext(ctx).Warnf(format, arg...)
}

func (d *TraceLogger) Debug(ctx context.Context, msg string) {
	logrus.WithField(TraceIdField, getTraceId(ctx)).WithContext(ctx).Debug(msg)
}
func (d *TraceLogger) Debugf(ctx context.Context, format string, arg ...interface{}) {
	logrus.WithField(TraceIdField, getTraceId(ctx)).WithContext(ctx).Debugf(format, arg...)
}

func BuildTraceLogger() ILogger {
	return &TraceLogger{
		log: getLogrus(),
	}
}
func getTraceId(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	traceId := ctx.Value(TraceIdField)
	if traceId == nil {
		traceId = ""
	}
	return traceId.(string)
}
