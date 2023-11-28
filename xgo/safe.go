package xgo

import (
	"context"
	"github.com/RavenHuo/go-pkg/log"
	"runtime"
	"runtime/debug"
)

func SafeFunc(f func(), ctx context.Context, panicLog string) {
	defer func() {
		if e := recover(); e != nil {
			_, file, line, _ := runtime.Caller(3)
			log.Errorf(ctx, "SafeFunc catch panic, business_log is %s", panicLog)
			log.Errorf(ctx, "recover. line:%s:%d, e:%v,stack:%s ", file, line, e, string(debug.Stack()))
		}
	}()
	f()
}

func SafeGo(f func(), ctx context.Context) {
	go func() {
		defer recoverContext(ctx)
		f()
	}()
}

func recoverContext(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf(ctx, "goroutine panic() err:%v, msg : %s", err, string(debug.Stack()))
			go func() {
				defer func() {
					if err := recover(); err != nil {
						log.Errorf(ctx, "goroutine panic() err : %s, msg :%s", err, string(debug.Stack()))
					}
				}()
			}()
		}
	}()
}
