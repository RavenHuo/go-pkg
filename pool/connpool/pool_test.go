package connpool

import (
	"context"
	"github.com/RavenHuo/go-pkg/log"
	"strconv"
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
		go runCmd(pool, &wg)
		// go goPing(pool, &wg)
	}
	wg.Wait()
	time.Sleep(5 * time.Second)
}

func runCmd(pool *ConnPool, wg *sync.WaitGroup) {
	defer wg.Done()
	conn, err := pool.Get(context.Background())
	if err != nil {
		log.Errorf(context.Background(), "get pool failed, err:%s", err)
		return
	}
	cmd := "Get raven " + strconv.Itoa(int(time.Now().UnixNano()))
	_, err = conn.Write([]byte(cmd))
	if err != nil {
		log.Errorf(context.Background(), "write failed %s", err)
		return
	}
	readByte, err := conn.ReadWithContext(context.Background(), time.Second)
	if err != nil {
		log.Errorf(context.Background(), "read failed %s", err)
		return
	}
	log.Infof(context.Background(), "read success writeByte:%s, readByte:%s", cmd, string(readByte))
}

func goPing(pool *ConnPool, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Second)
	conn, err := pool.Get(context.Background())
	if err != nil {
		log.Errorf(context.Background(), "get pool failed, err:%s", err)
		return
	}

	log.Infof(context.Background(), "ping localAddr:%s ,success:%+v", conn.GetNetConn().LocalAddr(), conn.isConnectionClosed())
	pool.Put(context.Background(), conn)
}
