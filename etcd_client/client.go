/**
 * @Author raven
 * @Description
 * @Date 2022/8/29
 **/
package etcd_client

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/RavenHuo/go-kit/log"
	"github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc"
	"io/ioutil"
	"time"
)

var ConfigErr = errors.New("init etcd config error")

type Client struct {
	c      *clientv3.Client
	config *EtcdConfig
	logger log.ILogger
}

func New(config *EtcdConfig, logger log.ILogger) (*Client, error) {
	if config == nil || len(config.Endpoints) == 0 {
		return nil, ConfigErr
	}
	if config.ConnectTimeout == 0 {
		config.ConnectTimeout = DefaultDialTimeoutSecond * time.Second
	}
	etcdClient, err := genClient(config)
	if err != nil {
		return nil, err
	}
	return &Client{
		c:      etcdClient,
		config: config,
		logger: logger,
	}, nil
}

func genClient(config *EtcdConfig) (*clientv3.Client, error) {
	conf := clientv3.Config{
		Endpoints:            config.Endpoints,
		DialTimeout:          config.ConnectTimeout,
		DialKeepAliveTime:    10 * time.Second,
		DialKeepAliveTimeout: 3 * time.Second,
		DialOptions: []grpc.DialOption{
			grpc.WithBlock(),
		},
		AutoSyncInterval: config.AutoSyncInterval,
	}

	if !config.Secure {
		conf.DialOptions = append(conf.DialOptions, grpc.WithInsecure())
	}

	if config.BasicAuth {
		conf.Username = config.UserName
		conf.Password = config.Password
	}

	tlsEnabled := false
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
	}

	if config.CaCert != "" {
		certBytes, err := ioutil.ReadFile(config.CaCert)
		if err != nil {
			return nil, err
		}

		caCertPool := x509.NewCertPool()
		ok := caCertPool.AppendCertsFromPEM(certBytes)

		if ok {
			tlsConfig.RootCAs = caCertPool
		}
		tlsEnabled = true
	}

	if config.CertFile != "" && config.KeyFile != "" {
		tlsCert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{tlsCert}
		tlsEnabled = true
	}

	if tlsEnabled {
		conf.TLS = tlsConfig
	}

	return clientv3.New(conf)
}
