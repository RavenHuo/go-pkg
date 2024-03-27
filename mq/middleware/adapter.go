package middleware

import (
	"context"
	mq "github.com/RavenHuo/go-pkg/mq"
)

// WrapperMultiHandler 中间件将多个 handler 合并成一个
func WrapperMultiHandler(handlers ...mq.EventHandler) mq.Handler {
	return func(eventInfo *mq.EventInfo) error {
		manager := GetHookManager()
		return Adapt(eventInfo, manager, handlers...)
	}
}

func Adapt(event *mq.EventInfo, hookManager *HookManager, handlers ...mq.EventHandler) error {
	bCtx := context.Background()
	for _, h := range handlers {
		preCtx, err := hookManager.BeforeRequest(bCtx, event, h)
		if err != nil {
			continue
		}
		err = h.Handler(event)
		hookManager.AfterRequest(preCtx, event, h, err)
	}
	return nil
}
