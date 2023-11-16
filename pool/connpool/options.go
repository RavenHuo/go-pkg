package connpool

import (
	"context"
	"github.com/RavenHuo/go-pkg/log"
	"net"
	"time"
)

type options struct {
	onDialer           func(context.Context) (net.Conn, error) // connection 执行方法
	onClose            func(net.Conn) error                    // close 执行方法
	queueSize          int                                     // 等待队列长度
	minConn            int                                     // 最少连接数
	maxConn            int                                     // 最大连接数
	waitTimeout        time.Duration                           // 连接池超时等待时间，
	idleTimeout        time.Duration                           // 链接空闲时间
	idleCheckFrequency time.Duration                           // 心跳检查间隔时间
	maxConnAge         time.Duration                           // 最长的链接时间
}

func defaultOpt() *options {
	return &options{
		onDialer: func(ctx context.Context) (net.Conn, error) {
			return net.Dial("tcp", "")
		},
		onClose: func(conn net.Conn) error {
			log.Infof(context.Background(), "conn close :%s", conn.LocalAddr())
			return nil
		},
		queueSize:          10,
		minConn:            5,
		idleCheckFrequency: time.Minute,
		waitTimeout:        time.Second,
		idleTimeout:        time.Minute,
	}
}

type Option func(options *options)

func WithOnDialer(onDialer func(context.Context) (net.Conn, error)) Option {
	return func(options *options) {
		options.onDialer = onDialer
	}
}
func WithOnClose(onClose func(net.Conn) error) Option {
	return func(options *options) {
		options.onClose = onClose
	}
}

func WithQueueSize(poolSize int) Option {
	return func(options *options) {
		options.queueSize = poolSize
	}
}

func WithMinConn(minIdleConn int) Option {
	return func(options *options) {
		options.minConn = minIdleConn
	}
}

func WithMaxConn(maxConn int) Option {
	return func(options *options) {
		options.maxConn = maxConn
	}
}

func WithWaitTimeout(waitTimeout time.Duration) Option {
	return func(options *options) {
		options.waitTimeout = waitTimeout
	}
}

func WithIdleTimeout(idleTimeout time.Duration) Option {
	return func(options *options) {
		options.idleTimeout = idleTimeout
	}
}

func WithIdleCheckFrequency(idleCheckFrequency time.Duration) Option {
	return func(options *options) {
		options.idleCheckFrequency = idleCheckFrequency
	}
}

func WithMaxConnAge(maxConnAge time.Duration) Option {
	return func(options *options) {
		options.maxConnAge = maxConnAge
	}
}
