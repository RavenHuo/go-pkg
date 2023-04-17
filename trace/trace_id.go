/**
 * @Author raven
 * @Description
 * @Date 2023/4/17
 **/
package trace

import (
	"context"
	"github.com/RavenHuo/go-kit/utils"
)

const TraceIdField = "trace-id"

func GetTraceId(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	traceId := ctx.Value(TraceIdField)
	if traceId == nil {
		traceId = "main"
	}
	return traceId.(string)
}

func SetTraceId(ctx context.Context, traceId string) context.Context {
	ctx = context.WithValue(ctx, TraceIdField, traceId)
	return ctx
}

func GenTraceId(ctx context.Context) context.Context {
	traceId := utils.GetUuid()
	ctx = context.WithValue(ctx, TraceIdField, traceId)
	return ctx
}
