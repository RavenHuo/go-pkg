/**
 * @Author raven
 * @Description
 * @Date 2021/9/15
 **/
package mongo

import (
	"github.com/sirupsen/logrus"
)

type Options struct {
	// 重试次数
	Retries int
	// 单个数据库链接的最大链接数
	PoolSize int
	// 多个数据库链接的最大链接数
	PoolMax int
	// 是否打印没有释放链接的行为
	PrintNoClose bool

	UserName string

	PassWord string

	Url string

	// 超时时间
	TimeOut int

	Logger *logrus.Logger
}

type Option func(options *Options)

func defaultOptions() *Options {
	return &Options{
		Retries:      3,
		PoolSize:     10,
		PoolMax:      20,
		PrintNoClose: true,
		UserName:     "raven",
		PassWord:     "123",
		Url:          "127.0.0.1:",
		TimeOut:      1000,
		Logger:       logrus.New(),
	}
}

func WithRetries(retry int) Option {
	return func(options *Options) {
		options.Retries = retry
	}
}

func WithPoolSize(poolSize int) Option {
	return func(options *Options) {
		options.PoolSize = poolSize
	}
}
func WithPoolMax(PoolMax int) Option {
	return func(options *Options) {
		options.PoolMax = PoolMax
	}
}
func WithUserName(UserName string) Option {
	return func(options *Options) {
		options.UserName = UserName
	}
}

func WithPassWord(PassWord string) Option {
	return func(options *Options) {
		options.PassWord = PassWord
	}
}

func WithUrl(Url string) Option {
	return func(options *Options) {
		options.Url = Url
	}
}

func WithTimeOut(TimeOut int) Option {
	return func(options *Options) {
		options.TimeOut = TimeOut
	}
}

func WithLogger(Logger *logrus.Logger) Option {
	return func(options *Options) {
		options.Logger = Logger
	}
}
