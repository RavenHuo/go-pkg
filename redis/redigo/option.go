/**
 * @Author raven
 * @Description
 * @Date 2022/7/26
 **/
package redigo

// Options redis配置参数
type Options struct {
	Network     string                                 // 通讯协议，默认为 tcp
	Addr        string                                 // redis服务的地址，默认为 127.0.0.1:6379
	Password    string                                 // redis鉴权密码
	Db          int                                    // 数据库
	MaxActive   int                                    // 最大活动连接数，值为0时表示不限制
	MaxIdle     int                                    // 最大空闲连接数
	IdleTimeout int                                    // 空闲连接的超时时间，超过该时间则关闭连接。单位为秒。默认值是5分钟。值为0时表示不关闭空闲连接。此值应该总是大于redis服务的超时时间。
	Prefix      string                                 // 键名前缀
	Marshal     func(v interface{}) ([]byte, error)    // 数据序列化方法，默认使用json.Marshal序列化
	Unmarshal   func(data []byte, v interface{}) error // 数据反序列化方法，默认使用json.Unmarshal序列化

	HeartBeatInternal int // 心跳时间 Second

	ConnectTimeout int // 链接超时 Millisecond
	ReadTimeout    int // 读超时 Millisecond
	WriteTimeout   int // 写超时 Millisecond
}
