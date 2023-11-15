package connpool

import (
	"context"
	"github.com/RavenHuo/go-pkg/log"
	"sync"
	"testing"
	"time"
)

func TestNewPool(t *testing.T) {
	opts := make([]Option, 0)
	opts = append(opts, WithMaxConn(10))
	opts = append(opts, WithWaitTimeout(2*time.Second))
	pool := New(opts...)
	wg := sync.WaitGroup{}
	defer pool.Close()
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go ping(pool, &wg)
	}
	wg.Wait()
}

func ping(pool *ConnPool, wg *sync.WaitGroup) {
	defer wg.Done()
	conn, err := pool.Get(context.Background())
	if err != nil {
		log.Errorf(context.Background(), "get pool failed, err:%s", err)
		return
	}
	conn.Ping()
	log.Infof(context.Background(), "ping success localAddr:%s", conn.GetNetConn().LocalAddr())
	time.Sleep(time.Second)
	pool.Put(context.Background(), conn)
}
