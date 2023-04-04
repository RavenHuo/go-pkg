/**
 * @Author raven
 * @Description
 * @Date 2021/8/27
 **/
package mongo

import (
	"errors"
	"github.com/sirupsen/logrus"
	"sync"
	"sync/atomic"
	"time"

	"gopkg.in/mgo.v2"
)

var clientPool ClientPool

var once = sync.Once{}

//ClientPool
type ClientPool struct {
	pools           map[string][]*Client
	lock            sync.Mutex
	totalUsing      int32
	poolUsing       map[string]int32
	retries         int
	retryIntervalMs int
	confPoolSize    int
	poolSize        int
	poolMax         int
	printNoClose    bool

	userName string
	passWord string
	url      string

	// 超时时间
	timeOut int
	logger  *logrus.Logger
}

//Client
type Client struct {
	session *mgo.Session
	url     string
	closeCh chan int
}

func NewMongo(opts ...Option) {
	once.Do(func() {
		options := defaultOptions()
		for _, opt := range opts {
			opt(options)
		}

		clientPool = ClientPool{
			confPoolSize: options.PoolSize,
			poolMax:      options.PoolMax,
			poolSize:     options.PoolSize,
			retries:      options.Retries,
			pools:        make(map[string][]*Client),
			poolUsing:    make(map[string]int32),
			passWord:     options.PassWord,
			userName:     options.UserName,
			timeOut:      options.TimeOut,
			url:          options.Url,
			logger:       options.Logger,
		}
	})
}

// getClient  获取 client对象
func (pool *ClientPool) getClient(dbName string, mode mgo.Mode) (*Client, error) {
	if dbName == "" {
		errStr := "mongodb database name is empty"
		pool.logger.Errorf(errStr)
		return nil, errors.New(errStr)
	}

	//从连接池获取
	client, err := pool.getMongoClientWithRetry(dbName)
	if client == nil || client.session == nil {
		return nil, err
	}

	client.session.SetMode(mode, true)
	return client, err
}

// getMongoClientWithRetry 重试的获取client
func (pool *ClientPool) getMongoClientWithRetry(dbName string) (*Client, error) {
	client, err := pool.getMongoClient(dbName)
	if !IsErrMaxConnectionLimited(err) {
		return client, err
	}

	// 重试 去获取链接
	interval := pool.retryIntervalMs
	for i := 0; i < pool.retries; i++ {
		if interval > 0 {
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
		client, err = pool.getMongoClient(dbName)
		if !IsErrMaxConnectionLimited(err) {
			break
		}

		interval *= 2

		pool.lock.Lock()
		using := pool.poolUsing[dbName]
		remain := len(pool.pools[dbName])
		total := pool.totalUsing
		pool.lock.Unlock()

		if i+1 == pool.retries {
			pool.logger.Warnf("reach max pool size, retry=%v, Url=%v, using=%v, remain=%v, total_using=%v", i, dbName, using, remain, total)
		}
	}
	return client, err
}

// 获取client
func (pool *ClientPool) getMongoClient(dbName string) (*Client, error) {
	pool.lock.Lock()
	atomic.AddInt32(&pool.totalUsing, 1)
	pool.poolUsing[dbName] += 1

	//超过连接数限制
	if int(atomic.LoadInt32(&pool.totalUsing)) > pool.poolMax || int(pool.poolUsing[dbName]) > pool.poolSize {
		atomic.AddInt32(&pool.totalUsing, -1)
		pool.poolUsing[dbName] -= 1
		pool.lock.Unlock()

		return nil, ErrMaxConnectionLimited
	}

	var client *Client
	var array = pool.pools[dbName]
	if array == nil {
		array = make([]*Client, 0, pool.poolSize)
		pool.pools[dbName] = array
	}
	var size = len(array)

	if size > 0 {
		client = array[size-1]
		pool.pools[dbName] = array[:size-1]
	}

	pool.lock.Unlock()

	if client == nil {
		session, err := pool.createMongoSession()
		if err != nil {
			pool.logger.Errorf("#mongo# dial server %s error: %v", dbName, err)
			pool.lock.Lock()
			atomic.AddInt32(&pool.totalUsing, -1)
			pool.poolUsing[dbName] -= 1
			pool.lock.Unlock()
			return nil, err
		} else {
			pool.logger.Infof("#mongo# server connected: Url=%s total_using:%d", dbName, pool.totalUsing)
			client = new(Client)
			client.session = session
			client.url = dbName
			client.closeCh = make(chan int)
		}
	}

	return client, nil
}

// 创建mongo session
func (pool *ClientPool) createMongoSession() (*mgo.Session, error) {
	info := &mgo.DialInfo{
		Addrs:    []string{pool.url},
		Password: pool.passWord,
		Username: pool.userName,
		Timeout:  time.Second * time.Duration(pool.timeOut),
	}
	session, err := mgo.DialWithInfo(info)
	return session, err
}

func (client *Client) Close() {
	if client.session == nil {
		return
	}
	clientPool.closeClient(client)
}

// 关闭链接
func (pool *ClientPool) closeClient(client *Client) {
	pool.lock.Lock()
	defer pool.lock.Unlock()

	if client == nil || client.session == nil {
		return
	}

	atomic.AddInt32(&pool.totalUsing, -1)
	pool.poolUsing[client.url] -= 1
	select {
	case client.closeCh <- 1:
	default:
	}

	//关闭连接
	var array = pool.pools[client.url]
	if len(array) >= pool.confPoolSize {
		client.session.Close()
		client.session = nil
		pool.logger.Infof("#mongo# server disconnected: Url=%s using:%d client len:%d total_using:%d", client.url, pool.poolUsing[client.url], len(array), pool.totalUsing)
		return
	}

	//保存连接
	pool.pools[client.url] = append(array, client)
}
