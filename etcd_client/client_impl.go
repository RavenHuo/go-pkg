/**
 * @Author raven
 * @Description
 * @Date 2022/8/29
 **/
package etcd_client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc"
	"io/ioutil"
	"time"
)

// 获取目录里面所有信息
func (client *Client) GetDirectory(ctx context.Context, directory string) (res map[string][]byte, err error) {
	kvMap := make(map[string][]byte)
	kv := clientv3.NewKV(client.c)
	rsp, getErr := kv.Get(ctx, directory, clientv3.WithPrefix())
	if getErr != nil {
		return nil, getErr
	}
	for _, item := range rsp.Kvs {
		kvMap[string(item.Key)] = item.Value
	}
	return kvMap, err
}

// 给key设置value
func (client *Client) PutKey(key string, value string, expireTTL int) (leaseId int64, err error) {
	ctx := context.Background()
	kv := clientv3.NewKV(client.c)

	if expireTTL > 0 {
		// 创建一个租约
		resp, err := client.c.Grant(context.TODO(), int64(expireTTL))
		if err != nil {
			return 0, err
		}

		if _, putErr := kv.Put(ctx, key, value, clientv3.WithLease(resp.ID)); putErr != nil {
			return 0, putErr
		}

		return int64(resp.ID), nil
	}

	if _, putErr := kv.Put(ctx, key, value); putErr != nil {
		return 0, putErr
	}

	return 0, nil
}

// 删除一个key
func (client *Client) DeleteKey(key string) (err error) {
	ctx := context.Background()
	kv := clientv3.NewKV(client.c)
	_, deleteErr := kv.Delete(ctx, key)
	if deleteErr != nil {
		return deleteErr
	}
	return nil
}

func (client *Client) genClient() (err error) {
	config := client.config
	conf := clientv3.Config{
		Endpoints:            config.Endpoints,
		DialTimeout:          config.ConnectTimeout,
		DialKeepAliveTime:    10 * time.Second,
		DialKeepAliveTimeout: 3 * time.Second,
		DialOptions: []grpc.DialOption{
			grpc.WithBlock(),
			// go-grpc-prometheus 后续接入
			//grpc.WithUnaryInterceptor(grpcprom.UnaryClientInterceptor),
			//grpc.WithStreamInterceptor(grpcprom.StreamClientInterceptor),
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
			return err
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
			return err
		}
		tlsConfig.Certificates = []tls.Certificate{tlsCert}
		tlsEnabled = true
	}

	if tlsEnabled {
		conf.TLS = tlsConfig
	}

	if client.c, err = clientv3.New(conf); err != nil {
		return err
	}

	return nil
}

func (client *Client) Close() (err error) {
	return client.c.Close()
}

// 获取原始client ptr
func (client *Client) GetClient() *clientv3.Client {
	return client.c
}