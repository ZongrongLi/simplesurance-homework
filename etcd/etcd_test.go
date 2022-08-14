package etcd

import (
	"context"
	"fmt"
	"log"
	"testing"

	"example.com/m/v2/config"
)

func TestEtcd_GetLock(t *testing.T) {
	InitEtcd(config.EtcdEndpoints, config.EtcdTTl)
	l, err := MustGetEtcd().GetLock(config.LockKey)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := l.Lock(ctx); err != nil {
		t.Errorf("%+v\n", err)
	}

	defer func() {
		if err := l.Unlock(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	fmt.Println("hello world")
}
