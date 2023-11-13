package connpool

import (
	"context"
	"net"
	"time"
)

type options struct {
	onDialer           func(context.Context) (net.Conn, error) // connection
	onClose            func(net.Conn) error                    // close
	poolSize           int                                     // 连接池的长度,  最少连接数<=连接池的长度<=最大连接数
	minIdleConn        int                                     // 最少连接数
	maxIdleConn        int                                     // 最大连接数
	poolTimeout        time.Duration                           // 连接池超时时间
	idleTimeout        time.Duration                           // 链接空闲时间
	idleCheckFrequency time.Duration                           // 心跳检查间隔时间
	maxConnAge         time.Duration                           // 最长的链接时间
}

func defaultOpt() *options {
	return &options{}
}

type Option func(options *options)
