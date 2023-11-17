package connpool

import (
	"context"
	"fmt"
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
		go goPing(pool, &wg)
	}
	wg.Wait()
	runCmd(pool)
}

func runCmd(pool *ConnPool) {
	conn, err := pool.Get(context.Background())
	if err != nil {
		log.Errorf(context.Background(), "get pool failed, err:%s", err)
		return
	}
	cmd := "Get raven"
	_, err = conn.Write([]byte(cmd))
	if err != nil {
		log.Infof(context.Background(), "write failed ")
		return
	}
	conn.netConn.SetReadDeadline(time.Now().Add(2 * time.Second))
	go func() {
		for {
			// 接收最大的数据字节数为512
			buf := make([]byte, 512)
			len, err := conn.Read(buf)
			if err != nil {
				fmt.Println("Error reading", err.Error())
				return //终止程序
			}
			fmt.Printf("Received data: %v\n", string(buf[:len]))
		}
	}()
	time.Sleep(2 * time.Second)
}

func goPing(pool *ConnPool, wg *sync.WaitGroup) {
	defer wg.Done()
	conn, err := pool.Get(context.Background())
	if err != nil {
		log.Errorf(context.Background(), "get pool failed, err:%s", err)
		return
	}
	log.Infof(context.Background(), "ping localAddr:%s ,success:%+v", conn.GetNetConn().LocalAddr(), conn.isConnectionClosed())
	pool.Put(context.Background(), conn)
}
