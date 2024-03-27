package middleware

import (
	"context"
	"github.com/RavenHuo/go-pkg/log"
	"github.com/RavenHuo/go-pkg/mq"
	"time"
)

var (
	hookManager *HookManager
)

func init() {
	hookManager = &HookManager{}
	Register(&LogHook{})
}

func GetHookManager() *HookManager {
	return hookManager
}

func Register(hook Hook) {
	hookManager.Register(hook)
}

type Hook interface {
	BeforeConsumer(ctx context.Context, info *mq.EventInfo, handler mq.EventHandler) (context.Context, error)
	AfterConsumer(ctx context.Context, info *mq.EventInfo, handler mq.EventHandler, err error)
}

type HookManager struct {
	hooks []Hook
}

func (manager *HookManager) Register(hook Hook) {
	manager.hooks = append(manager.hooks, hook)
}
func (manager *HookManager) BeforeRequest(ctx context.Context, info *mq.EventInfo, handler mq.EventHandler) (context.Context, error) {
	var err error
	for _, h := range manager.hooks {
		ctx, err = h.BeforeConsumer(ctx, info, handler)
		if err != nil {
			return nil, err
		}
	}
	return ctx, err
}

func (manager *HookManager) AfterRequest(ctx context.Context, info *mq.EventInfo, handler mq.EventHandler, err error) {
	for i := len(manager.hooks) - 1; i >= 0; i-- {
		manager.hooks[i].AfterConsumer(ctx, info, handler, err)
	}
}

const (
	timeCostCtxKey = "mq_start"
)

type LogHook struct {
}

func (l *LogHook) BeforeConsumer(ctx context.Context, info *mq.EventInfo, handler mq.EventHandler) (context.Context, error) {
	ctx = context.WithValue(ctx, timeCostCtxKey, time.Now())
	return ctx, nil
}

func (l *LogHook) AfterConsumer(ctx context.Context, info *mq.EventInfo, handler mq.EventHandler, err error) {
	requestStartCtx := ctx.Value(timeCostCtxKey)
	handlerName := handler.Name()
	if requestTime, ok := requestStartCtx.(time.Time); ok {
		cost := time.Now().Sub(requestTime).Milliseconds()
		log.Infof(ctx, "[kafka] start:%d detail:%v,handlerName:%s cost:%v ms ", requestTime.UnixNano()/1e6, info, handlerName, cost)
	}
	if err != nil {
		log.Infof(ctx, "[kafka] detail:%v,handlerName:%s,err:%s ", info, handlerName, err)
	}
}
