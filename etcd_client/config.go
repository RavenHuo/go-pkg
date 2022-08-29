/**
 * @Author raven
 * @Description
 * @Date 2022/8/29
 **/
package etcd_client

import "time"

const (
	DefaultDialTimeoutSecond = 10
)

type EtcdConfig struct {
	Endpoints []string `json:"endpoints" yaml:"endpoints"`
	CertFile  string   `json:"certFile" yaml:"certFile"`
	KeyFile   string   `json:"keyFile" yaml:"keyFile"`
	CaCert    string   `json:"caCert" yaml:"caCert"`
	BasicAuth bool     `json:"basicAuth" yaml:"basicAuth"`
	UserName  string   `json:"userName" yaml:"userName"`
	Password  string   `json:"-" yaml:""`
	// 连接超时时间
	ConnectTimeout time.Duration `json:"connectTimeout" yaml:"connectTimeout"`
	Secure         bool          `json:"secure" yaml:"secure"`
	// 自动同步member list的间隔
	AutoSyncInterval time.Duration `json:"autoAsyncInterval" yaml:"autoAsyncInterval"`
	TTL              int           `json:"ttl" yaml:"ttl"` // 单位：s
}
