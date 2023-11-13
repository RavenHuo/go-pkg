package connpool

import (
	"context"
	"errors"
	"github.com/RavenHuo/go-pkg/log"
	"sync"
	"sync/atomic"
	"time"
)

var (
	// ErrClosed performs any operation on the closed client will return this error.
	ErrClosed = errors.New("redis: client is closed")

	// ErrPoolTimeout timed out waiting to get a connection from the connection pool.
	ErrPoolTimeout = errors.New("redis: connection pool timeout")
)

// ConnPool 连接池
// 1、核心连接数
// 2、最大连接数
// 3、连接健康检查(TODO)
// 4、空闲连接管理
// 5、连接达到最长时间的时候重置连接
type ConnPool struct {
	opt          *options
	queue        chan struct{} // 这里代表着所有的连接，chan里面能获取到struct，就代表连接池还有链接，无需等待，如果
	connsMu      *sync.RWMutex // 读写锁，提高性能
	conns        []*Conn
	idleConns    []*Conn
	poolSize     int    // 当前连接池大小
	idleConnsLen int    // 当前空闲连接数量
	_closed      uint32 // 关闭状态 atomic 0=open，1=close
	closedCh     chan struct{}
}

func New(opts ...Option) *ConnPool {
	opt := defaultOpt()
	for _, o := range opts {
		o(opt)
	}

	pool := &ConnPool{
		opt:       opt,
		queue:     make(chan struct{}),
		connsMu:   &sync.RWMutex{},
		conns:     make([]*Conn, 0, opt.poolSize),
		idleConns: make([]*Conn, 0, opt.poolSize),
		poolSize:  0,
		closedCh:  make(chan struct{}),
	}
	pool.connsMu.Lock()
	defer pool.connsMu.Unlock()
	pool.addMinIdleConns()

	if opt.idleTimeout > 0 && opt.idleCheckFrequency > 0 {
		go func() {
			pool.idleCheck()
		}()
	}

	return pool
}

// 判断是否关闭
func (p *ConnPool) closed() bool {
	return atomic.LoadUint32(&p._closed) == 1
}

// Close 关闭连接池
func (p *ConnPool) Close() error {
	if !atomic.CompareAndSwapUint32(&p._closed, 0, 1) {
		return ErrClosed
	}
	// 关闭channel
	close(p.closedCh)

	var firstErr error
	p.connsMu.Lock()
	defer p.connsMu.Unlock()
	for _, cn := range p.conns {
		if err := p.closeConn(cn); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	p.conns = nil
	p.poolSize = 0
	p.idleConns = nil
	p.idleConnsLen = 0

	return firstErr
}

// 关闭连接
func (p *ConnPool) closeConn(cn *Conn) error {
	if p.opt.onClose != nil && cn.netConn != nil {
		_ = p.opt.onClose(cn.netConn)
	}
	return cn.Close()
}

// 空闲连接检查
func (p *ConnPool) idleCheck() {
	ticker := time.NewTicker(p.opt.idleCheckFrequency)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 假如连接池是关闭了的，就结束这个定时任务
			if p.closed() {
				return
			}
			_, err := p.CheckIdleConns()
			if err != nil {
				log.Errorf(context.Background(), "CheckIdleConns failed: %s", err)
				continue
			}
		case <-p.closedCh:
			return
		}
	}
}

// CheckIdleConns 检查空闲连接
func (p *ConnPool) CheckIdleConns() (int, error) {
	var n int
	for {
		p.getTurn()

		p.connsMu.Lock()
		cn := p.RemoveIdleConn()
		p.connsMu.Unlock()

		p.freeTurn()

		// 这里cn不为空的话，那就是移除了一个连接
		if cn != nil {
			// 移除了一个连接之后需要判断是否满足最小连接池
			p.addMinIdleConns()

			_ = p.closeConn(cn)
			n++
		} else {
			break
		}
	}
	return n, nil
}

// 写queue
func (p *ConnPool) getTurn() {
	p.queue <- struct{}{}
}

// queue 出
func (p *ConnPool) freeTurn() {
	<-p.queue
}

// 判断链接是否可用，返回false的时候需要 重置连接
func (p *ConnPool) isStaleConn(cn *Conn) bool {
	if p.opt.idleTimeout == 0 && p.opt.maxConnAge == 0 {
		return false
	}

	now := time.Now()
	// 空闲时间 判断
	if p.opt.idleTimeout > 0 && now.Sub(cn.UsedAt()) >= p.opt.idleTimeout {
		return true
	}
	// 最长存活时间判断
	if p.opt.maxConnAge > 0 && now.Sub(cn.createdAt) >= p.opt.maxConnAge {
		return true
	}

	return false
}

// RemoveIdleConn 移除 空闲链接
func (p *ConnPool) RemoveIdleConn() *Conn {
	if len(p.idleConns) == 0 {
		return nil
	}

	// TODO 每次只检查第一个连接？
	cn := p.idleConns[0]
	if !p.isStaleConn(cn) {
		return nil
	}

	p.idleConns = append(p.idleConns[:0], p.idleConns[1:]...)
	p.idleConnsLen--
	p.removeConn(cn)

	return cn
}

// 从连接池中移除连接
func (p *ConnPool) removeConn(cn *Conn) {
	for i, c := range p.conns {
		// 注意这里需要使用指针才能判断出来
		if c == cn {
			p.conns = append(p.conns[:i], p.conns[i+1:]...)
			if cn.pooled {
				p.poolSize--
			}
			return
		}
	}
}

// 判断最小连接数
func (p *ConnPool) checkMinIdleConns() bool {
	if p.opt.minIdleConn == 0 {
		return false
	}
	if p.poolSize < p.opt.poolSize && p.idleConnsLen < p.opt.minIdleConn {
		return true
	}
	return false
}

// 添加最少的空闲链接
func (p *ConnPool) addMinIdleConns() {
	for p.checkMinIdleConns() {
		p.poolSize++
		p.idleConnsLen++

		go func() {
			err := p.addIdleConn()
			if err != nil && err != ErrClosed {
				p.connsMu.Lock()
				p.poolSize--
				p.idleConnsLen--
				p.connsMu.Unlock()
			}
		}()
	}
}

// 连接池添加连接
func (p *ConnPool) addIdleConn() error {
	cn, err := p.dialConn(context.TODO(), true)
	if err != nil {
		return err
	}

	p.connsMu.Lock()
	defer p.connsMu.Unlock()

	// It is not allowed to add new connections to the closed connection pool.
	if p.closed() {
		_ = cn.Close()
		return ErrClosed
	}

	p.conns = append(p.conns, cn)
	p.idleConns = append(p.idleConns, cn)
	return nil
}

// 创建链接
// @param pooled: 是否在池子里面
func (p *ConnPool) dialConn(ctx context.Context, pooled bool) (*Conn, error) {
	if p.closed() {
		return nil, ErrClosed
	}
	// 创建原生连接
	netConn, err := p.opt.onDialer(ctx)
	if err != nil {
		log.Errorf(ctx, "pool dialConn failed, pooled:%v, err:%s", pooled, err)
		return nil, err
	}
	cn := NewConn(netConn)
	cn.pooled = pooled
	return cn, nil
}

// Len 连接池长度
func (p *ConnPool) Len() int {
	p.connsMu.Lock()
	n := len(p.conns)
	p.connsMu.Unlock()
	return n
}

// IdleLen 空闲
func (p *ConnPool) IdleLen() int {
	p.connsMu.Lock()
	n := p.idleConnsLen
	p.connsMu.Unlock()
	return n
}
