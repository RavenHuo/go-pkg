package connpool

import (
	"net"
	"sync/atomic"
	"time"
)

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

func (cn *Conn) Ping() bool {

	conn, err := net.DialTimeout("ip:icmp", cn.RemoteAddr().String(), 2*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()

	return true
}
