package etcd

import (
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
)

var EtcdV3 *Etcd
var initEtcdOnce sync.Once

func InitEtcd(endpoints []string, timeout time.Duration) {
	initEtcdOnce.Do(func() {
		var err error
		EtcdV3, err = newEtcd(endpoints, timeout)
		if err != nil {
			panic(err)
		}
	})
}

func MustGetEtcd() *Etcd {
	if EtcdV3 == nil {
		panic("etcd init error")
	}
	return EtcdV3
}

type Etcd struct {
	endpoints []string
	client    *clientv3.Client
	kv        clientv3.KV
	timeout   time.Duration
}

func newEtcd(endpoints []string, timeout time.Duration) (*Etcd, error) {
	conf := clientv3.Config{
		Username:    "root",
		Password:    "Cet1f7qTjo",
		Endpoints:   endpoints,
		DialTimeout: timeout,
	}
	client, err := clientv3.New(conf)
	if err != nil {
		return nil, err
	}

	return &Etcd{
		endpoints: endpoints,
		client:    client,
		kv:        clientv3.NewKV(client),
		timeout:   timeout,
	}, nil
}

func (etcd *Etcd) GetLock(path string) (*concurrency.Mutex, error) {
	s, err := concurrency.NewSession(etcd.client)
	if err != nil {
		return nil, err
	}
	l := concurrency.NewMutex(s, path)

	return l, nil
}
