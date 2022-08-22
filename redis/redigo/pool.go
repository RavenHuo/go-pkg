/**
 * @Author raven
 * @Description
 * @Date 2022/7/26
 **/
package redigo

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/RavenHuo/go-kit/redis/redigo/conn"
	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
)

type Pool struct {
	*redis.Pool
	name              string
	marshal           func(v interface{}) ([]byte, error)
	unmarshal         func(data []byte, v interface{}) error
	heartBeatInternal int
}

// New 根据配置参数创建redis工具实例
func NewPool(options *Options) (*Pool, error) {
	r := &Pool{}
	err := r.StartPool(options)
	return r, err
}

// StartPool 使用 Options 初始化redis，并在程序进程退出时关闭连接池。
func (p *Pool) StartPool(opts *Options) error {

	if opts.Network == "" {
		opts.Network = "tcp"
	}
	if opts.Addr == "" {
		opts.Addr = "127.0.0.1:6379"
	}
	if opts.MaxIdle == 0 {
		opts.MaxIdle = 3
	}
	if opts.IdleTimeout == 0 {
		opts.IdleTimeout = 300
	}
	if opts.Marshal == nil {
		p.marshal = json.Marshal
	}
	if opts.Unmarshal == nil {
		p.unmarshal = json.Unmarshal
	}
	if opts.HeartBeatInternal == 0 {
		opts.HeartBeatInternal = 10
	}
	pool := &redis.Pool{
		MaxActive:   opts.MaxActive,
		MaxIdle:     opts.MaxIdle,
		IdleTimeout: time.Duration(opts.IdleTimeout) * time.Second,

		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(opts.Network, opts.Addr)
			if err != nil {
				return nil, err
			}
			if opts.Password != "" {
				if _, err := conn.Do("AUTH", opts.Password); err != nil {
					conn.Close()
					return nil, err
				}
			}
			if _, err := conn.Do("SELECT", opts.Db); err != nil {
				conn.Close()
				return nil, err
			}
			return conn, err
		},

		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}
	p.Pool = pool
	p.heartBeatInternal = opts.HeartBeatInternal
	p.heartbeat()
	return nil
}

func (p *Pool) ClosePool() error {
	if p.Pool == nil {
		return errors.New("pool is nil")
	}
	return p.Pool.Close()
}

func (p *Pool) GetConn() (*conn.RedisConn, error) {
	redisConn := p.Pool.Get()
	if redisConn.Err() != nil {
		return nil, redisConn.Err()
	}
	return conn.WrapperRedisConn(redisConn), nil
}

func (p *Pool) heartbeat() {
	go func() {
		for {
			select {
			case <-time.After(time.Second * time.Duration(p.heartBeatInternal)):
				func() {
					redisConn, err := p.GetConn()
					if err != nil {
						logrus.Errorf("ping getConn err:%s", err)
						return
					}
					defer redisConn.Close()
					stringCmder := redisConn.Ping(context.Background())
					if stringCmder != nil && stringCmder.Err() != nil {
						logrus.Errorf("ping getConn err:%s", err)
					}
				}()
			}
		}
	}()
}
