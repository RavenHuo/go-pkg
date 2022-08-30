package etcd_client

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"time"
)

type WatchHandler func(ctx context.Context, event *Event)

type Event struct {
	Type mvccpb.Event_EventType
	Kv   *KeyValue
}

type KeyValue struct {
	Key   []byte
	Value []byte
}

// Watch A watch only tells the latest revision
type Watch struct {
	// 对于观看进度响应， header.revision 指示进度。保证在此流中收到的所有未来事件的修订号都高于 header.revision 号。
	// 用于版本控制
	revision     int64
	cancel       context.CancelFunc
	eventChan    chan *Event
	incipientKVs []*KeyValue
}

// C ...
func (w *Watch) C() chan *Event {
	return w.eventChan
}

// 监听前缀key，并返回chan （非阻塞执行）
func (client *client) WatchPrefix(prefixContext context.Context, prefix string) (chan *Event, error) {
	resp, err := client.c.Get(prefixContext, prefix, clientv3.WithPrefix(), clientv3.WithKeysOnly())
	if err != nil {
		return nil, err
	}

	kvs := resp.Kvs
	reversion := resp.Header.Revision

	incipientKVs := make([]*KeyValue, len(resp.Kvs))
	start := time.Now()
	for index, kv := range kvs {
		resp, err = client.c.Get(prefixContext, string(kv.Key))
		if err != nil || resp == nil || len(resp.Kvs) != 1 {
			continue
		}

		incipientKVs[index] = &KeyValue{Key: kv.Key, Value: resp.Kvs[0].Value}
		reversion = resp.Header.Revision
	}

	w := &Watch{
		revision:  reversion,
		eventChan: make(chan *Event, 100),
	}
	client.logger.Infof(prefixContext, "init get prefix:%s cost:%dms", prefix, time.Now().Sub(start).Milliseconds())

	w.incipientKVs = incipientKVs

	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		w.cancel = cancel
		rch := client.c.Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithCreatedNotify(), clientv3.WithRev(w.revision))
		for {
			for n := range rch {
				// etcd的最小版本
				if n.CompactRevision > w.revision {
					w.revision = n.CompactRevision
				}
				// 当前版本
				if n.Header.GetRevision() > w.revision {
					w.revision = n.Header.GetRevision()
				}
				if err := n.Err(); err != nil {
					client.logger.Errorf(ctx, "etcd watch prefix found error:%v, prefix:%v", n.Err(), prefix)
					continue
				}
				for _, ev := range n.Events {
					event := &Event{
						Type: clientv3.EventTypePut,
						Kv:   &KeyValue{Key: ev.Kv.Key, Value: ev.Kv.Value},
					}
					if ev.Type == clientv3.EventTypeDelete {
						event.Type = clientv3.EventTypeDelete
					}

					client.logger.Infof(ctx, "watch type:%d key:%s", event.Type, string(event.Kv.Key))
					begin := time.Now().UnixNano() / 1e9
					w.eventChan <- event
					end := time.Now().UnixNano() / 1e9
					client.logger.Infof(ctx, "etcd watch directory take times: %dms, mod: WatchPrefix", end-begin)
				}
			}
			ctx, cancel := context.WithCancel(context.Background())
			w.cancel = cancel
			if w.revision > 0 {
				rch = client.c.Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithCreatedNotify(), clientv3.WithRev(w.revision))
			} else {
				rch = client.c.Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithCreatedNotify())
			}
		}
	}()

	return w.eventChan, nil
}

// 监听目录 阻塞执行
func (client *client) WatchDirectory(ctx context.Context, directory string, handler WatchHandler) error {
	watchCh := client.c.Watch(ctx, directory, clientv3.WithPrefix())
	for resp := range watchCh {
		for _, ev := range resp.Events {
			client.logger.Infof(ctx, "watch data:%d %s", ev.Type, string(ev.Kv.Key))
			begin := time.Now().UnixNano() / 1e9
			event := &Event{
				Type: clientv3.EventTypePut,
				Kv:   &KeyValue{Key: ev.Kv.Key, Value: ev.Kv.Value},
			}
			if ev.Type == clientv3.EventTypeDelete {
				event.Type = clientv3.EventTypeDelete
			}
			handler(ctx, event)
			end := time.Now().UnixNano() / 1e9
			client.logger.Infof(ctx, "etcd watch directory take times: %dms, mod: client_impl", end-begin)
		}
	}
	return nil
}

// Close close watch
func (w *Watch) Close() error {
	if w.cancel != nil {
		w.cancel()
	}
	return nil
}
