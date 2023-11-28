package connpool

import (
	"context"
	"fmt"
	"github.com/RavenHuo/go-pkg/log"
	"net"
	"sync/atomic"
	"time"
)

var noDeadline = time.Time{}

type Conn struct {
	usedAt    int64     // 使用时间 atomic
	netConn   net.Conn  // 真正的连接
	Inited    bool      // 是否初始化，需要注意的是new了不代表是初始化了
	pooled    bool      // 是否在池子里面
	createdAt time.Time // 创建时间
}

func NewConn(netConn net.Conn) *Conn {
	cn := &Conn{
		netConn:   netConn,
		createdAt: time.Now(),
	}
	cn.SetUsedAt(time.Now())
	return cn
}

func (cn *Conn) Close() error {
	return cn.netConn.Close()
}
func (cn *Conn) RemoteAddr() net.Addr {
	if cn.netConn != nil {
		return cn.netConn.RemoteAddr()
	}
	return nil
}

func (cn *Conn) UsedAt() time.Time {
	unix := atomic.LoadInt64(&cn.usedAt)
	return time.Unix(unix, 0)
}

func (cn *Conn) SetUsedAt(tm time.Time) {
	atomic.StoreInt64(&cn.usedAt, tm.Unix())
}

func (cn *Conn) SetNetConn(netConn net.Conn) {
	cn.netConn = netConn
}

func (cn *Conn) GetNetConn() net.Conn {
	return cn.netConn
}

// 写
func (cn *Conn) Write(b []byte) (int, error) {
	return cn.netConn.Write(b)
}

// 读
func (cn *Conn) Read(b []byte) (int, error) {
	return cn.netConn.Read(b)
}

func (cn *Conn) ReadWithContext(ctx context.Context, timeout time.Duration) ([]byte, error) {
	if err := cn.netConn.SetReadDeadline(cn.deadline(ctx, timeout)); err != nil {
		return nil, err
	}
	buffer := make([]byte, 512)

	for {
		// 接收最大的数据字节数为512
		len, err := cn.netConn.Read(buffer)
		if err != nil {
			break
		}
		buffer = buffer[:len]
	}

	return buffer, nil
}

func (cn *Conn) isConnectionClosed() bool {
	conn := cn.netConn
	// 设置短暂的读取超时时间
	timeoutDuration := 1 * time.Second
	err := conn.SetReadDeadline(time.Now().Add(timeoutDuration))
	if err != nil {
		fmt.Println("Error setting read deadline:", err)
		return true
	}

	// 尝试从连接中读取一个字节
	_, err = conn.Read([]byte{})
	if err != nil {
		// 检查错误类型，如果是超时错误，说明连接已断开
		netErr, ok := err.(net.Error)
		if ok && netErr.Timeout() {
			return true
		}

		// 其他错误可能是连接真的断开了
		log.Errorf(context.Background(), "Error reading from connection: %s", err)
		return true
	}

	// 重置读取超时时间
	err = conn.SetReadDeadline(time.Time{})
	if err != nil {
		log.Errorf(context.Background(), "Error resetting read deadline: %s", err)
	}

	// 如果以上都通过，则连接仍然打开
	return false
}

func (cn *Conn) deadline(ctx context.Context, timeout time.Duration) time.Time {
	tm := time.Now()
	cn.SetUsedAt(tm)

	if timeout > 0 {
		tm = tm.Add(timeout)
	}

	if ctx != nil {
		deadline, ok := ctx.Deadline()
		if ok {
			if timeout == 0 {
				return deadline
			}
			if deadline.Before(tm) {
				return deadline
			}
			return tm
		}
	}

	if timeout > 0 {
		return tm
	}

	return noDeadline
}

func (cn *Conn) Ping() bool {

	conn, err := net.DialTimeout("ip:icmp", cn.RemoteAddr().String(), 2*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()

	return true
}
