## **go-pkg**



golang的共用pak 库，用于支持各种基础功能，包括但不限于以下功能



- [Conf](#Conf)
- [Encode](#Encode)
- [Redis](#Redis)
- [Log](#Log)
- [Mongo](#Mongo)
- [Etcd](#Etcd)

------------------------

## Installation

The only requirement is the [Go Programming Language](https://golang.org/dl/)

```
go get -u github.com/RavenHuo/go-pkg
```



------------------

## Usage



## Conf

用于读取配置文件，获取配置

```
//初始化配置
InitConf()
// 获取命名 server.name 的配置
Get("server.name")
// 获取类型为string 命名为name的配置
GetString("name")
// 获取类型为int 命名为name的配置
GetInt("name")
// 获取命名为name的配置，并序列化为registerConfig对象
GetObject("register", &registerConfig)
```



## Encode

通过策略模式，实现不同格式的序列化以及反序列化，包括以下格式：

- ini

- json

- toml

- xml

- yaml

```
// Encoder represents a format encoder
type Encoder interface {
    Encode(interface{}) ([]byte, error)
    Decode([]byte, interface{}) error
    Name() string
}
  
```

  

## Redis

- go-redis

封装 github.com/go-redis/redis/v8，加入Lock及UnLock方法

- distribution_lock

分布式锁，参照redis-sission框架，实现了golang的分布式锁

- 使用redis的setnx 实现分布式加锁

- 使用watch dog 协程监控及redis的eval，实现了分布式锁的续期

- 使用redis的eval 实现分布式锁的解锁



## Log

基于 [logrus](https://github.com/sirupsen/logrus)的日志封装，将context透传到日志打印中，方便打印trace_id



## Mongo

基于[mgo](gopkg.in/mgo.v2)的 mongo 连接池封装



## Etcd

使用选项模式，基于[etcdV3](go.etcd.io/etcd/client/v3)的 etcd 链接封装

实现了Get，PUt，Delete，Watch 等功能

